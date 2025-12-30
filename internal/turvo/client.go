package turvo

import (
	"encoding/json"
	"fmt"
	"sync"
	"time"

	"github.com/go-resty/resty/v2"
	"github.com/lwlach/turvo-integration-backend/internal/models"
)

type Config struct {
	BaseURL      string
	APIKey       string // Deprecated: Use ClientName and ClientSecret instead
	ClientName   string
	ClientSecret string
	Username     string
	Password     string
}

type Client struct {
	baseURL      string
	httpClient   *resty.Client
	apiKey       string
	clientName   string
	clientSecret string
	username     string
	password     string

	// Token management
	token       string
	tokenExpiry time.Time
	tokenMutex  sync.RWMutex
}

func NewClient(cfg Config) *Client {
	client := resty.New().
		SetBaseURL(cfg.BaseURL).
		SetTimeout(30*time.Second).
		SetHeader("Content-Type", "application/json").
		SetHeader("Accept", "application/json").
		SetHeader("x-api-key", cfg.APIKey)

	turvoClient := &Client{
		baseURL:      cfg.BaseURL,
		httpClient:   client,
		apiKey:       cfg.APIKey,
		clientName:   cfg.ClientName,
		clientSecret: cfg.ClientSecret,
		username:     cfg.Username,
		password:     cfg.Password,
	}

	return turvoClient
}

// authenticate retrieves an access token from Turvo's authentication API
func (c *Client) authenticate() error {
	if c.clientName == "" || c.clientSecret == "" {
		return fmt.Errorf("clientName and clientSecret are required for authentication")
	}
	if c.username == "" || c.password == "" {
		return fmt.Errorf("username and password are required for authentication")
	}

	authReq := models.TurvoAuthRequest{
		GrantType:    "password",
		ClientID:     c.clientName,
		ClientSecret: c.clientSecret,
		Username:     c.username,
		Password:     c.password,
		Scope:        "read+trust+write",
		Type:         "business",
	}

	var authResp models.TurvoAuthResponse

	// Use a separate client for auth to avoid circular auth issues
	authClient := resty.New().
		SetHeader("Content-Type", "application/json").
		SetBaseURL("https://my-sandbox-publicapi.turvo.com").
		SetTimeout(30*time.Second).
		SetHeader("x-api-key", c.apiKey).
		SetQueryParam("client_id", c.clientName).
		SetQueryParam("client_secret", c.clientSecret)

	resp, err := authClient.R().
		SetBody(authReq).
		SetResult(&authResp).
		Post("/v1/oauth/token")

	if err != nil {
		return fmt.Errorf("failed to authenticate: %w", err)
	}

	if resp.IsError() {
		return fmt.Errorf("authentication failed: %s - %s", resp.Status(), string(resp.Body()))
	}

	if authResp.AccessToken == "" {
		return fmt.Errorf("authentication response missing access token")
	}

	// Store token and calculate expiry
	c.tokenMutex.Lock()
	c.token = authResp.AccessToken
	// Set expiry time (default to 1 hour if not provided, with 5 minute buffer)
	expiresIn := authResp.ExpiresIn
	if expiresIn == 0 {
		expiresIn = 3600 // Default to 1 hour
	}
	c.tokenExpiry = time.Now().Add(time.Duration(expiresIn-300) * time.Second) // 5 min buffer
	c.tokenMutex.Unlock()

	// Update the HTTP client with the new token
	c.httpClient.SetHeader("Authorization", fmt.Sprintf("%s %s", authResp.TokenType, authResp.AccessToken))

	return nil
}

// ensureAuthenticated ensures we have a valid token, refreshing if necessary
func (c *Client) ensureAuthenticated() error {
	c.tokenMutex.RLock()
	tokenValid := c.token != "" && time.Now().Before(c.tokenExpiry)
	c.tokenMutex.RUnlock()

	if tokenValid {
		return nil
	}

	// Need to authenticate
	return c.authenticate()
}

// getAuthenticatedRequest returns a resty request with authentication applied
func (c *Client) getAuthenticatedRequest() (*resty.Request, error) {
	if err := c.ensureAuthenticated(); err != nil {
		return nil, err
	}

	c.tokenMutex.RLock()
	token := c.token
	c.tokenMutex.RUnlock()

	return c.httpClient.R().SetAuthToken(token), nil
}

// GetShipments fetches all shipments from Turvo (alias for ListShipments for backward compatibility)
func (c *Client) GetShipments() ([]models.TurvoShipment, error) {
	return c.ListShipments()
}

// ListShipments fetches all shipments from Turvo
func (c *Client) ListShipments() ([]models.TurvoShipment, error) {
	return c.ListShipmentsWithFilters(models.TurvoShipmentFilters{})
}

// ListShipmentsWithFilters fetches shipments from Turvo with filters and pagination
func (c *Client) ListShipmentsWithFilters(filters models.TurvoShipmentFilters) ([]models.TurvoShipment, error) {
	var response models.TurvoShipmentsListResponse

	req, err := c.getAuthenticatedRequest()
	if err != nil {
		return nil, err
	}

	// Build query parameters
	if filters.Status != "" {
		req.SetQueryParam("status[eq]", filters.Status)
	}
	if filters.CustomerID != "" {
		req.SetQueryParam("customerId[eq]", filters.CustomerID)
	}
	if filters.PickupDateGte != "" {
		req.SetQueryParam("pickupDate[gte]", filters.PickupDateGte)
	}
	if filters.PickupDateLte != "" {
		req.SetQueryParam("pickupDate[lte]", filters.PickupDateLte)
	}
	if filters.Start > 0 {
		req.SetQueryParam("start", fmt.Sprintf("%d", filters.Start))
	}
	if filters.PageSize > 0 {
		req.SetQueryParam("pageSize", fmt.Sprintf("%d", filters.PageSize))
	}

	resp, err := req.
		SetResult(&response).
		Get("/v1/shipments/list")

	if err != nil {
		return nil, fmt.Errorf("failed to list shipments: %w", err)
	}

	if resp.IsError() {
		// If unauthorized, try to re-authenticate once
		if resp.StatusCode() == 401 {
			if authErr := c.authenticate(); authErr != nil {
				return nil, fmt.Errorf("authentication failed: %w", authErr)
			}
			// Retry the request
			req, err := c.getAuthenticatedRequest()
			if err != nil {
				return nil, err
			}
			// Re-apply query parameters
			if filters.Status != "" {
				req.SetQueryParam("status[eq]", filters.Status)
			}
			if filters.CustomerID != "" {
				req.SetQueryParam("customerId[eq]", filters.CustomerID)
			}
			if filters.PickupDateGte != "" {
				req.SetQueryParam("pickupDate[gte]", filters.PickupDateGte)
			}
			if filters.PickupDateLte != "" {
				req.SetQueryParam("pickupDate[lte]", filters.PickupDateLte)
			}
			if filters.Start > 0 {
				req.SetQueryParam("start", fmt.Sprintf("%d", filters.Start))
			}
			if filters.PageSize > 0 {
				req.SetQueryParam("pageSize", fmt.Sprintf("%d", filters.PageSize))
			}
			resp, err = req.SetResult(&response).Get("/v1/shipments/list")
			if err != nil {
				return nil, fmt.Errorf("failed to list shipments after re-auth: %w", err)
			}
			if resp.IsError() {
				return nil, fmt.Errorf("turvo API error: %s - %s", resp.Status(), string(resp.Body()))
			}
		} else {
			return nil, fmt.Errorf("turvo API error: %s - %s", resp.Status(), string(resp.Body()))
		}
	}

	return response.Details.Shipments, nil
}

// ListShipmentsWithFiltersAndPagination fetches shipments with filters and returns pagination info
func (c *Client) ListShipmentsWithFiltersAndPagination(filters models.TurvoShipmentFilters) ([]models.TurvoShipment, models.TurvoPagination, error) {
	var response models.TurvoShipmentsListResponse

	req, err := c.getAuthenticatedRequest()
	if err != nil {
		return nil, models.TurvoPagination{}, err
	}

	// Build query parameters
	if filters.Status != "" {
		req.SetQueryParam("status[eq]", filters.Status)
	}
	if filters.CustomerID != "" {
		req.SetQueryParam("customerId[eq]", filters.CustomerID)
	}
	if filters.PickupDateGte != "" {
		req.SetQueryParam("pickupDate[gte]", filters.PickupDateGte)
	}
	if filters.PickupDateLte != "" {
		req.SetQueryParam("pickupDate[lte]", filters.PickupDateLte)
	}
	if filters.Start > 0 {
		req.SetQueryParam("start", fmt.Sprintf("%d", filters.Start))
	}
	if filters.PageSize > 0 {
		req.SetQueryParam("pageSize", fmt.Sprintf("%d", filters.PageSize))
	}

	resp, err := req.
		SetResult(&response).
		Get("/v1/shipments/list")

	if err != nil {
		return nil, models.TurvoPagination{}, fmt.Errorf("failed to list shipments: %w", err)
	}

	if resp.IsError() {
		// If unauthorized, try to re-authenticate once
		if resp.StatusCode() == 401 {
			if authErr := c.authenticate(); authErr != nil {
				return nil, models.TurvoPagination{}, fmt.Errorf("authentication failed: %w", authErr)
			}
			// Retry the request
			req, err := c.getAuthenticatedRequest()
			if err != nil {
				return nil, models.TurvoPagination{}, err
			}
			// Re-apply query parameters
			if filters.Status != "" {
				req.SetQueryParam("status[eq]", filters.Status)
			}
			if filters.CustomerID != "" {
				req.SetQueryParam("customerId[eq]", filters.CustomerID)
			}
			if filters.PickupDateGte != "" {
				req.SetQueryParam("pickupDate[gte]", filters.PickupDateGte)
			}
			if filters.PickupDateLte != "" {
				req.SetQueryParam("pickupDate[lte]", filters.PickupDateLte)
			}
			if filters.Start > 0 {
				req.SetQueryParam("start", fmt.Sprintf("%d", filters.Start))
			}
			if filters.PageSize > 0 {
				req.SetQueryParam("pageSize", fmt.Sprintf("%d", filters.PageSize))
			}
			resp, err = req.SetResult(&response).Get("/v1/shipments/list")
			if err != nil {
				return nil, models.TurvoPagination{}, fmt.Errorf("failed to list shipments after re-auth: %w", err)
			}
			if resp.IsError() {
				return nil, models.TurvoPagination{}, fmt.Errorf("turvo API error: %s - %s", resp.Status(), string(resp.Body()))
			}
		} else {
			return nil, models.TurvoPagination{}, fmt.Errorf("turvo API error: %s - %s", resp.Status(), string(resp.Body()))
		}
	}

	return response.Details.Shipments, response.Details.Pagination, nil
}

// CreateShipment creates a new shipment in Turvo
func (c *Client) CreateShipment(shipment *models.TurvoShipmentCreate) (*models.TurvoShipmentCreateResponse, error) {
	var response models.TurvoShipmentCreateResponse
	var errorResponse models.TurvoShipmentCreateErrorResponse

	req, err := c.getAuthenticatedRequest()
	if err != nil {
		return nil, err
	}

	resp, err := req.
		SetBody(shipment).
		SetResult(&response).
		SetError(&errorResponse).
		Post("/v1/shipments")

	if err != nil {
		return nil, fmt.Errorf("failed to create shipment: %w", err)
	}

	if resp.IsError() {
		// If unauthorized, try to re-authenticate once
		if resp.StatusCode() == 401 {
			if authErr := c.authenticate(); authErr != nil {
				return nil, fmt.Errorf("authentication failed: %w", authErr)
			}
			// Retry the request
			req, err := c.getAuthenticatedRequest()
			if err != nil {
				return nil, err
			}
			resp, err = req.SetBody(shipment).SetResult(&response).Post("/v1/shipments")
			if err != nil {
				return nil, fmt.Errorf("failed to create shipment after re-auth: %w", err)
			}
			if resp.IsError() {
				return nil, fmt.Errorf("turvo API error: %s - %s", errorResponse.Status, errorResponse.Details.ErrorMessage)
			}
		} else {
			return nil, fmt.Errorf("turvo API error: %s - %s", resp.Status(), string(resp.Body()))
		}
	}

	if response.Status != "SUCCESS" {
		err := json.Unmarshal(resp.Body(), &errorResponse)
		if err != nil {
			return nil, fmt.Errorf("failed to unmarshal error response: %w", err)
		}
		return nil, fmt.Errorf("turvo API error: %s - %s", errorResponse.Status, errorResponse.Details.ErrorMessage)
	}

	return &response, nil
}

// GetShipment fetches a single shipment by ID from Turvo
func (c *Client) GetShipment(shipmentID int) (*models.TurvoShipmentCreateDetails, error) {
	var response models.TurvoShipmentResponse

	req, err := c.getAuthenticatedRequest()
	if err != nil {
		return nil, err
	}

	resp, err := req.
		SetResult(&response).
		Get(fmt.Sprintf("/v1/shipments/%d", shipmentID))

	if err != nil {
		return nil, fmt.Errorf("failed to get shipment: %w", err)
	}

	if resp.IsError() {
		// If unauthorized, try to re-authenticate once
		if resp.StatusCode() == 401 {
			if authErr := c.authenticate(); authErr != nil {
				return nil, fmt.Errorf("authentication failed: %w", authErr)
			}
			// Retry the request
			req, err := c.getAuthenticatedRequest()
			if err != nil {
				return nil, err
			}
			resp, err = req.SetResult(&response).Get(fmt.Sprintf("/v1/shipments/%d", shipmentID))
			if err != nil {
				return nil, fmt.Errorf("failed to get shipment after re-auth: %w", err)
			}
			if resp.IsError() {
				return nil, fmt.Errorf("turvo API error: %s - %s", resp.Status(), string(resp.Body()))
			}
		} else {
			return nil, fmt.Errorf("turvo API error: %s - %s", resp.Status(), string(resp.Body()))
		}
	}

	if response.Status != "SUCCESS" {
		return nil, fmt.Errorf("turvo API error: %s", response.Status)
	}

	return &response.Details, nil
}
