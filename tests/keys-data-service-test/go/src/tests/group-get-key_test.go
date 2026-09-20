package tests

import (
	"encoding/json"
	"net/http"
	"testing"

	"keys-data-service-test/settings"
	"keys-data-service-test/utilities"
)

func TestGroupGetKeyFeature(tester *testing.T) {
	utilities.SetupTest(tester)

	setBody := map[string]any{
		"value": "test group get key value",
	}

	setResponse := utilities.SendTestPost(tester, settings.TargetAddress+"/service/data/keys/test-group-get/test-group-get-key", setBody)
	if setResponse.StatusCode != http.StatusOK {
		tester.Fatalf("group get key setup set status = %d, want %d", setResponse.StatusCode, http.StatusOK)
	}

	setResponse.Body.Close()

	getResponse := utilities.SendTestGet(tester, settings.TargetAddress+"/service/data/keys/test-group-get/test-group-get-key")
	if getResponse.StatusCode != http.StatusOK {
		tester.Fatalf("group get key status = %d, want %d", getResponse.StatusCode, http.StatusOK)
	}

	result := utilities.ParseValue(tester, getResponse)
	if result.Value == nil {
		tester.Fatal("group get key value is nil")
	}

	var value string
	if issue := json.Unmarshal(result.Value, &value); issue != nil {
		tester.Fatalf("group get key value unmarshal with issue, %v", issue)
	}

	if value != "test group get key value" {
		tester.Fatalf("group get key value = %q, want %q", value, "test group get key value")
	}
}

func TestGroupGetKeyNotFoundFeature(tester *testing.T) {
	utilities.SetupTest(tester)

	getResponse := utilities.SendTestGet(tester, settings.TargetAddress+"/service/data/keys/test-group-not-found/test-group-missing-key")
	if getResponse.StatusCode != http.StatusNotFound {
		tester.Fatalf("group get key not found status = %d, want %d", getResponse.StatusCode, http.StatusNotFound)
	}

	getResponse.Body.Close()
}
