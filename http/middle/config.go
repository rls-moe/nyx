package middle

import (
	"context"
	"go.rls.moe/nyx/config"
	"net/http"
)

func ConfigCtx(config *config.Config) func(http.Handler) http.Handler {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			r = r.WithContext(context.WithValue(r.Context(), configKey, config))
			next.ServeHTTP(w, r)
		})
	}
}

func GetConfig(r *http.Request) *config.Config {
	val := r.Context().Value(configKey)
	if val == nil {
		panic("Config Middleware not configured")
	}
	return val.(*config.Config)
}
