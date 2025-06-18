// SPDX-FileCopyrightText: 2025 BROKE DA EAR LLC <https://brokedaear.com>
//
// SPDX-License-Identifier: Apache-2.0

package server

import (
	"context"
	"sync"

	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/health/grpc_health_v1"
	"google.golang.org/grpc/status"
)

// HealthServer implements the gRPC health checking protocol.
// It maintains a registry of service health statuses and provides
// both point-in-time checks and streaming watch functionality.
type HealthServer struct {
	grpc_health_v1.UnimplementedHealthServer
	mu       sync.RWMutex
	services map[string]grpc_health_v1.HealthCheckResponse_ServingStatus
	watchers map[string][]*healthWatcher
	logger   Logger
}

// healthWatcher represents a streaming health check watcher.
type healthWatcher struct {
	service string
	stream  grpc_health_v1.Health_WatchServer
}

// NewHealthServer creates a new health server instance.
func NewHealthServer(logger Logger) *HealthServer {
	return &HealthServer{
		UnimplementedHealthServer: grpc_health_v1.UnimplementedHealthServer{},
		mu:                        sync.RWMutex{},
		services:                 make(map[string]grpc_health_v1.HealthCheckResponse_ServingStatus),
		watchers:                 make(map[string][]*healthWatcher),
		logger:                   logger,
	}
}

// SetServingStatus sets the serving status for a specific service.
// Use empty string for the overall server health status.
func (h *HealthServer) SetServingStatus(service string, status grpc_health_v1.HealthCheckResponse_ServingStatus) {
	h.mu.Lock()
	defer h.mu.Unlock()
	
	h.services[service] = status
	h.logger.Debug("health status updated", "service", service, "status", status.String())
	
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

// Check implements the health check RPC method.
func (h *HealthServer) Check(_ context.Context, req *grpc_health_v1.HealthCheckRequest) (*grpc_health_v1.HealthCheckResponse, error) {
	h.mu.RLock()
	defer h.mu.RUnlock()
	
	service := req.GetService()
	servingStatus, exists := h.services[service]
	
	if !exists {
		h.logger.Debug("health check for unknown service", "service", service)
		return nil, status.Errorf(codes.NotFound, "service %q not found", service)
	}
	
	h.logger.Debug("health check performed", "service", service, "status", servingStatus.String())
	return &grpc_health_v1.HealthCheckResponse{
		Status: servingStatus,
	}, nil
}

// Watch implements the streaming health check RPC method.
func (h *HealthServer) Watch(req *grpc_health_v1.HealthCheckRequest, stream grpc_health_v1.Health_WatchServer) error {
	service := req.GetService()
	
	h.mu.Lock()
	servingStatus, exists := h.services[service]
	if !exists {
		h.mu.Unlock()
		h.logger.Debug("watch request for unknown service", "service", service)
		return status.Errorf(codes.NotFound, "service %q not found", service)
	}
	
	// Add watcher to the list.
	watcher := &healthWatcher{
		service: service,
		stream:  stream,
	}
	h.watchers[service] = append(h.watchers[service], watcher)
	h.mu.Unlock()
	
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
	<-stream.Context().Done()
	
	// Remove watcher from the list.
	h.mu.Lock()
	defer h.mu.Unlock()
	
	watchers := h.watchers[service]
	for i, w := range watchers {
		if w == watcher {
			h.watchers[service] = append(watchers[:i], watchers[i+1:]...)
			break
		}
	}
	
	h.logger.Debug("health watch ended", "service", service)
	return nil
}

// Shutdown gracefully stops the health server and notifies all watchers.
func (h *HealthServer) Shutdown() {
	h.mu.Lock()
	defer h.mu.Unlock()
	
	// Set all services to NOT_SERVING.
	for service := range h.services {
		h.services[service] = grpc_health_v1.HealthCheckResponse_NOT_SERVING
		
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
	
	h.logger.Info("health server shutdown complete")
}