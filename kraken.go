package kraken

import (
	"bytes"
	"encoding/json"
	"errors"
	"fmt"
	"net/http"
	"reflect"
	"strings"
)

const (
	methodGet = "GET"
)

var (
	basePath = "https://api.twitch.tv/kraken"
	queryTag = "query"
)

// HTTPClient ...
type HTTPClient interface {
	Do(req *http.Request) (*http.Response, error)
}

// Client ...
type Client struct {
	clientID    string
	accessToken string
	userAgent   string
	httpClient  HTTPClient
}

// ResponseCommon ...
type ResponseCommon struct {
	StatusCode   int
	Error        string `json:"error"`
	ErrorStatus  int    `json:"status"`
	ErrorMessage string `json:"message"`
}

// Response ...
type Response struct {
	ResponseCommon
	Data interface{}
}

// NewClient ... It is concurrecy safe.
func NewClient(clientID string, httpClient HTTPClient) (*Client, error) {
	c := &Client{}
	c.clientID = clientID

	if c.clientID == "" {
		return nil, errors.New("clientID cannot be an empty string")
	}

	if httpClient != nil {
		c.httpClient = httpClient
	} else {
		c.httpClient = http.DefaultClient
	}

	return c, nil
}

func (c *Client) get(path string, data, params interface{}) (*Response, error) {
	resp := &Response{
		Data: data,
	}

	req, err := c.newRequest(methodGet, path, params)
	if err != nil {
		return nil, err
	}

	err = c.doRequest(req, resp)
	if err != nil {
		return nil, err
	}

	return resp, nil
}

func buildQueryString(req *http.Request, v interface{}) (string, error) {
	isNil, err := isZero(v)
	if err != nil {
		return "", err
	}

	if isNil {
		return "", nil
	}

	query := req.URL.Query()
	t := reflect.TypeOf(v).Elem()

	for i := 0; i < t.NumField(); i++ {
		field := t.Field(i)
		tag := field.Tag.Get(queryTag)
		defaultValue := ""

		// Get the default value from the struct tag
		if strings.Contains(tag, ",") {
			tagSlice := strings.Split(tag, ",")

			tag = tagSlice[0]
			defaultValue = tagSlice[1]
		}

		// Get the value passed in by the user
		value := fmt.Sprintf("%v", reflect.Indirect(reflect.ValueOf(v)).FieldByName(field.Name))

		// If no value was set by the user, use the default
		// value specified in the struct tag.
		if value == "" || value == "0" {
			if defaultValue == "" {
				continue
			}

			value = defaultValue
		}

		query.Add(tag, value)
	}

	return query.Encode(), nil
}

func isZero(v interface{}) (bool, error) {
	t := reflect.TypeOf(v)
	if !t.Comparable() {
		return false, fmt.Errorf("type is not comparable: %v", t)
	}
	return v == reflect.Zero(t).Interface(), nil
}

func (c *Client) doRequest(req *http.Request, resp *Response) error {
	c.setRequestHeaders(req)

	response, err := c.httpClient.Do(req)
	if err != nil {
		return fmt.Errorf("Failed to execute API request: %s", err.Error())
	}
	defer response.Body.Close()

	setResponseStatusCode(resp, "StatusCode", response.StatusCode)

	// Only attempt to decode the response if JSON was returned
	if resp.StatusCode < 500 {
		decoder := json.NewDecoder(response.Body)
		if resp.StatusCode < 400 {
			// Successful request
			err = decoder.Decode(resp.Data)

		} else {
			// Failed request
			err = decoder.Decode(resp)
		}

		if err != nil {
			return fmt.Errorf("Failed to decode API response: %s", err.Error())
		}
	}

	return nil
}

func (c *Client) setRequestHeaders(req *http.Request) {
	req.Header.Set("Accept", "application/vnd.twitchtv.v5+json")
	req.Header.Set("Client-ID", c.clientID)

	if c.accessToken != "" {
		req.Header.Set("Authorization", "OAuth "+c.accessToken)
	}
	if c.userAgent != "" {
		req.Header.Set("User-Agent", c.userAgent)
	}
}

func setResponseStatusCode(v interface{}, fieldName string, code int) {
	s := reflect.ValueOf(v).Elem()
	field := s.FieldByName(fieldName)
	field.SetInt(int64(code))
}

// SetAccessToken ...
func (c *Client) SetAccessToken(accessToken string) {
	c.accessToken = accessToken
}

// SetUserAgent ...
func (c *Client) SetUserAgent(userAgent string) {
	c.userAgent = userAgent
}

func (c *Client) newRequest(method, path string, params interface{}) (*http.Request, error) {
	url := concatString([]string{basePath, path})

	req, err := http.NewRequest(method, url, nil)
	if err != nil {
		return nil, err
	}

	if method == methodGet && params != nil {
		query, err := buildQueryString(req, params)
		if err != nil {
			return nil, err
		}

		req.URL.RawQuery = query
	}

	return req, nil
}

// concatString concatenates each of the strings provided by
// strs in the order they are presented. You may also pass in
// an optional delimiter to be appended along with the strings.
func concatString(strs []string, delimiter ...string) string {
	var buffer bytes.Buffer
	appendDelimiter := len(delimiter) > 0

	for _, str := range strs {
		s := str
		if appendDelimiter {
			s = concatString([]string{s, delimiter[0]})
		}
		buffer.Write([]byte(s))
	}

	if appendDelimiter {
		return strings.TrimSuffix(buffer.String(), delimiter[0])
	}

	return buffer.String()
}
