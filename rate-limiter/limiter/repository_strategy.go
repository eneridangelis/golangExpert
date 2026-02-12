package limiter

type RepositoryStrategy interface {
	IncrementAndGet(key string, expirationSeconds int) (int64, error)
	Close() error
}
