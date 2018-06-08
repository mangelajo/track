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
	"github.com/mangelajo/track/pkg/storecache"
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

func setupQuery(u *url.URL, query *BugListQuery) (url string, referer string){

	u.Path = "buglist.cgi"

	var advMatches int = 1

	if query.CustomQuery != "" {
		u.RawQuery = query.CustomQuery + "&ctype=csv&human=1"
		url := u.String()
		u.RawQuery = query.CustomQuery
		referer := u.String()
		return url, referer
	}

	q := u.Query()

	q.Set("ctype","csv")
	q.Set("columnlist","product,component,bug_severity,cf_pm_score,"+
		                    "assigned_to,bug_status,short_desc,changeddate,resolution")
	q.Set("human", "1") // If we don't use this flag it will ignore some filters
	q.Set("query_format", "advanced")
	q.Set("limit", strconv.Itoa(query.Limit))
	q.Set("offset", strconv.Itoa(query.Offset))

	if query.Order == "" {
		q.Set("order", "changeddate DESC")
	} else {
		q.Set("order", query.Order)
	}

	if query.Product != "" {
		q.Set("product", query.Product)
	}

	if query.Component != "" {
		q.Set("component", query.Component)
	}

	for _, bs := range query.BugStatus {
		q.Add("bug_status", bs)
	}

	if query.Classification != "" {
		q.Set("classification", query.Classification)
	}

	if query.WhiteBoard != "" {
		q.Set("cf_internal_whiteboard", query.WhiteBoard)
		q.Set("cf_internal_whiteboard_type", "substring")
	}

	if query.AssignedTo != "" {
		q.Set("assigned_to", query.AssignedTo)
	}

	if query.FlagRequestee != "" {
		q.Set(fmt.Sprintf("f%d", advMatches), "requestees.login_name")
		q.Set(fmt.Sprintf("o%d", advMatches), "substring")
		q.Set(fmt.Sprintf("v%d", advMatches), query.FlagRequestee)
		advMatches++
	}

	u.RawQuery = q.Encode()
	return u.String(), ""
}

// bugList list of last changed bugs
func (client *bugzillaCGIClient) bugList(query *BugListQuery) ([]Bug, error) {

	u, err := url.Parse(client.bugzillaAddr)
	if err != nil {
		return nil, err
	}

	url, referer := setupQuery(u, query)

	//url = https://bugzilla.mozilla.org/buglist.cgi?format=simple&limit=4&query_format=advanced&offset=400&order=changeddate%20DESC

	req, err := newHTTPRequest("GET", url , nil)
	req.Header.Set("Upgrade-Insecure-Request", "1")
	req.Header.Set("DNT", "1")

	// For some weird reason it doesn't return the right list if no Referer is set
	if referer != "" {
		req.Header.Set("Referer", referer)
	}
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

	bugList, err := parseBugzCSV(res.Body)

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

func (client *bugzillaCGIClient) getBug(id int, currentTimestamp string, getXml bool) (xml *[]byte, cached bool, err error) {

	xml, err = storecache.RetrieveCache(id, currentTimestamp,getXml)

	if err == nil {
		return xml, true, err
	}

	u, err := url.Parse(client.bugzillaAddr)
	if err != nil {
		return nil, false, err
	}

	u.Path = "show_bug.cgi"
	q := u.Query()
	if getXml {
		q.Set("ctype", "xml")
		q.Set("excludefield", "attachmentdata")
	}
	q.Set("id", strconv.Itoa(id))

	u.RawQuery = q.Encode()

	//url = https://bugzilla.mozilla.org/show_bug.cgi?...
	req, err := newHTTPRequest("GET", u.String(), nil)
	if err != nil {
		return nil, false, err
	}

	res, err := client.httpClient.Do(req)
	defer func() {
		if res != nil && res.Body != nil {
			res.Body.Close()
		}
	}()
	if err != nil {
		if strings.Contains(err.Error(), "use of closed network connection") {
			return nil, false, fmt.Errorf("Timeout occured while accessing %v", req.URL)
		}
		return nil, false, err
	}


	body, err := ioutil.ReadAll(res.Body)

	storecache.StoreCache(id, currentTimestamp, &body, getXml)

	return &body, false,err
}

func (client *bugzillaCGIClient) bugInfo(id int, currentTimestamp string) (*Cbug,  bool, error) {

	var bugzilla Cbugzilla

	body, cached, err := client.getBug(id, currentTimestamp, true)
	err = xml.Unmarshal(*body ,&bugzilla)

	if err != nil {
		// invalidate cache
		storecache.StoreCache(id,"xxx", body, true)
		return nil, false, err
	}
	return bugzilla.Cbug, cached, err
}

func (client *bugzillaCGIClient) bugInfoHTML(id int, currentTimestamp string) (*[]byte,  bool, error) {


	body, cached, err := client.getBug(id, currentTimestamp, false)

	if err != nil {
		// invalidate cache
		storecache.StoreCache(id,"xxx", body, false)
		return nil, false, err
	}
	return body, cached, err
}

