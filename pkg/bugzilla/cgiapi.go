package bugzilla

import (
	"fmt"
	"io/ioutil"
	"net/http"
	"net/url"
	"regexp"
	"strconv"
	"strings"
	"encoding/xml"
)

// bugzillaCGIClient bugzilla REST API client
type bugzillaCGIClient struct {
	bugzillaAddr string
	httpClient   *http.Client
}

// NewCGIClient creates a helper json rpc client for regular HTTP based endpoints
func newCGIClient(addr string, httpClient *http.Client) (*bugzillaCGIClient, error) {
	return &bugzillaCGIClient{
		bugzillaAddr: addr,
		httpClient:   httpClient,
	}, nil
}

// setBugzillaLoginCookie visits bugzilla page to obtain login cookie
func (client *bugzillaCGIClient) setBugzillaLoginCookie(loginURL string) (err error) {
	req, err := newHTTPRequest("GET", loginURL, nil)
	if err != nil {
		return err
	}

	res, err := client.httpClient.Do(req)
	defer func() {
		if res != nil && res.Body != nil {
			res.Body.Close()
		}
	}()

	if err != nil {
		if strings.Contains(err.Error(), "use of closed network connection") {
			return fmt.Errorf("Timeout occured while accessing %v", loginURL)
		}
		return err
	}
	return nil
}

// getBugzillaLoginToken returns Bugzilla_login_token input field value. Requires login cookie to be set
func (client *bugzillaCGIClient) getBugzillaLoginToken(loginURL string) (loginToken string, err error) {
	req, err := newHTTPRequest("GET", loginURL, nil)
	if err != nil {
		return "", err
	}

	res, err := client.httpClient.Do(req)
	defer func() {
		if res != nil && res.Body != nil {
			res.Body.Close()
		}
	}()
	if err != nil {
		if strings.Contains(err.Error(), "use of closed network connection") {
			return "", fmt.Errorf("Timeout occured while accessing %v", loginURL)
		}
		return "", err
	}
	//<input type="hidden" name="Bugzilla_login_token" value="1435647781-eV7m3mhmosArYikHPtaisTliTn7e3kKOZ-RhiX-Qz1A">
	r := regexp.MustCompile(`name="Bugzilla_login_token"\s+value="(?P<value>[\d\w-]+)"`)
	body, err := ioutil.ReadAll(res.Body)
	if err != nil {
		return "", err
	}
	match := r.FindStringSubmatch(string(body))
	a := make(map[string]string)
	for i, name := range r.SubexpNames() {
		a[name] = match[i]
	}
	return a["value"], nil
}

// Login allows to login using Bugzilla CGI API
func (client *bugzillaCGIClient) login(login string, password string) (err error) {
	u, err := url.Parse(client.bugzillaAddr)
	if err != nil {
		return err
	}
	u.Path = "index.cgi"
	loginURL := u.String()

	err = client.setBugzillaLoginCookie(loginURL)
	if err != nil {
		return err
	}

	loginToken, err := client.getBugzillaLoginToken(loginURL)
	if err != nil {
		return err
	}

	data := url.Values{}
	data.Set("Bugzilla_login", login)
	data.Set("Bugzilla_password", password)
	data.Set("Bugzilla_login_token", loginToken)

	req, err := newHTTPRequest("POST", loginURL, strings.NewReader(data.Encode()))
	if err != nil {
		return err
	}
	req.Header.Set("Accept", "text/html")
	req.Header.Set("Content-Type", "application/x-www-form-urlencoded")

	res, err := client.httpClient.Do(req)
	defer func() {
		if res != nil && res.Body != nil {
			res.Body.Close()
		}
	}()
	if err != nil {
		if strings.Contains(err.Error(), "use of closed network connection") {
			return fmt.Errorf("Timeout occured while accessing %v", loginURL)
		}
		return err
	}
	return nil
}

// bugList list of last changed bugs
func (client *bugzillaCGIClient) bugList(limit int, offset int) ([]Bug, error) {
	u, err := url.Parse(client.bugzillaAddr)
	if err != nil {
		return nil, err
	}
	u.Path = "buglist.cgi"
	q := u.Query()
	q.Set("ctype", "rdf")
	q.Set("query_format", "advanced")
	q.Set("limit", strconv.Itoa(limit))
	q.Set("offset", strconv.Itoa(offset))
	q.Set("order", "changeddate DESC")
	q.Set("product", "Red Hat OpenStack")
	q.Set("component", "python-networking-ovn")
	q.Set("bug_status", "NEW")
	q.Set("classification", "Red Hat")


	u.RawQuery = q.Encode()

	//url = https://bugzilla.mozilla.org/buglist.cgi?format=simple&limit=4&query_format=advanced&offset=400&order=changeddate%20DESC
	req, err := newHTTPRequest("GET", u.String(), nil)
	if err != nil {
		return nil, err
	}
	req.Header.Set("Accept", "text/csv")

	res, err := client.httpClient.Do(req)
	defer func() {
		if res != nil && res.Body != nil {
			res.Body.Close()
		}
	}()
	if err != nil {
		if strings.Contains(err.Error(), "use of closed network connection") {
			return nil, fmt.Errorf("Timeout occured while accessing %v", req.URL)
		}
		return nil, err
	}

	bugList, err := parseBugzRDF(res.Body)
	if err != nil {
		return nil, err
	}
	results := make([]Bug, len(bugList))
	for i := range bugList {
		b, err := NewBugFromBzBug(bugList[i])
		if err != nil {
			return nil, err
		}
		results[i] = *b
	}
	return results, err
}


// bugList list of last changed bugs
func (client *bugzillaCGIClient) bugInfo(id int) (*Cbug, error) {
	u, err := url.Parse(client.bugzillaAddr)
	if err != nil {
		return nil, err
	}

	u.Path = "show_bug.cgi"
	q := u.Query()
	q.Set("ctype", "xml")
	q.Set("id", strconv.Itoa(id))

	u.RawQuery = q.Encode()

	//url = https://bugzilla.mozilla.org/show_bug.cgi?id=xxxxx&ctype=xml
	req, err := newHTTPRequest("GET", u.String(), nil)
	if err != nil {
		return nil, err
	}
	//req.Header.Set("Accept", "text/xml")

	res, err := client.httpClient.Do(req)
	defer func() {
		if res != nil && res.Body != nil {
			res.Body.Close()
		}
	}()
	if err != nil {
		if strings.Contains(err.Error(), "use of closed network connection") {
			return nil, fmt.Errorf("Timeout occured while accessing %v", req.URL)
		}
		return nil, err
	}

	var bugzilla Cbugzilla

	body, err := ioutil.ReadAll(res.Body)
	err = xml.Unmarshal(body ,&bugzilla)

	if err != nil {
		return nil, err
	}


	return bugzilla.Cbug, err
}
