package cache

import (
	"github.com/spf13/viper"
)

type ICache interface {
	Store(key string, data string) error
	Load(key string) (string, error)
	List() ([]string, error)
	Exists(key string) bool
	IsCacheDisabled() bool
}

func New(noCache bool, remoteCache bool) ICache {
	if remoteCache {
		return NewS3Cache(noCache)
	}
	return &FileBasedCache{
		noCache: noCache,
	}
}

// CacheProvider is the configuration for the cache provider when using a remote cache
type CacheProvider struct {
	BucketName string `mapstructure:"bucketname"`
	Region     string `mapstructure:"region"`
}

func FetchCacheInfo() (CacheProvider, error) {
	var cacheInfo CacheProvider
	err := viper.UnmarshalKey("cache", &cacheInfo)
	if err != nil {
		return cacheInfo, err
	}
	return cacheInfo, nil
}

func RemoteCacheEnabled() (bool, error) {
	// load remote cache if it is configured

	cacheInfo, err := FetchCacheInfo()
	if err != nil {
		return false, err
	}
	if cacheInfo.BucketName != "" && cacheInfo.Region != "" {
		return true, nil
	}
	return false, nil
}

func CacheAlreadyConfigured(cacheInfo CacheProvider) bool {
	if cacheInfo.BucketName != "" {
		return true
	}
	return false

}

func AddRemoteCache(bucketName string, region string, cacheInfo CacheProvider) error {
	cacheInfo.BucketName = bucketName
	cacheInfo.Region = region
	viper.Set("cache", cacheInfo)
	err := viper.WriteConfig()
	if err != nil {
		return err
	}
	return nil
}

func RemoveRemoteCache(bucketName string, cacheInfo CacheProvider) error {
	cacheInfo = CacheProvider{}
	viper.Set("cache", cacheInfo)
	err := viper.WriteConfig()
	if err != nil {
		return err
	}

	return nil

}
