package middle

import (
	"context"
	"github.com/tidwall/buntdb"
	"go.rls.moe/nyx/config"
	"go.rls.moe/nyx/resources"
	"net/http"
)

func Database(c *config.Config) (func(http.Handler) http.Handler, error) {
	db, err := buntdb.Open(c.DB.File)
	if err != nil {
		return nil, err
	}
	if err = resources.InitialSetup(db); err != nil {
		return nil, err
	}
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			r = r.WithContext(context.WithValue(r.Context(),
				dbCtxKey, db))
			next.ServeHTTP(w, r)
		})
	}, nil
}

func GetDB(r *http.Request) *buntdb.DB {
	val := r.Context().Value(dbCtxKey)
	if val == nil {
		panic("DB Middleware not configured")
	}
	return val.(*buntdb.DB)
}
