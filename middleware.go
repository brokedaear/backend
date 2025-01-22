package main

import (
	"fmt"
	"net/http"
	"sync"
	"time"

	"github.com/tomasen/realip"
	"golang.org/x/time/rate"
)

func (app *app) recoverPanic(next http.Handler) http.Handler {
	return http.HandlerFunc(
		func(w http.ResponseWriter, r *http.Request) {
			defer func() {
				if err := recover(); err != nil {
					w.Header().Set("Connect", "close")
					app.serverError(w, r, fmt.Errorf("%s", err))
				}
			}()
		},
	)
}

// rateLimit limits the rate of requests using the golang.org/x/time/rate
// package. It also handles race conditions.
func (app *app) rateLimit(next http.Handler) http.Handler {
	type client struct {
		limiter  *rate.Limiter
		lastSeen time.Time
	}

	var (
		mu      sync.Mutex
		clients = make(map[string]*client)
	)

	// Creates a go routine that checks when a client was last seen,
	// so that we can delete old clients who haven't been seen in a
	// while.
	go func() {
		for {
			time.Sleep(time.Minute)

			mu.Lock()

			for ip, client := range clients {
				if time.Since(client.lastSeen) > 3*time.Minute {
					delete(clients, ip)
				}
			}

			mu.Unlock()
		}
	}()

	return http.HandlerFunc(
		func(w http.ResponseWriter, r *http.Request) {
			if app.config.limiter.enabled {
				ip := realip.FromRequest(r)

				mu.Lock()

				if _, found := clients[ip]; !found {
					clients[ip] = &client{
						limiter: rate.NewLimiter(
							rate.Limit(app.config.limiter.rps),
							app.config.limiter.burst,
						),
					}
				}

				clients[ip].lastSeen = time.Now()

				if !clients[ip].limiter.Allow() {
					mu.Unlock()
					app.rateLimitExceededResponse(w, r)
					return
				}

				mu.Unlock()
			}

			next.ServeHTTP(w, r)
		},
	)
}

func (app *app) enableCORS(next http.Handler) http.Handler {
	return http.HandlerFunc(
		func(w http.ResponseWriter, r *http.Request) {
			w.Header().Set("Access-Control-Allow-Origin", "*")
			w.Header().Set(
				"Access-Control-Allow-Methods",
				"GET, POST, PUT, DELETE, OPTIONS",
			)
			w.Header().Set(
				"Access-Control-Allow-Headers",
				"Content-Type, Authorization",
			)
			w.Header().Set("Access-Control-Allow-Credentials", "true")

			// If the request is for the OPTIONS method, return immediately with a 200 status
			// as this is a preflight request
			if r.Method == "OPTIONS" {
				w.WriteHeader(http.StatusOK)
				return
			}

			next.ServeHTTP(w, r)
		},
	)
}

//func (app *application) authenticate(next http.Handler) http.Handler {
//	return http.HandlerFunc(
//		func(w http.ResponseWriter, r *http.Request) {
//			w.Header().Add("Vary", "Authorization")
//
//			authorizationHeader := r.Header.Get("Authorization")
//
//			if authorizationHeader == "" {
//				r = app.contextSetUser(r, data.AnonymousUser)
//				next.ServeHTTP(w, r)
//				return
//			}
//
//			headerParts := strings.Split(authorizationHeader, " ")
//			if len(headerParts) != 2 || headerParts[0] != "Bearer" {
//				app.invalidAuthenticationTokenResponse(w, r)
//				return
//			}
//
//			token := headerParts[1]
//
//			v := validator.New()
//
//			if data.ValidateTokenPlaintext(v, token); !v.Valid() {
//				app.invalidAuthenticationTokenResponse(w, r)
//				return
//			}
//
//			user, err := app.models.Users.GetForToken(
//				data.ScopeAuthentication,
//				token,
//			)
//			if err != nil {
//				switch {
//				case errors.Is(err, data.ErrRecordNotFound):
//					app.invalidAuthenticationTokenResponse(w, r)
//				default:
//					app.serverError(w, r, err)
//				}
//				return
//			}
//
//			r = app.contextSetUser(r, user)
//
//			next.ServeHTTP(w, r)
//		},
//	)
//}
//
//func (app *application) requireAuthenticatedUser(next http.HandlerFunc) http.HandlerFunc {
//	return http.HandlerFunc(
//		func(w http.ResponseWriter, r *http.Request) {
//			user := app.contextGetUser(r)
//
//			if user.IsAnonymous() {
//				app.authenticationRequiredResponse(w, r)
//				return
//			}
//
//			next.ServeHTTP(w, r)
//		},
//	)
//}
