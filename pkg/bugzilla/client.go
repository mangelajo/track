package bugzilla

import (
	"fmt"
	"io"
	"net/http"
	"net/http/cookiejar"
	"reflect"
	"time"
)

const httpTimeout int = 60
const userAgent string = "bugzilla go client"

// newHTTPClient creates HTTP client for HTTP based endpoints
func newHTTPClient() (*http.Client, error) {
	timeout := time.Duration(time.Duration(httpTimeout) * time.Second)
	cookieJar, err := cookiejar.New(nil)
	if err != nil {
		return nil, err
	}

	client := http.Client{
		Jar:     cookieJar,
		Timeout: timeout,
	}
	return &client, nil
}

// newHTTPRequest creates HTTP request
func newHTTPRequest(method string, urlStr string, body io.Reader) (*http.Request, error) {
	req, err := http.NewRequest(method, urlStr, body)
	if err != nil {
		return nil, err
	}
	req.Header.Set("User-Agent", userAgent)
	req.Header.Set("Accept-Charset", "utf-8")
	//todo: support gzip
	//req.Header.Set("Accept-Encoding", "gzip, deflate")
	req.Header.Set("Pragma", "no-cache")
	req.Header.Set("Cache-Control", "no-cache")
	return req, nil
}

// Client bugzilla client
type Client struct {
	bugzillaAddress  string
	bugzillaLogin    string
	bugzillaPassword string
	bugzillaToken    string
	cgi              *bugzillaCGIClient
	json             *bugzillaJSONRPCClient
}

// NewClient creates bugzilla Client instance
func NewClient(bugzillaAddress string, bugzillaLogin string, bugzillaPassword string) (client *Client, err error) {
	httpClient, err := newHTTPClient()
	if err != nil {
		return nil, err
	}
	cgiClient, err := newCGIClient(bugzillaAddress, httpClient)
	if err != nil {
		return nil, err
	}
	jsonClient, err := newJSONRPCClient(bugzillaAddress, httpClient)
	if err != nil {
		return nil, err
	}
	err = cgiClient.login(bugzillaLogin, bugzillaPassword)
	if err != nil {
		return nil, err
	}
	token, err := jsonClient.login(bugzillaLogin, bugzillaPassword)
	if err != nil {
		return nil, err
	}

	client = &Client{
		bugzillaAddress:  bugzillaAddress,
		bugzillaLogin:    bugzillaLogin,
		bugzillaPassword: bugzillaPassword,
		bugzillaToken:    token,
		cgi:              cgiClient,
		json:             jsonClient,
	}

	return client, nil
}

type BugListQuery struct {
	Limit int
	Offset int
	Order string
	Classification string
	Product string
	Component string
	BugStatus []string
	WhiteBoard string
	AssignedTo string

}

// BugList list of last changed bugs
func (client *Client) BugList(query *BugListQuery) ([]Bug, error) {
	return client.cgi.bugList(query)
}

// BugzillaVersion returns Bugzilla version
func (client *Client) BugzillaVersion() (version string, err error) {
	return client.json.bugzillaVersion()
}

// BugInfo returns information about single bugzilla ticket
func (client *Client) BugInfo(id int) (bugInfo map[string]interface{}, err error) {
	bugsInfo, err := client.BugsInfo([]int{id})
	if err != nil {
		return nil, err
	}
	if len(bugsInfo) != 1 {
		return nil, fmt.Errorf("invalid length of array, expected = 1, got = %v", len(bugsInfo))
	}
	return bugsInfo[0], nil
}

func (client *Client) ShowBug(id int)  (bug *Cbug, err error) {
	return client.cgi.bugInfo(id)
}

// BugsInfo returns information about selected bugzilla tickets
func (client *Client) BugsInfo(idList []int) (bugInfo []map[string]interface{}, err error) {
	bugsInfo, err := client.json.bugsInfo(idList, client.bugzillaToken)
	if err != nil {
		return nil, err
	}
	if val, ok := bugsInfo["bugs"]; ok {
		if parsedArray, ok := val.([]interface{}); ok {
			result := make([]map[string]interface{}, len(parsedArray), len(parsedArray))
			for i := range result {
				if parsedItem, ok := parsedArray[i].(map[string]interface{}); ok {
					result[i] = parsedItem
				} else {
					return nil, fmt.Errorf("could not parse BugsInfo item in %v %v", reflect.TypeOf(parsedArray[i]), parsedArray[i])
				}
			}
			return result, nil
		}
		if parsedMap, ok := val.(map[string]interface{}); ok {
			return []map[string]interface{}{parsedMap}, nil
		}
		return nil, fmt.Errorf("could not parse BugsInfo result in %v %v", reflect.TypeOf(val), val)
	}
	return nil, fmt.Errorf("no 'bugs' field in %v", bugsInfo)
}

// BugHistory returns history of selected bugzilla ticket
func (client *Client) BugHistory(id int) (bugInfo map[string]interface{}, err error) {
	return client.BugsHistory([]int{id})
}

// BugsHistory returns history of selected bugzilla tickets
func (client *Client) BugsHistory(idList []int) (bugInfo map[string]interface{}, err error) {
	return client.json.bugsHistory(idList, client.bugzillaToken)
}

// AddComment adds comment for selected bugzilla ticket
func (client *Client) AddComment(id int, comment string) (bugInfo map[string]interface{}, err error) {
	return client.json.addComment(id, client.bugzillaToken, comment)
}
