// SPDX-FileCopyrightText: 2025 BROKE DA EAR LLC <https://brokedaear.com>
//
// SPDX-License-Identifier: Apache-2.0

package main

import (
	"encoding/json"
	"io"
	"net/http"

	"github.com/stripe/stripe-go"
	"github.com/stripe/stripe-go/webhook"
)

const MAX_BODY_BYTES = int64(65536)

// https://docs.stripe.com/webhooks?lang=go

func (a *app) StripePayment(w http.ResponseWriter, req *http.Request) {
	req.Body = http.MaxBytesReader(w, req.Body, MAX_BODY_BYTES)
	payload, err := io.ReadAll(req.Body)
	if err != nil {
		a.logger.Error("Error reading request body", "msg", err)
		w.WriteHeader(http.StatusServiceUnavailable)
		return
	}

	// If you are testing your webhook locally with the Stripe CLI you
	// can find the endpoint's secret by running `stripe listen`
	// Otherwise, find your endpoint's secret in your webhook settings
	// in the Developer Dashboard

	// Pass the request body and Stripe-Signature header to ConstructEvent,
	// along with the webhook signing key.
	event, err := webhook.ConstructEvent(payload, req.Header.Get("Stripe-Signature"),
		a.config.secrets.stripeWebhookSecret)
	if err != nil {
		a.logger.Error("Error verifying webhook signature", "msg", err)
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	// Unmarshal the event data into an appropriate struct depending on its Type
	switch event.Type {
	case "payment_intent.succeeded":
		var paymentIntent stripe.PaymentIntent

		err := json.Unmarshal(event.Data.Raw, &paymentIntent)
		if err != nil {
			a.logger.Error("Error parsing webhook JSON", "msg", err)
			w.WriteHeader(http.StatusBadRequest)
			return
		}
		a.logger.Info("PaymentIntent successful!")

	case "payment_method.attached":
		var paymentMethod stripe.PaymentMethod

		err := json.Unmarshal(event.Data.Raw, &paymentMethod)
		if err != nil {
			a.logger.Error("Error parsing webhook JSON", "msg", err)
			w.WriteHeader(http.StatusBadRequest)
			return
		}

		a.logger.Info("PaymentMethod attached to a customer!")

	default:
		a.logger.Info("Unhandled event type")
	}

	w.WriteHeader(http.StatusOK)
}
