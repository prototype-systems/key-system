package data

import (
	"fmt"
	"keys-data-service/libraries"
	"keys-data-service/settings"
	"os"
	"path/filepath"
	"sync"

	"github.com/tidwall/buntdb"
)

type Cache struct {
	name       string
	connection *buntdb.DB
	path       string
}

var (
	dataCache        *Cache
	dataCacheControl sync.RWMutex
)

func SetCache(cache *Cache) {
	dataCacheControl.Lock()
	defer dataCacheControl.Unlock()

	dataCache = cache
}

func GetCache() *Cache {
	dataCacheControl.RLock()
	defer dataCacheControl.RUnlock()

	return dataCache
}

func (cache *Cache) DB() *buntdb.DB {
	return cache.connection
}

func (cache *Cache) Name() string {
	return cache.name
}

func (cache *Cache) Path() string {
	return cache.path
}

func Open() (*Cache, error) {
	cacheDirectory := os.Getenv(settings.CacheDirectoryVariable)

	if !libraries.CheckDirectory(cacheDirectory) {
		executablePath, issue := os.Executable()
		if issue != nil {
			return nil, fmt.Errorf("get process path on open with issue, %w", issue)
		}

		resolvedExecutablePath, issue := filepath.EvalSymlinks(executablePath)
		if issue != nil {
			return nil, fmt.Errorf("get resolved process path on open with issue, %w", issue)
		}

		cacheDirectory = filepath.Join(filepath.Dir(resolvedExecutablePath), settings.DefaultCacheDirectory)
	}

	if !libraries.CheckDirectory(cacheDirectory) {
		if issue := os.MkdirAll(cacheDirectory, 0755); issue != nil {
			return nil, fmt.Errorf("create cache directory on open with issue, %w", issue)
		}
	}

	cachePath := filepath.Join(cacheDirectory, settings.DefaultCacheFile)

	connection, issue := buntdb.Open(":memory:")
	if issue != nil {
		return nil, issue
	}

	if libraries.CheckFile(cachePath) {
		cacheFileHandle, issue := os.Open(cachePath)
		if issue != nil {
			return nil, fmt.Errorf("open cache with issue, %w", issue)
		}
		defer cacheFileHandle.Close()

		if issue := connection.Load(cacheFileHandle); issue != nil {
			return nil, fmt.Errorf("load cache snapshot on open with issue, %w", issue)
		}
	}

	if issue := Index(connection); issue != nil {
		return nil, issue
	}

	return &Cache{
		name:       settings.DefaultCacheName,
		connection: connection,
		path:       cachePath,
	}, nil
}

func (cache *Cache) Persist() error {
	if cache.path == "" {
		return nil
	}

	if issue := os.MkdirAll(filepath.Dir(cache.path), 0755); issue != nil {
		return fmt.Errorf("create cache directory on persist with issue, %w", issue)
	}

	persistFile, issue := os.Create(cache.path)
	if issue != nil {
		return fmt.Errorf("create cache file on persist with issue, %w", issue)
	}
	defer persistFile.Close()

	if issue := cache.connection.Save(persistFile); issue != nil {
		return fmt.Errorf("save cache snapshot on persist with issue, %w", issue)
	}

	return nil
}

func (cache *Cache) Close() error {
	if cache.path != "" {
		if issue := os.MkdirAll(filepath.Dir(cache.path), 0755); issue != nil {
			return fmt.Errorf("create cache directory on close with issue, %w", issue)
		}

		fileHandle, issue := os.Create(cache.path)
		if issue != nil {
			return fmt.Errorf("create cache file on close with issue, %w", issue)
		}
		defer fileHandle.Close()

		if issue := cache.connection.Save(fileHandle); issue != nil {
			return fmt.Errorf("save cache snapshot on close with issue, %w", issue)
		}
	}

	return cache.connection.Close()
}
