package kraken

import (
	"net/http"
	"net/http/httptest"
	"testing"
)

func newMockClient(clientID string, mockHandler func(http.ResponseWriter, *http.Request)) *Client {
	mc := &Client{}
	mc.clientID = clientID
	mc.httpClient = &mockHTTPClient{mockHandler}

	return mc
}

func newMockHandler(statusCode int, json string) func(http.ResponseWriter, *http.Request) {
	return func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(statusCode)
		w.Write([]byte(json))
	}
}

type mockHTTPClient struct {
	mockHandler func(http.ResponseWriter, *http.Request)
}

func (mtc *mockHTTPClient) Do(req *http.Request) (*http.Response, error) {
	rr := httptest.NewRecorder()
	handler := http.HandlerFunc(mtc.mockHandler)
	handler.ServeHTTP(rr, req)

	return rr.Result(), nil
}

func TestNewClient(t *testing.T) {
	t.Parallel()

	testCases := []struct {
		clientID   string
		httpClient HTTPClient
		shouldFail bool
	}{
		{"client-id", nil, false},
		{"client-id", &http.Client{}, false},
		{"", nil, true},
	}

	for _, testCase := range testCases {
		c, err := NewClient(testCase.clientID, testCase.httpClient)
		if err != nil {
			if testCase.shouldFail {
				continue
			}
			t.Error(err)
		}

		if c.clientID != testCase.clientID {
			t.Errorf("expected clientID to be %s, got %s", testCase.clientID, c.clientID)
		}
	}
}
