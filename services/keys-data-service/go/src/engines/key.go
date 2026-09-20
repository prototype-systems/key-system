package engines

import (
	"encoding/json"
	"fmt"

	"github.com/tidwall/buntdb"

	"keys-data-service/data"
	"keys-data-service/models"
)

func GetKey(cache *data.Cache, reference string) (*models.Key, error) {
	if reference == "" {
		return nil, fmt.Errorf("get key with issue, undefined key reference")
	}

	var key models.Key

	issue := cache.DB().View(func(transaction *buntdb.Tx) error {
		rawKey, issue := transaction.Get("key:" + reference)
		if issue != nil {
			return fmt.Errorf("get key by reference %s with issue, %w", reference, issue)
		}

		if issue := json.Unmarshal([]byte(rawKey), &key); issue != nil {
			return fmt.Errorf("parse key on get with issue, %w", issue)
		}

		return nil
	})

	if issue != nil {
		return nil, issue
	}

	return &key, nil
}
