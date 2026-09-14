package libraries

import (
	"encoding/json"
	"fmt"
	"hash/fnv"
)

func ConstructReference(data ...any) (string, error) {
	hash := fnv.New64a()

	for _, item := range data {
		bytes, issue := json.Marshal(item)

		if issue != nil {
			return "", issue
		}

		hash.Write(bytes)
	}

	return fmt.Sprintf("%016x", hash.Sum64()), nil
}
