package middleware

import (
	"net/http"
	"time"

	"go.uber.org/zap"
)

type Function func(next http.HandlerFunc) http.HandlerFunc

type Onion struct {
	middlewares []Function
	log         *zap.Logger
	duration    time.Duration
}

func NewOnion(lg *zap.Logger) *Onion {
	return &Onion{
		log: lg,
	}
}

func (o *Onion) Apply(h http.HandlerFunc) http.HandlerFunc {
	for i := range o.middlewares {
		h = o.middlewares[i](h)
	}
	return h
}

func (o *Onion) AppendMiddleware(mw ...Function) {
	o.middlewares = append(o.middlewares, mw...)
}

func (o *Onion) LogRequestResponse(next http.HandlerFunc) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		o.log.Info("Request recieved",
			zap.String("method", r.Method),
			zap.String("requestURI", r.RequestURI),
			zap.String("host", r.Host),
			zap.Time("time", time.Now()),
		)

		next(w, r)

		o.log.Info("Response sent",
			zap.Duration("time to handle", o.duration),
		)
	}
}

func (o *Onion) Timer(next http.HandlerFunc) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		start := time.Now()
		next(w, r)
		o.duration = time.Since(start)
	}
}
