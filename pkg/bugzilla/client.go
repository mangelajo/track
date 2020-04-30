package bugzilla

import (
	"fmt"
	"io"
	"net/http"
	"net/http/cookiejar"
	"reflect"
	"time"

	"github.com/howeyc/gopass"
	"bufio"
	"os"
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


type GetAuthFunc func() ([]*http.Cookie, *string)
type StoreAuthFunc func(cookies []*http.Cookie, authToken string)

// NewClient creates bugzilla Client instance
func NewClient(bugzillaAddress string, bugzillaLogin string, bugzillaPassword string,
	           getAuth GetAuthFunc,
	           	storeAuth StoreAuthFunc) (client *Client, err error) {

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

	var authToken *string
	var cookies []*http.Cookie

	if getAuth != nil {
		cookies, authToken = getAuth()
		if cookies != nil && authToken != nil {
			jsonClient.SetCookies(cookies)
			cgiClient.SetCookies(cookies)
		}

	}

	if authToken == nil || cookies == nil {

		fmt.Println("Oh, we don't have an auth token or cookies...")

		if bugzillaLogin == "" {
			reader := bufio.NewReader(os.Stdin)

			fmt.Print("Bugzilla Email: ")
			bugzillaLogin, _ = reader.ReadString('\n')

		}

		if bugzillaPassword == "" {

			fmt.Print("Bugzilla Password: ")
			pass, _ := gopass.GetPasswdMasked()
			bugzillaPassword = string(pass)
		}

		fmt.Print("Authenticating to bugzilla via JSON.... ")
		token, err := jsonClient.login(bugzillaLogin, bugzillaPassword)
		authToken = &token
		if err != nil {
			return nil, err
		}

		fmt.Print("Authenticating to bugzilla via CGI.... ")
		if err := cgiClient.login(bugzillaLogin, bugzillaPassword); err != nil {
			return nil, err
		}

		cookies = cgiClient.GetCookies()
		if storeAuth != nil {
			storeAuth(cookies, *authToken)
		}
		fmt.Println("done")
	}

	client = &Client{
		bugzillaAddress:  bugzillaAddress,
		bugzillaLogin:    bugzillaLogin,
		bugzillaPassword: bugzillaPassword,
		bugzillaToken:    *authToken,
		cgi:              cgiClient,
		json:             jsonClient,
	}

	return client, nil
}

type BugListQuery struct {
	CustomQuery string
	Limit int
	Offset int
	Order string
	Classification string
	Product string
	Component string
	BugStatus []string
	WhiteBoard string
	AssignedTo string
	FlagRequestee string
	TargetRelease string
	TargetMilestone string

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

func (client *Client) ShowBug(id int, currentTimestamp string)  (bug *Cbug, cached bool, err error) {
	return client.cgi.bugInfo(id, currentTimestamp)
}

func (client *Client) ShowBugHTML(id int, currentTimestamp string)  (html *[]byte, cached bool, err error) {
	return client.cgi.bugInfoHTML(id, currentTimestamp)
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
