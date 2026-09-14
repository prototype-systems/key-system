package tests

import (
	"net/http"
	"testing"

	"keys-data-service-test/go/settings"
	"keys-data-service-test/go/utilities"
)

func TestGetKeysFeature(tester *testing.T) {
	utilities.SetupTest(tester)

	pushBody := map[string]any{
		"reference": "test-get-keys-reference",
		"name":      "test-get-keys",
		"value":     "test get keys value",
	}

	pushResponse := utilities.SendTestPost(tester, settings.TargetAddress+"/service/data/keys/push", pushBody)
	if pushResponse.StatusCode != http.StatusOK {
		tester.Fatalf("get keys setup push status = %d, want %d", pushResponse.StatusCode, http.StatusOK)
	}

	pushResponse.Body.Close()

	getResponse := utilities.SendTestGet(tester, settings.TargetAddress+"/service/data/keys")
	if getResponse.StatusCode != http.StatusOK {
		tester.Fatalf("get keys status = %d, want %d", getResponse.StatusCode, http.StatusOK)
	}

	result := utilities.ParseKeys(tester, getResponse)
	if result.Total < 1 {
		tester.Fatalf("get keys total = %d, want at least 1", result.Total)
	}

	if len(result.Keys) < 1 {
		tester.Fatalf("get keys count = %d, want at least 1", len(result.Keys))
	}

	if result.Limit <= 0 {
		tester.Fatalf("get keys limit = %d, want positive value", result.Limit)
	}

	if result.Pages < 1 {
		tester.Fatalf("get keys pages = %d, want at least 1", result.Pages)
	}
}
