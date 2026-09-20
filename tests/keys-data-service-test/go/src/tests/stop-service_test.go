package tests

import (
	"net/http"
	"testing"

	"keys-data-service-test/settings"
	"keys-data-service-test/utilities"
)

func TestStopServiceFeature(tester *testing.T) {
	test := utilities.SetupTest(tester)

	stopResponse := utilities.SendTestPost(tester, settings.TargetAddress+"/service/data/keys/stop")
	if stopResponse.StatusCode != http.StatusOK {
		tester.Fatalf("stop status = %d, want %d", stopResponse.StatusCode, http.StatusOK)
	}

	operation := utilities.ParseOperation(tester, stopResponse)
	if operation.Procedure != "stop" || operation.Status != "initiated" {
		tester.Fatalf("stop operation = %#v, want procedure stop and status initiated", operation)
	}

	utilities.CheckTestHealthy(tester, test)
}
