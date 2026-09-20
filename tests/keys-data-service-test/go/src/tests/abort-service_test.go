package tests

import (
	"net/http"
	"testing"

	"keys-data-service-test/settings"
	"keys-data-service-test/utilities"
)

func TestAbortServiceFeature(tester *testing.T) {
	test := utilities.SetupTest(tester)

	abortResponse := utilities.SendTestPost(tester, settings.TargetAddress+"/service/data/keys/abort")
	if abortResponse.StatusCode != http.StatusOK {
		tester.Fatalf("abort status = %d, want %d", abortResponse.StatusCode, http.StatusOK)
	}

	operation := utilities.ParseOperation(tester, abortResponse)
	if operation.Procedure != "abort" || operation.Status != "initiated" {
		tester.Fatalf("abort operation = %#v, want procedure abort and status initiated", operation)
	}

	utilities.CheckTestHealthy(tester, test)

	killResponse := utilities.SendTestPost(tester, settings.TargetAddress+"/service/data/keys/kill")
	if killResponse.StatusCode != http.StatusOK {
		tester.Fatalf("kill status = %d, want %d", killResponse.StatusCode, http.StatusOK)
	}

	utilities.CheckTestDone(tester, test)
}
