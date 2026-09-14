package tests

import (
	"net/http"
	"testing"

	"keys-data-service-test/go/settings"
	"keys-data-service-test/go/utilities"
)

func TestGetKeyFeature(tester *testing.T) {
	utilities.SetupTest(tester)

	pushBody := map[string]any{
		"reference": "test-get-key-reference",
		"name":      "test-get-key",
		"value":     "test get key value",
	}

	pushResponse := utilities.SendTestPost(tester, settings.TargetAddress+"/service/data/keys/push", pushBody)
	if pushResponse.StatusCode != http.StatusOK {
		tester.Fatalf("get key setup push status = %d, want %d", pushResponse.StatusCode, http.StatusOK)
	}

	pushResponse.Body.Close()

	getResponse := utilities.SendTestGet(tester, settings.TargetAddress+"/service/data/keys/test-get-key-reference")
	if getResponse.StatusCode != http.StatusOK {
		tester.Fatalf("get key status = %d, want %d", getResponse.StatusCode, http.StatusOK)
	}

	key := utilities.ParseKey(tester, getResponse)
	if key.Reference != "test-get-key-reference" {
		tester.Fatalf("get key reference = %q, want %q", key.Reference, "test-get-key-reference")
	}

	if key.Name != "test-get-key" {
		tester.Fatalf("get key name = %q, want %q", key.Name, "test-get-key")
	}
}
