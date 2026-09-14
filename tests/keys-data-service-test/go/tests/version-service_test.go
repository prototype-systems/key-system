package tests

import (
	"net/http"
	"testing"

	"keys-data-service-test/go/settings"
	"keys-data-service-test/go/utilities"
)

func TestVersionServiceFeature(tester *testing.T) {
	utilities.SetupTest(tester)

	response := utilities.SendTestGet(tester, settings.TargetAddress+"/service/data/keys/version")
	if response.StatusCode != http.StatusOK {
		tester.Fatalf("version status = %d, want %d", response.StatusCode, http.StatusOK)
	}

	version := utilities.ParseVersion(tester, response)
	if version.Service != "keys-data-service" {
		tester.Fatalf("version service = %q, want %q", version.Service, "keys-data-service")
	}

	if version.Number == "" {
		tester.Fatal("version number is empty")
	}
}
