package engines

import (
	"fmt"

	"keys-data-service/data"
	"keys-data-service/libraries"
	"keys-data-service/models"
)

func AddKey(cache *data.Cache, key models.Key) (string, error) {
	reference, issue := libraries.ConstructReference(key.Name, key.Value, key.Group)
	if issue != nil {
		return "", fmt.Errorf("add key on construct reference with issue, %w", issue)
	}

	key.Reference = reference

	issue = PushKey(cache, key)
	if issue != nil {
		return "", fmt.Errorf("add key with issue, %w", issue)
	}

	return reference, nil
}
