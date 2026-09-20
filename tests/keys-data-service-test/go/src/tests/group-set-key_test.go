package tests

import (
	"net/http"
	"testing"

	"keys-data-service-test/settings"
	"keys-data-service-test/utilities"
)

func TestGroupSetKeyFeature(tester *testing.T) {
	utilities.SetupTest(tester)

	body := map[string]any{
		"value": "test group set key value",
	}

	response := utilities.SendTestPost(tester, settings.TargetAddress+"/service/data/keys/test-group/test-group-set-key", body)
	if response.StatusCode != http.StatusOK {
		tester.Fatalf("group set key status = %d, want %d", response.StatusCode, http.StatusOK)
	}

	reference := utilities.ParseReference(tester, response)
	if reference == "" {
		tester.Fatal("group set key reference is empty")
	}
}

func TestGroupSetKeyUpdateFeature(tester *testing.T) {
	utilities.SetupTest(tester)

	firstBody := map[string]any{
		"value": "test group set key initial value",
	}

	firstResponse := utilities.SendTestPost(tester, settings.TargetAddress+"/service/data/keys/test-group-update/test-group-update-key", firstBody)
	if firstResponse.StatusCode != http.StatusOK {
		tester.Fatalf("group set key first set status = %d, want %d", firstResponse.StatusCode, http.StatusOK)
	}

	firstReference := utilities.ParseReference(tester, firstResponse)
	if firstReference == "" {
		tester.Fatal("group set key first reference is empty")
	}

	secondBody := map[string]any{
		"value": "test group set key updated value",
	}

	secondResponse := utilities.SendTestPost(tester, settings.TargetAddress+"/service/data/keys/test-group-update/test-group-update-key", secondBody)
	if secondResponse.StatusCode != http.StatusOK {
		tester.Fatalf("group set key second set status = %d, want %d", secondResponse.StatusCode, http.StatusOK)
	}

	secondReference := utilities.ParseReference(tester, secondResponse)
	if secondReference != firstReference {
		tester.Fatalf("group set key update reference = %q, want %q", secondReference, firstReference)
	}
}
