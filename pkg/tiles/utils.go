package tiles

import (
	"net/http"
	"time"
)

// Create a Client for control over HTTP client settings.
// Client is safe for concurrent use by multiple goroutines
// and for efficiency should only be created once and re-used.
var client = &http.Client{
	Timeout: time.Second * 30,
}

// GetTile sends http.Get request to WMS Server
// and returns response content.
func GetTile() {
	// To do
}
