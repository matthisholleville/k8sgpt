package cache

import (
	"os"
	"path/filepath"

	"github.com/adrg/xdg"
	"github.com/k8sgpt-ai/k8sgpt/pkg/util"
)

var _ (ICache) = (*FileBasedCache)(nil)

type FileBasedCache struct {
	noCache bool
}

func (f *FileBasedCache) IsCacheDisabled() bool {
	return f.noCache
}

func (*FileBasedCache) Exists(key string) (bool, error) {
	path, err := xdg.CacheFile(filepath.Join("k8sgpt", key))

	if err != nil {
		return false, err
	}

	exists, err := util.FileExists(path)

	if err != nil {
		return false, err
	}

	return exists, nil
}

func (*FileBasedCache) Load(key string) (string, error) {
	path, err := xdg.CacheFile(filepath.Join("k8sgpt", key))

	if err != nil {
		return "", err
	}

	data, err := os.ReadFile(path)

	if err != nil {
		return "", err
	}

	return string(data), nil
}

func (*FileBasedCache) Store(key string, data string) error {
	path, err := xdg.CacheFile(filepath.Join("k8sgpt", key))

	if err != nil {
		return err
	}

	return os.WriteFile(path, []byte(data), 0600)
}
