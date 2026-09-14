package engines

import (
	"encoding/json"
	"fmt"

	"github.com/tidwall/buntdb"

	"keys-data-service/data"
	"keys-data-service/libraries"
	"keys-data-service/settings"
)

func GroupCountKeys(cache *data.Cache, group string) (int, error) {
	var issues libraries.Issues
	var issue error

	if group == "" {
		return 0, issues.Raise("get group keys count with issue, undefined group")
	}

	var pivot []byte
	if pivot, issue = json.Marshal(map[string]string{"group": group}); issue != nil {
		return 0, issues.Raise("get group keys count on build pivot with issue, %w", issue)
	}

	var count int
	if issue = cache.DB().View(func(transaction *buntdb.Tx) error {
		issue = transaction.AscendEqual(settings.GroupIndexName, string(pivot), func(_, _ string) bool {
			count++

			return true
		})

		if issue != nil {
			issues.Raise("iterate keys on get group keys count with issue, %w", issue)
		}

		return issues.Flush()

	}); issue != nil {
		return 0, fmt.Errorf("get group keys count with issue, %w", issue)
	}

	return count, nil
}
