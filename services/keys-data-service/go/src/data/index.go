package data

import (
	"fmt"
	"keys-data-service/settings"

	"github.com/tidwall/buntdb"
)

func Index(connection *buntdb.DB) error {
	if issue := connection.CreateIndex(
		settings.DefaultCacheName,
		"key:*",
		buntdb.IndexString,
	); issue != nil {
		return fmt.Errorf("create default index on setup with issue, %w", issue)
	}

	if issue := connection.CreateIndex(
		settings.ReferenceIndexName,
		"key:*",
		buntdb.IndexJSON("reference"),
	); issue != nil {
		return fmt.Errorf("create reference index on setup with issue, %w", issue)
	}

	if issue := connection.CreateIndex(
		settings.GroupIndexName,
		"key:*",
		buntdb.IndexJSON("group"),
	); issue != nil {
		return fmt.Errorf("create group index on setup with issue, %w", issue)
	}

	return nil
}
