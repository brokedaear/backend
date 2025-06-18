// SPDX-FileCopyrightText: 2025 BROKE DA EAR LLC <https://brokedaear.com>
//
// SPDX-License-Identifier: Apache-2.0

package server

import (
	"context"
	"sync"

	"go.opentelemetry.io/otel/attribute"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/health/grpc_health_v1"
	"google.golang.org/grpc/status"

	"backend.brokedaear.com/internal/common/telemetry"
)

// TelemetryHealthServer is an enhanced health server with telemetry integration.
type TelemetryHealthServer struct {
	grpc_health_v1.UnimplementedHealthServer
	mu               sync.RWMutex
	services         map[string]grpc_health_v1.HealthCheckResponse_ServingStatus
	watchers         map[string][]*healthWatcher
	logger           Logger
	telemetry        telemetry.Telemetry
	healthMetrics    *HealthMetrics
}

// NewTelemetryHealthServer creates a new health server with telemetry integration.
func NewTelemetryHealthServer(logger Logger, tel telemetry.Telemetry) (*TelemetryHealthServer, error) {
	healthMetrics, err := NewHealthMetrics(tel)
	if err != nil {
		return nil, err
	}

	return &TelemetryHealthServer{
		UnimplementedHealthServer: grpc_health_v1.UnimplementedHealthServer{},
		mu:                        sync.RWMutex{},
		services:                  make(map[string]grpc_health_v1.HealthCheckResponse_ServingStatus),
		watchers:                  make(map[string][]*healthWatcher),
		logger:                    logger,
		telemetry:                 tel,
		healthMetrics:             healthMetrics,
	}, nil
}

// SetServingStatus sets the serving status for a specific service with telemetry.
func (h *TelemetryHealthServer) SetServingStatus(service string, status grpc_health_v1.HealthCheckResponse_ServingStatus) {
	ctx := context.Background()
	
	// Start a trace for this operation
	ctx, span := h.telemetry.TraceStart(ctx, "health.set_status")
	defer span.End()
	
	span.SetAttributes(
		attribute.String("service", service),
		attribute.String("status", status.String()),
	)

	h.mu.Lock()
	defer h.mu.Unlock()
	
	oldStatus, existed := h.services[service]
	h.services[service] = status
	
	// Record metrics
	h.healthMetrics.RecordHealthStatus(ctx, service, status)
	
	// Log status change
	if existed && oldStatus != status {
		h.logger.Info("health status changed", 
			"service", service, 
			"old_status", oldStatus.String(),
			"new_status", status.String())
	} else if !existed {
		h.logger.Info("health status set", 
			"service", service, 
			"status", status.String())
	}
	
	// Notify watchers of the status change.
	if watchers, exists := h.watchers[service]; exists {
		for _, watcher := range watchers {
			err := watcher.stream.Send(&grpc_health_v1.HealthCheckResponse{
				Status: status,
			})
			if err != nil {
				h.logger.Warn("failed to send health status to watcher", "service", service, "error", err)
			}
		}
	}
}

// Check implements the health check RPC method with telemetry.
func (h *TelemetryHealthServer) Check(ctx context.Context, req *grpc_health_v1.HealthCheckRequest) (*grpc_health_v1.HealthCheckResponse, error) {
	// Start a trace for this health check
	ctx, span := h.telemetry.TraceStart(ctx, "health.check")
	defer span.End()
	
	service := req.GetService()
	span.SetAttributes(attribute.String("service", service))
	
	h.mu.RLock()
	servingStatus, exists := h.services[service]
	h.mu.RUnlock()
	
	if !exists {
		span.SetAttributes(attribute.Bool("found", false))
		h.healthMetrics.RecordHealthCheck(ctx, service, grpc_health_v1.HealthCheckResponse_SERVICE_UNKNOWN, false)
		h.logger.Debug("health check for unknown service", "service", service)
		return nil, status.Errorf(codes.NotFound, "service %q not found", service)
	}
	
	span.SetAttributes(
		attribute.Bool("found", true),
		attribute.String("status", servingStatus.String()),
	)
	
	h.healthMetrics.RecordHealthCheck(ctx, service, servingStatus, true)
	h.logger.Debug("health check performed", "service", service, "status", servingStatus.String())
	
	return &grpc_health_v1.HealthCheckResponse{
		Status: servingStatus,
	}, nil
}

// Watch implements the streaming health check RPC method with telemetry.
func (h *TelemetryHealthServer) Watch(req *grpc_health_v1.HealthCheckRequest, stream grpc_health_v1.Health_WatchServer) error {
	ctx := stream.Context()
	ctx, span := h.telemetry.TraceStart(ctx, "health.watch")
	defer span.End()
	
	service := req.GetService()
	span.SetAttributes(attribute.String("service", service))
	
	h.mu.Lock()
	servingStatus, exists := h.services[service]
	if !exists {
		h.mu.Unlock()
		span.SetAttributes(attribute.Bool("found", false))
		h.logger.Debug("watch request for unknown service", "service", service)
		return status.Errorf(codes.NotFound, "service %q not found", service)
	}
	
	// Add watcher to the list and record metric.
	watcher := &healthWatcher{
		service: service,
		stream:  stream,
	}
	h.watchers[service] = append(h.watchers[service], watcher)
	h.mu.Unlock()
	
	// Record watcher added
	h.healthMetrics.RecordWatcherChange(ctx, service, 1)
	span.SetAttributes(attribute.Bool("found", true))
	
	// Send initial status.
	err := stream.Send(&grpc_health_v1.HealthCheckResponse{
		Status: servingStatus,
	})
	if err != nil {
		h.logger.Warn("failed to send initial health status", "service", service, "error", err)
		return err
	}
	
	h.logger.Debug("health watch started", "service", service)
	
	// Keep the stream open until context is cancelled.
	<-ctx.Done()
	
	// Remove watcher from the list and record metric.
	h.mu.Lock()
	defer h.mu.Unlock()
	
	watchers := h.watchers[service]
	for i, w := range watchers {
		if w == watcher {
			h.watchers[service] = append(watchers[:i], watchers[i+1:]...)
			break
		}
	}
	
	// Record watcher removed
	h.healthMetrics.RecordWatcherChange(ctx, service, -1)
	h.logger.Debug("health watch ended", "service", service)
	
	return nil
}

// Shutdown gracefully stops the health server with telemetry.
func (h *TelemetryHealthServer) Shutdown() {
	ctx := context.Background()
	ctx, span := h.telemetry.TraceStart(ctx, "health.shutdown")
	defer span.End()
	
	h.mu.Lock()
	defer h.mu.Unlock()
	
	// Set all services to NOT_SERVING and record metrics.
	for service := range h.services {
		h.services[service] = grpc_health_v1.HealthCheckResponse_NOT_SERVING
		h.healthMetrics.RecordHealthStatus(ctx, service, grpc_health_v1.HealthCheckResponse_NOT_SERVING)
		
		// Notify watchers.
		if watchers, exists := h.watchers[service]; exists {
			for _, watcher := range watchers {
				err := watcher.stream.Send(&grpc_health_v1.HealthCheckResponse{
					Status: grpc_health_v1.HealthCheckResponse_NOT_SERVING,
				})
				if err != nil {
					h.logger.Warn("failed to send shutdown status to watcher", "service", service, "error", err)
				}
			}
		}
	}
	
	h.logger.Info("telemetry health server shutdown complete")
}