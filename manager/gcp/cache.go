package gcp

type Cache interface {
	Get(key string) (string, bool)
}
