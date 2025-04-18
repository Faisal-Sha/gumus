package crawler

import (
	"encoding/json"
	"fmt"
	"net/http"
	"time"
	"trendyol-tracker/pkg/models"
)

type Crawler struct {
	client     *http.Client
	baseURL    string
	userAgent  string
	rateLimit  time.Duration
	lastAccess time.Time
}

func NewCrawler() *Crawler {
	return &Crawler{
		client: &http.Client{
			Timeout: time.Second * 30,
			Transport: &http.Transport{
				MaxIdleConns:        100,
				IdleConnTimeout:     90 * time.Second,
				DisableCompression:  true,
				DisableKeepAlives:   false,
			},
		},
		baseURL:   "https://www.trendyol.com",
		userAgent: "Mozilla/5.0 (Windows NT 10.0; Win64; x64) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/91.0.4472.124 Safari/537.36",
		rateLimit: time.Millisecond * 500, // 500ms between requests
	}
}

func (c *Crawler) GetCategories() ([]string, error) {
	// Respect rate limiting
	c.waitForRateLimit()

	req, err := http.NewRequest("GET", c.baseURL, nil)
	if err != nil {
		return nil, fmt.Errorf("failed to create request: %w", err)
	}

	req.Header.Set("User-Agent", c.userAgent)
	
	resp, err := c.client.Do(req)
	if err != nil {
		return nil, fmt.Errorf("failed to fetch categories: %w", err)
	}
	defer resp.Body.Close()

	// Parse HTML and extract category links
	// TODO: Implement HTML parsing logic
	categories := []string{
		"/kadin",
		"/erkek",
		"/cocuk",
		"/ev-yasam",
		"/supermarket",
		"/kozmetik",
		"/ayakkabi-canta",
		"/elektronik",
		"/spor-outdoor",
	}

	return categories, nil
}

func (c *Crawler) GetProductsFromCategory(categoryPath string, page int) ([]models.Product, error) {
	c.waitForRateLimit()

	apiURL := fmt.Sprintf("%s/api/products%s?page=%d", c.baseURL, categoryPath, page)
	req, err := http.NewRequest("GET", apiURL, nil)
	if err != nil {
		return nil, fmt.Errorf("failed to create request: %w", err)
	}

	req.Header.Set("User-Agent", c.userAgent)
	req.Header.Set("Accept", "application/json")

	resp, err := c.client.Do(req)
	if err != nil {
		return nil, fmt.Errorf("failed to fetch products: %w", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode == http.StatusTooManyRequests {
		// Handle rate limiting
		time.Sleep(time.Second * 30)
		return c.GetProductsFromCategory(categoryPath, page)
	}

	var products []models.Product
	if err := json.NewDecoder(resp.Body).Decode(&products); err != nil {
		return nil, fmt.Errorf("failed to decode products: %w", err)
	}

	return products, nil
}

func (c *Crawler) GetProductDetails(productID string) (*models.Product, error) {
	c.waitForRateLimit()

	apiURL := fmt.Sprintf("%s/api/products/%s/detail", c.baseURL, productID)
	req, err := http.NewRequest("GET", apiURL, nil)
	if err != nil {
		return nil, fmt.Errorf("failed to create request: %w", err)
	}

	req.Header.Set("User-Agent", c.userAgent)
	req.Header.Set("Accept", "application/json")

	resp, err := c.client.Do(req)
	if err != nil {
		return nil, fmt.Errorf("failed to fetch product details: %w", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode == http.StatusTooManyRequests {
		// Handle rate limiting
		time.Sleep(time.Second * 30)
		return c.GetProductDetails(productID)
	}

	var product models.Product
	if err := json.NewDecoder(resp.Body).Decode(&product); err != nil {
		return nil, fmt.Errorf("failed to decode product details: %w", err)
	}

	return &product, nil
}

func (c *Crawler) waitForRateLimit() {
	elapsed := time.Since(c.lastAccess)
	if elapsed < c.rateLimit {
		time.Sleep(c.rateLimit - elapsed)
	}
	c.lastAccess = time.Now()
}

// Helper function to handle retries with exponential backoff
func (c *Crawler) retryRequest(req *http.Request, maxRetries int) (*http.Response, error) {
	var lastErr error
	for i := 0; i < maxRetries; i++ {
		resp, err := c.client.Do(req)
		if err != nil {
			lastErr = err
			waitTime := time.Second * time.Duration(1<<uint(i))
			time.Sleep(waitTime)
			continue
		}
		return resp, nil
	}
	return nil, fmt.Errorf("max retries exceeded: %w", lastErr)
}
