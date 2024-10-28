package main

import (
	"context"
	"fmt"
	"net"
	"net/http"
)

type contentKey string

const contentUserKey contentKey = "user_ip"

func (app *application) ipFromContext(ctx context.Context) string {
	return ctx.Value(contentUserKey).(string)
}

func (app *application) addIPToContext(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		var ctx = context.Background()
		ip, err := getIP(r)
		if err != nil {
			ip, _, _ := net.SplitHostPort(r.RemoteAddr)
			if len(ip) == 0 {
				ip = "unknown"
			}
			ctx = context.WithValue(r.Context(), contentUserKey, ip)
		} else {
			ctx = context.WithValue(r.Context(), contentUserKey, ip)
		}
		next.ServeHTTP(w, r.WithContext(ctx))
	})
}

func getIP(r *http.Request) (string, error) {
	ip, _, err := net.SplitHostPort(r.RemoteAddr)

	if err != nil {
		return "unknown", err
	}

	userIp := net.ParseIP(ip)
	if userIp == nil {
		return "", fmt.Errorf("userip: %q is not IP:port", r.RemoteAddr)
	}

	forward := r.Header.Get("X-Forwarded-For")
	if len(forward) > 0 {
		ip = forward
	}

	return ip, nil
}

func (app *application) auth(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if !app.Session.Exists(r.Context(), "user_id") {
			app.Session.Put(r.Context(), "error", "Log in first!")
			http.Redirect(w, r, "/", http.StatusTemporaryRedirect)
			return
		}
		next.ServeHTTP(w, r)
	})
}
