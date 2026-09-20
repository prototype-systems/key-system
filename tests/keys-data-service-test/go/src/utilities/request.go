package utilities

import (
	"bytes"
	"encoding/json"
	"net/http"
)

func SendPost(address string, body ...any) (*http.Response, error) {
	var request *http.Request
	var issue error

	if len(body) > 0 && body[0] != nil {
		var data []byte

		data, issue = json.Marshal(body[0])
		if issue != nil {
			return nil, issue
		}

		request, issue = http.NewRequest(http.MethodPost, address, bytes.NewReader(data))
		if issue != nil {
			return nil, issue
		}

		request.Header.Set("Content-Type", "application/json")
	} else {
		request, issue = http.NewRequest(http.MethodPost, address, nil)
		if issue != nil {
			return nil, issue
		}
	}

	response, issue := http.DefaultClient.Do(request)
	if issue != nil {
		return nil, issue
	}

	return response, nil
}

func SendDelete(address string) (*http.Response, error) {
	request, issue := http.NewRequest(http.MethodDelete, address, nil)
	if issue != nil {
		return nil, issue
	}

	response, issue := http.DefaultClient.Do(request)
	if issue != nil {
		return nil, issue
	}

	return response, nil
}

func SendGet(address string) (*http.Response, error) {
	request, issue := http.NewRequest(http.MethodGet, address, nil)
	if issue != nil {
		return nil, issue
	}

	response, issue := http.DefaultClient.Do(request)
	if issue != nil {
		return nil, issue
	}

	return response, nil
}
