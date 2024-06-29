package middleware

import (
	"net/http"
	"runtime/debug"

	"go.uber.org/zap"
)

func LogMiddleware(next http.HandlerFunc, logger zap.SugaredLogger) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		logger.Infow(
			"entering",
			"method", r.Method,
			"url", r.RequestURI,
		)

		next(w, r)

		logger.Infow(
			"leaving",
			"method", r.Method,
			"url", r.RequestURI,
		)
	}
}

func PanicMiddleware(next http.HandlerFunc, logger zap.SugaredLogger) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		defer func() {
			if err := recover(); err != nil {
				logger.Errorw(
					"panic happened",
					"reason", err,
					"stackTrace", string(debug.Stack()),
				)

				w.WriteHeader(http.StatusInternalServerError)
			}
		}()

		next(w, r)
	}
}

func ErrorMiddleware(next func(http.ResponseWriter, *http.Request) error, logger zap.SugaredLogger) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		if err := next(w, r); err != nil {
			logger.Errorw(
				"error happened",
				"url", r.RequestURI,
				"reason", err,
			)

			w.WriteHeader(http.StatusBadRequest)
		}
	}
}

func Accept(f func(http.ResponseWriter, *http.Request) error, logger zap.SugaredLogger) http.HandlerFunc {
	middlewares := []func(next http.HandlerFunc, logger zap.SugaredLogger) http.HandlerFunc{
		PanicMiddleware,
		LogMiddleware,
	}

	prelude := ErrorMiddleware(f, logger)

	for i := range middlewares {
		prelude = middlewares[i](prelude, logger)
	}

	return prelude
}
