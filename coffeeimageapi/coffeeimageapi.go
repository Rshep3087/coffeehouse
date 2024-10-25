// coffeeimageapi provides a client that satisfies the ImageAPI interface
// for the coffeehouse application.
package coffeeimageapi

import (
	"context"
	"encoding/json"
	"net/http"
	"time"

	"github.com/rshep3087/coffeehouse/web"
	"go.opentelemetry.io/contrib/instrumentation/net/http/otelhttp"
)

var _ web.CoffeeImageProvider = (*Client)(nil)

// Client is a client for the coffee image API
type Client struct {
	// URL is the base URL for the coffee image API
	URL string
	// client is the HTTP client to use for requests
	client *http.Client
}

// NewClient creates a new Client with the given URL
func NewClient(url string) *Client {
	client := &http.Client{
		Timeout:   5 * time.Second,
		Transport: otelhttp.NewTransport(http.DefaultTransport),
	}

	return &Client{URL: url, client: client}
}

type randomResponse struct {
	// File is the URL of the random image
	File string `json:"file"`
}

// GetImageURL returns a URL for a random coffee image
func (c *Client) GetImageURL(ctx context.Context) string {

	req, err := http.NewRequestWithContext(ctx, http.MethodGet, c.URL+"/random.json", nil)
	if err != nil {
		return ""
	}

	resp, err := c.client.Do(req)
	if err != nil {
		return ""
	}
	defer resp.Body.Close()

	var r randomResponse
	if err := json.NewDecoder(resp.Body).Decode(&r); err != nil {
		return ""
	}

	return r.File
}
