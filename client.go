package gopaapi5

import (
	"encoding/json"
	"errors"
	"io/ioutil"
	"net/http"

	"github.com/utekaravinash/gopaapi5/api"
)

var (
	ErrEmptyAccessKey    = errors.New("Empty access key")
	ErrEmptySecretKey    = errors.New("Empty secret key")
	ErrEmptyAssociateTag = errors.New("Empty associate tag")
	ErrInvalidLocale     = errors.New("Invalid locale")
)

// Client stores AccessKey, SecretKey, and, AssociateTag; and exposes GetBrowseNodes, GetItems, GetVariations, and SearchItems operations.
type Client struct {
	AccessKey    string
	SecretKey    string
	AssociateTag string
	Locale       api.Locale
	partnerType  string
	service      string
	host         string
	region       string
	marketplace  string
	httpClient   *http.Client
	testing      bool
}

// NewClient accepts Access Key, Secrete Key, Associate Tag, Locale and returns a new client
func NewClient(accessKey, secretKey, associateTag string, locale api.Locale) (*Client, error) {

	if accessKey == "" {
		return nil, ErrEmptyAccessKey
	}

	if secretKey == "" {
		return nil, ErrEmptySecretKey
	}

	if associateTag == "" {
		return nil, ErrEmptyAssociateTag
	}

	if !locale.IsValid() {
		return nil, ErrInvalidLocale
	}

	client := &Client{
		AccessKey:    accessKey,
		SecretKey:    secretKey,
		AssociateTag: associateTag,
		partnerType:  "Associates",
		service:      "ProductAdvertisingAPIv1",
		host:         locale.Host(),
		region:       locale.Region(),
		marketplace:  locale.Marketplace(),
		httpClient:   &http.Client{},
		testing:      false,
	}

	return client, nil
}

// send sends a http request to Amazon Product Advertising service and returns response or error
func (c *Client) send(req *request, v interface{}) error {

	// Construct http request
	err := req.build()
	if err != nil {
		return err
	}

	// Sign http request by adding Authorization header
	err = req.sign()
	if err != nil {
		return err
	}

	// Send http request and receive response or return error if any
	resp, err := c.httpClient.Do(req.httpReq)
	if err != nil {
		return err
	}
	defer resp.Body.Close()

	// Read response body
	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return err
	}

	// Unmarshal response body to operation specific response struct
	err = json.Unmarshal(body, v)
	if err != nil {
		return err
	}

	return nil
}
