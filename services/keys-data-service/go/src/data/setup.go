package data

import "github.com/tidwall/buntdb"

func Setup(cacheName string, connection *buntdb.DB, path string) (*Cache, error) {
	if issue := Index(connection); issue != nil {
		return nil, issue
	}

	return &Cache{name: cacheName, connection: connection, path: path}, nil
}
