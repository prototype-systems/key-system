package tests

import (
	"net/http"
	"testing"

	"keys-data-service-test/go/settings"
	"keys-data-service-test/go/utilities"
)

func TestPushKeyFeature(tester *testing.T) {
	utilities.SetupTest(tester)

	body := map[string]any{
		"reference": "test-push-reference",
		"name":      "test-push-key",
		"value":     "test push key value",
	}

	response := utilities.SendTestPost(tester, settings.TargetAddress+"/service/data/keys/push", body)
	if response.StatusCode != http.StatusOK {
		tester.Fatalf("push key status = %d, want %d", response.StatusCode, http.StatusOK)
	}

	reference := utilities.ParseReference(tester, response)
	if reference != "test-push-reference" {
		tester.Fatalf("push key reference = %q, want %q", reference, "test-push-reference")
	}
}
