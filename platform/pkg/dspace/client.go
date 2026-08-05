package dspace

import (
	"context"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"time"

	"go.uber.org/zap"
)

// Client represents a DSpace API client
type Client struct {
	baseURL    string
	apiKey     string
	httpClient *http.Client
	logger     *zap.SugaredLogger
}

// Config holds DSpace configuration
type Config struct {
	URL    string
	APIKey string
}

// NewClient creates a new DSpace API client
func NewClient(cfg Config, logger *zap.SugaredLogger) *Client {
	return &Client{
		baseURL: cfg.URL,
		apiKey:  cfg.APIKey,
		httpClient: &http.Client{
			Timeout: 30 * time.Second,
		},
		logger: logger,
	}
}

// Community represents a DSpace community
type Community struct {
	ID       string `json:"id"`
	Name     string `json:"name"`
	Handle   string `json:"handle"`
	Metadata []struct {
		Key   string   `json:"key"`
		Value []string `json:"value"`
	} `json:"metadata"`
}

// Collection represents a DSpace collection
type Collection struct {
	ID       string `json:"id"`
	Name     string `json:"name"`
	Handle   string `json:"handle"`
	Metadata []struct {
		Key   string   `json:"key"`
		Value []string `json:"value"`
	} `json:"metadata"`
}

// Item represents a DSpace item (resource)
type Item struct {
	ID       string `json:"id"`
	Name     string `json:"name"`
	Handle   string `json:"handle"`
	Metadata []struct {
		Key   string   `json:"key"`
		Value []string `json:"value"`
	} `json:"metadata"`
	Bitstreams []struct {
		ID       string `json:"id"`
		Name     string `json:"name"`
		MimeType string `json:"mimeType"`
		Size     int64  `json:"sizeBytes"`
	} `json:"bitstreams"`
}

// GetCommunities retrieves communities from DSpace
func (c *Client) GetCommunities(ctx context.Context) ([]Community, error) {
	url := fmt.Sprintf("%s/rest/communities", c.baseURL)
	resp, err := c.do(ctx, http.MethodGet, url, nil)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		body, _ := io.ReadAll(resp.Body)
		return nil, fmt.Errorf("dspace error: %d - %s", resp.StatusCode, string(body))
	}

	var communities []Community
	if err := json.NewDecoder(resp.Body).Decode(&communities); err != nil {
		return nil, fmt.Errorf("decoding communities: %w", err)
	}

	c.logger.Infow("Retrieved communities from DSpace", "count", len(communities))
	return communities, nil
}

// GetCollectionsByID retrieves collections within a community
func (c *Client) GetCollectionsByID(ctx context.Context, communityID string) ([]Collection, error) {
	url := fmt.Sprintf("%s/rest/communities/%s/collections", c.baseURL, communityID)
	resp, err := c.do(ctx, http.MethodGet, url, nil)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		body, _ := io.ReadAll(resp.Body)
		return nil, fmt.Errorf("dspace error: %d - %s", resp.StatusCode, string(body))
	}

	var collections []Collection
	if err := json.NewDecoder(resp.Body).Decode(&collections); err != nil {
		return nil, fmt.Errorf("decoding collections: %w", err)
	}

	c.logger.Infow("Retrieved collections from DSpace", "community", communityID, "count", len(collections))
	return collections, nil
}

// GetItemsByCollection retrieves items within a collection
func (c *Client) GetItemsByCollection(ctx context.Context, collectionID string, limit, offset int) ([]Item, error) {
	url := fmt.Sprintf("%s/rest/collections/%s/items?limit=%d&offset=%d", c.baseURL, collectionID, limit, offset)
	resp, err := c.do(ctx, http.MethodGet, url, nil)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		body, _ := io.ReadAll(resp.Body)
		return nil, fmt.Errorf("dspace error: %d - %s", resp.StatusCode, string(body))
	}

	var items []Item
	if err := json.NewDecoder(resp.Body).Decode(&items); err != nil {
		return nil, fmt.Errorf("decoding items: %w", err)
	}

	c.logger.Debugw("Retrieved items from DSpace", "collection", collectionID, "count", len(items))
	return items, nil
}

// GetItemByID retrieves a specific item
func (c *Client) GetItemByID(ctx context.Context, itemID string) (*Item, error) {
	url := fmt.Sprintf("%s/rest/items/%s", c.baseURL, itemID)
	resp, err := c.do(ctx, http.MethodGet, url, nil)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		body, _ := io.ReadAll(resp.Body)
		return nil, fmt.Errorf("dspace error: %d - %s", resp.StatusCode, string(body))
	}

	var item Item
	if err := json.NewDecoder(resp.Body).Decode(&item); err != nil {
		return nil, fmt.Errorf("decoding item: %w", err)
	}

	return &item, nil
}

// GetBitstream retrieves bitstream (file) content
func (c *Client) GetBitstream(ctx context.Context, bitstreamID string) (io.ReadCloser, error) {
	url := fmt.Sprintf("%s/rest/bitstreams/%s/retrieve", c.baseURL, bitstreamID)
	resp, err := c.do(ctx, http.MethodGet, url, nil)
	if err != nil {
		return nil, err
	}

	if resp.StatusCode != http.StatusOK {
		body, _ := io.ReadAll(resp.Body)
		resp.Body.Close()
		return nil, fmt.Errorf("dspace error: %d - %s", resp.StatusCode, string(body))
	}

	return resp.Body, nil
}

// do makes an HTTP request to DSpace
func (c *Client) do(ctx context.Context, method, url string, body io.Reader) (*http.Response, error) {
	req, err := http.NewRequestWithContext(ctx, method, url, body)
	if err != nil {
		return nil, fmt.Errorf("creating request: %w", err)
	}

	req.Header.Set("Accept", "application/json")
	req.Header.Set("Content-Type", "application/json")

	if c.apiKey != "" {
		req.Header.Set("Authorization", fmt.Sprintf("Bearer %s", c.apiKey))
	}

	return c.httpClient.Do(req)
}
