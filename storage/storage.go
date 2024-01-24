package storage

import "io"

type Storage interface {
	List(bucket, prefix string) ([]string, error)
	Get(bucket, key string) (io.ReadCloser, error)
}
