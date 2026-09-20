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

func GroupGetKey(cache *data.Cache, group string, name string) (json.RawMessage, error) {
	var issues libraries.Issues
	var issue error

	if group == "" {
		return nil, issues.Raise("get key by group with issue, undefined group")
	}

	if name == "" {
		return nil, issues.Raise("get key by group with issue, undefined key name")
	}

	var pivot []byte
	if pivot, issue = json.Marshal(map[string]string{"group": group}); issue != nil {
		return nil, issues.Raise("get key by group on build pivot with issue, %w", issue)
	}

	var value json.RawMessage

	if issue = cache.DB().View(func(transaction *buntdb.Tx) error {
		issue = transaction.AscendEqual(settings.GroupIndexName, string(pivot), func(_, raw string) bool {
			var candidate models.Key
			if issue = json.Unmarshal([]byte(raw), &candidate); issue != nil {
				issues.Raise("parse key on get by group with issue, %w", issue)

				return false
			}

			if candidate.Name == name {
				value = candidate.Value

				return false
			}

			return true
		})

		if issue != nil {
			issues.Raise("iterate keys on get by group with issue, %w", issue)
		}

		return issues.Flush()
	}); issue != nil {
		return nil, fmt.Errorf("get key by group with issue, %w", issue)
	}

	if value == nil {
		return nil, nil
	}

	return value, nil
}
