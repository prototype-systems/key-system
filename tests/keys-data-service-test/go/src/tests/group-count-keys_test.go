package tests

import (
	"net/http"
	"testing"

	"keys-data-service-test/settings"
	"keys-data-service-test/utilities"
)

func TestGroupCountKeysFeature(tester *testing.T) {
	utilities.SetupTest(tester)

	firstBody := map[string]any{
		"value": "test group count first value",
	}

	firstResponse := utilities.SendTestPost(tester, settings.TargetAddress+"/service/data/keys/test-group-count/test-group-count-first", firstBody)
	if firstResponse.StatusCode != http.StatusOK {
		tester.Fatalf("group count keys setup first set status = %d, want %d", firstResponse.StatusCode, http.StatusOK)
	}

	firstResponse.Body.Close()

	secondBody := map[string]any{
		"value": "test group count second value",
	}

	secondResponse := utilities.SendTestPost(tester, settings.TargetAddress+"/service/data/keys/test-group-count/test-group-count-second", secondBody)
	if secondResponse.StatusCode != http.StatusOK {
		tester.Fatalf("group count keys setup second set status = %d, want %d", secondResponse.StatusCode, http.StatusOK)
	}

	secondResponse.Body.Close()

	countResponse := utilities.SendTestGet(tester, settings.TargetAddress+"/service/data/keys/test-group-count/count")
	if countResponse.StatusCode != http.StatusOK {
		tester.Fatalf("group count keys status = %d, want %d", countResponse.StatusCode, http.StatusOK)
	}

	count := utilities.ParseCount(tester, countResponse)
	if count < 2 {
		tester.Fatalf("group count keys count = %d, want at least 2", count)
	}
}

func TestGroupCountKeysEmptyFeature(tester *testing.T) {
	utilities.SetupTest(tester)

	countResponse := utilities.SendTestGet(tester, settings.TargetAddress+"/service/data/keys/test-group-count-empty/count")
	if countResponse.StatusCode != http.StatusOK {
		tester.Fatalf("group count keys empty status = %d, want %d", countResponse.StatusCode, http.StatusOK)
	}

	count := utilities.ParseCount(tester, countResponse)
	if count != 0 {
		tester.Fatalf("group count keys empty count = %d, want 0", count)
	}
}
