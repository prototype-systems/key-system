package tests

import (
	"net/http"
	"testing"

	"keys-data-service-test/settings"
	"keys-data-service-test/utilities"
)

func TestAddKeyFeature(tester *testing.T) {
	utilities.SetupTest(tester)

	body := map[string]any{
		"name":  "test-add-key",
		"value": "test add key value",
	}

	response := utilities.SendTestPost(tester, settings.TargetAddress+"/service/data/keys/add", body)
	if response.StatusCode != http.StatusOK {
		tester.Fatalf("add key status = %d, want %d", response.StatusCode, http.StatusOK)
	}

	reference := utilities.ParseReference(tester, response)
	if reference == "" {
		tester.Fatal("add key reference is empty")
	}
}
