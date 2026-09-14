package engines

import (
	"encoding/json"
	"fmt"
	"time"

	"github.com/tidwall/buntdb"

	"keys-data-service/data"
	"keys-data-service/models"
)

func PushKey(cache *data.Cache, key models.Key) error {
	currentTime := time.Now()
	if key.CreatedAt.IsZero() {
		key.CreatedAt = currentTime
	}
	key.UpdatedAt = currentTime

	keyData, issue := json.Marshal(key)
	if issue != nil {
		return fmt.Errorf("parse key on push with issue, %w", issue)
	}

	return cache.DB().Update(func(transaction *buntdb.Tx) error {
		_, _, issue := transaction.Set("key:"+key.Reference, string(keyData), nil)
		if issue != nil {
			return fmt.Errorf("push key with issue, %w", issue)
		}

		return nil
	})
}
