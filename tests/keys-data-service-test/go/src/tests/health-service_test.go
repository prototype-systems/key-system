package tests

import (
	"net/http"
	"testing"

	"keys-data-service-test/settings"
	"keys-data-service-test/utilities"
)

func TestHealthServiceFeature(tester *testing.T) {
	utilities.SetupTest(tester)

	response := utilities.SendTestGet(tester, settings.TargetAddress+"/service/data/keys/health")
	if response.StatusCode != http.StatusOK {
		tester.Fatalf("health status = %d, want %d", response.StatusCode, http.StatusOK)
	}

	health := utilities.ParseHealth(tester, response)
	if health.Status != "healthy" {
		tester.Fatalf("health status = %q, want %q", health.Status, "healthy")
	}
}
