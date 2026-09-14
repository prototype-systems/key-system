package tests

import (
	"net/http"
	"testing"

	"keys-data-service-test/go/settings"
	"keys-data-service-test/go/utilities"
)

func TestPersistFeature(tester *testing.T) {
	utilities.SetupTest(tester)

	response := utilities.SendTestPost(tester, settings.TargetAddress+"/service/data/keys/persist")
	if response.StatusCode != http.StatusOK {
		tester.Fatalf("persist status = %d, want %d", response.StatusCode, http.StatusOK)
	}

	operation := utilities.ParseOperation(tester, response)
	if operation.Procedure != "persist" || operation.Status != "completed" {
		tester.Fatalf("persist operation = %#v, want procedure persist and status completed", operation)
	}
}
