package testutils

import (
	"net/http"
	"net/http/httptest"
	"net/url"
)

// SetupServer setups a test server for mocking API Calls
func SetupServer() (mux *http.ServeMux, hostname, port string, teardown func()) {
	mux = http.NewServeMux()

	// server is a test HTTP server used to provide mock API responses
	server := httptest.NewServer(mux)
	url, _ := url.Parse(server.URL)
	return mux, url.Path, url.Port(), server.Close
}
