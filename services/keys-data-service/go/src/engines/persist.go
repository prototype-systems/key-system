package engines

import (
	"fmt"

	"keys-data-service/data"
)

func PersistCache(cache *data.Cache) error {
	if issue := cache.Persist(); issue != nil {
		return fmt.Errorf("persist cache with issue, %w", issue)
	}

	return nil
}
