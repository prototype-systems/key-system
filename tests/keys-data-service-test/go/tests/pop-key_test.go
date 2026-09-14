package tests

import (
	"net/http"
	"testing"

	"keys-data-service-test/go/settings"
	"keys-data-service-test/go/utilities"
)

func TestPopKeyFeature(tester *testing.T) {
	utilities.SetupTest(tester)

	pushBody := map[string]any{
		"reference": "test-pop-reference",
		"name":      "test-pop-key",
		"value":     "test pop key value",
	}

	pushResponse := utilities.SendTestPost(tester, settings.TargetAddress+"/service/data/keys/push", pushBody)
	if pushResponse.StatusCode != http.StatusOK {
		tester.Fatalf("pop key setup push status = %d, want %d", pushResponse.StatusCode, http.StatusOK)
	}

	pushResponse.Body.Close()

	popResponse := utilities.SendTestDelete(tester, settings.TargetAddress+"/service/data/keys/pop/test-pop-reference")
	if popResponse.StatusCode != http.StatusOK {
		tester.Fatalf("pop key status = %d, want %d", popResponse.StatusCode, http.StatusOK)
	}

	key := utilities.ParseKey(tester, popResponse)
	if key.Reference != "test-pop-reference" {
		tester.Fatalf("pop key reference = %q, want %q", key.Reference, "test-pop-reference")
	}

	if key.Name != "test-pop-key" {
		tester.Fatalf("pop key name = %q, want %q", key.Name, "test-pop-key")
	}
}
