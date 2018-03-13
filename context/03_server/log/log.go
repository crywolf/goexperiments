package log

import (
	"context"
	"log"
	"math/rand"
	"net/http"
)

type requestKey int

const requestIDKey requestKey = 42

// Println logs the message with request ID
func Println(ctx context.Context, msg string) {
	requestID, ok := ctx.Value(requestIDKey).(int64)
	if !ok {
		log.Println("could not find request ID in context")
		return
	}

	log.Printf("[%d] %s", requestID, msg)
}

// Decorate adds request ID to request's context
func Decorate(handlerFunc http.HandlerFunc) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		ctx := r.Context()
		requestID := rand.Int63()
		ctx = context.WithValue(ctx, requestIDKey, requestID)
		handlerFunc(w, r.WithContext(ctx))
	}
}
