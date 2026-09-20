package engines

import (
	"encoding/json"
	"fmt"
	"time"

	"github.com/tidwall/buntdb"

	"keys-data-service/data"
	"keys-data-service/libraries"
	"keys-data-service/models"
	"keys-data-service/settings"
)

func GroupSetKey(cache *data.Cache, group string, name string, value json.RawMessage) (string, error) {
	var issues libraries.Issues
	var issue error

	if group == "" {
		return "", issues.Raise("set key by group with issue, undefined group")
	}

	if name == "" {
		return "", issues.Raise("set key by group with issue, undefined key name")
	}

	if libraries.CheckNullValue(value) {
		return "", issues.Raise("set key by group with issue, undefined key value")
	}

	var pivot []byte
	if pivot, issue = json.Marshal(map[string]string{"group": group}); issue != nil {
		return "", issues.Raise("set key by group on build pivot with issue, %w", issue)
	}

	var reference string

	if issue = cache.DB().Update(func(transaction *buntdb.Tx) error {
		var target *models.Key

		issue = transaction.AscendEqual(settings.GroupIndexName, string(pivot), func(_, raw string) bool {
			var candidate models.Key
			if issue = json.Unmarshal([]byte(raw), &candidate); issue != nil {
				issues.Raise("parse key on set by group with issue, %w", issue)

				return false
			}

			if candidate.Name == name {
				target = &candidate

				return false
			}

			return true
		})

		if issue != nil {
			issues.Raise("iterate keys on set by group with issue, %w", issue)
		}

		if !issues.Empty() {
			return issues.Flush()
		}

		currentTime := time.Now()

		if target != nil {
			target.Value = value
			target.UpdatedAt = currentTime

			var keyData []byte
			if keyData, issue = json.Marshal(target); issue != nil {
				return issues.Raise("parse key on set by group update with issue, %w", issue)
			}

			if _, _, issue = transaction.Set("key:"+target.Reference, string(keyData), nil); issue != nil {
				return issues.Raise("set key by group on update with issue, %w", issue)
			}

			reference = target.Reference

			return nil
		}

		var newReference string
		if newReference, issue = libraries.ConstructReference(name, value, group); issue != nil {
			return issues.Raise("set key by group on construct reference with issue, %w", issue)
		}

		key := models.Key{
			Reference: newReference,
			Name:      name,
			Group:     group,
			Value:     value,
			CreatedAt: currentTime,
			UpdatedAt: currentTime,
		}

		var keyData []byte
		if keyData, issue = json.Marshal(key); issue != nil {
			return issues.Raise("parse key on set by group insert with issue, %w", issue)
		}

		if _, _, issue = transaction.Set("key:"+newReference, string(keyData), nil); issue != nil {
			return issues.Raise("set key by group on insert with issue, %w", issue)
		}

		reference = newReference

		return nil

	}); issue != nil {
		return "", fmt.Errorf("set key by group with issue, %w", issue)
	}

	return reference, nil
}
