package engines

import (
	"encoding/json"
	"fmt"

	"github.com/tidwall/buntdb"

	"keys-data-service/data"
	"keys-data-service/libraries"
	"keys-data-service/models"
	"keys-data-service/settings"
)

func GetKeys(cache *data.Cache, skip int, limit int) ([]models.Key, error) {
	var issues libraries.Issues
	var issue error

	if skip < 0 {
		skip = 0
	}

	if limit <= 0 {
		limit = 10
	}

	var keys []models.Key

	if issue = cache.DB().View(func(transaction *buntdb.Tx) error {
		index := 0
		count := 0

		issue = transaction.Ascend(settings.ReferenceIndexName, func(reference, value string) bool {
			if index < skip {
				index++
				return true
			}

			if count >= limit {
				return false
			}

			var key models.Key
			if issue = json.Unmarshal([]byte(value), &key); issue != nil {
				issues.Raise("get keys on reference %s with issue, %w", reference, issue)

				return false
			}

			keys = append(keys, key)

			count++
			index++

			return true
		})

		if issue != nil {
			issues.Raise("get keys on iterate with issue, %w", issue)
		}

		return issues.Flush()
	}); issue != nil {
		return nil, fmt.Errorf("get keys with issue, %w", issue)
	}

	return keys, nil
}

func GetKeysCount(cache *data.Cache) (int, error) {
	var count int

	issue := cache.DB().View(func(transaction *buntdb.Tx) error {
		transaction.Ascend(cache.Name(), func(key, value string) bool {
			count++
			return true
		})

		return nil
	})

	if issue != nil {
		return 0, fmt.Errorf("get keys count with issue, %w", issue)
	}

	return count, nil
}
