package middle

type ctxKey int64

const (
	configKey ctxKey = iota
	dbCtxKey
	sessionKey
)
