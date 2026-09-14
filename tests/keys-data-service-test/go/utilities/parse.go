package utilities

import (
	"encoding/json"
	"keys-data-service-test/go/models"
	"net/http"
	"testing"
)

func ParseOperation(tester *testing.T, response *http.Response) models.Operation {
	tester.Helper()

	defer response.Body.Close()

	var result models.OperationResponse
	if issue := json.NewDecoder(response.Body).Decode(&result); issue != nil {
		tester.Fatalf("parse operation with issue, %v", issue)
	}

	return result.Operation
}

func ParseHealth(tester *testing.T, response *http.Response) models.Health {
	tester.Helper()

	defer response.Body.Close()

	var result models.HealthResponse
	if issue := json.NewDecoder(response.Body).Decode(&result); issue != nil {
		tester.Fatalf("parse health with issue, %v", issue)
	}

	return result.Health
}

func ParseVersion(tester *testing.T, response *http.Response) models.Version {
	tester.Helper()

	defer response.Body.Close()

	var result models.VersionResponse
	if issue := json.NewDecoder(response.Body).Decode(&result); issue != nil {
		tester.Fatalf("parse version with issue, %v", issue)
	}

	return result.Version
}

func ParseReference(tester *testing.T, response *http.Response) string {
	tester.Helper()

	defer response.Body.Close()

	var result models.ReferenceResponse
	if issue := json.NewDecoder(response.Body).Decode(&result); issue != nil {
		tester.Fatalf("parse reference with issue, %v", issue)
	}

	return result.Reference
}

func ParseKey(tester *testing.T, response *http.Response) models.Key {
	tester.Helper()

	defer response.Body.Close()

	var result models.KeyResponse
	if issue := json.NewDecoder(response.Body).Decode(&result); issue != nil {
		tester.Fatalf("parse key with issue, %v", issue)
	}

	return result.Key
}

func ParseKeys(tester *testing.T, response *http.Response) models.KeysResponse {
	tester.Helper()

	defer response.Body.Close()

	var result models.KeysResponse
	if issue := json.NewDecoder(response.Body).Decode(&result); issue != nil {
		tester.Fatalf("parse keys with issue, %v", issue)
	}

	return result
}

func ParseCount(tester *testing.T, response *http.Response) int {
	tester.Helper()

	defer response.Body.Close()

	var result models.CountResponse
	if issue := json.NewDecoder(response.Body).Decode(&result); issue != nil {
		tester.Fatalf("parse count with issue, %v", issue)
	}

	return result.Count
}

func ParseValue(tester *testing.T, response *http.Response) models.ValueResponse {
	tester.Helper()

	defer response.Body.Close()

	var result models.ValueResponse
	if issue := json.NewDecoder(response.Body).Decode(&result); issue != nil {
		tester.Fatalf("parse value with issue, %v", issue)
	}

	return result
}
