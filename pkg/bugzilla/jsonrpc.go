package bugzilla

import (
	"bytes"
	"encoding/json"
	"errors"
	"fmt"
	"net/http"
	"net/url"
	"sync"
)

type clientRequest struct {
	Method string         `json:"method"`
	Params [1]interface{} `json:"params"`
	ID     uint64         `json:"id"`
}

type bugzillaError struct {
	Message string
	Code    uint64
}

type clientResponse struct {
	ID     uint64           `json:"id"`
	Result *json.RawMessage `json:"result"`
	Error  *bugzillaError   `json:"error"`
}

// bugzillaJSONRPCClient bugzilla JSON RPC client
type bugzillaJSONRPCClient struct {
	bugzillaAddr string
	jsonRPCAddr  string
	httpClient   *http.Client
	seq          uint64
	m            sync.Mutex
}

// newJSONRPCClient creates a helper json rpc client for regular HTTP based endpoints
func newJSONRPCClient(addr string, httpClient *http.Client) (*bugzillaJSONRPCClient, error) {
	u, err := url.Parse(addr)
	if err != nil {
		return nil, err
	}
	u.Path = "jsonrpc.cgi"

	return &bugzillaJSONRPCClient{
		bugzillaAddr: addr,
		jsonRPCAddr:  u.String(),
		httpClient:   httpClient,
		seq:          0,
		m:            sync.Mutex{},
	}, nil
}

// Login allows to login using Bugzilla JSONRPC API, returns token
func (client *bugzillaJSONRPCClient) login(login string, password string) (token string, err error) {
	args := make(map[string]interface{})
	args["login"] = login
	args["password"] = password
	args["remember"] = true

	var result map[string]interface{}
	err = client.call("User.login", &args, &result)
	if err != nil {
		return "", err
	}
	token, ok := result["token"].(string)
	if !ok {
		return "", fmt.Errorf("could not parse token %v", result["token"])
	}
	return token, err
}

func (client *bugzillaJSONRPCClient) GetCookies() []*http.Cookie {
	url, _ := url.Parse(client.bugzillaAddr)
	cookies := client.httpClient.Jar.Cookies(url)
	return cookies
}

func (client *bugzillaJSONRPCClient) SetCookies(cookies []*http.Cookie) {
	url, _ := url.Parse(client.bugzillaAddr)
	client.httpClient.Jar.SetCookies(url, cookies)
}


// bugzillaVersion returns Bugzilla version
func (client *bugzillaJSONRPCClient) bugzillaVersion() (version string, err error) {
	var result map[string]interface{}
	err = client.call("Bugzilla.version", nil, &result)
	if err != nil {
		return "", err
	}
	version, ok := result["version"].(string)
	if !ok {
		return "", fmt.Errorf("could not parse token %v", result["version"])
	}
	return version, nil
}

// bugsInfo returns information about selected bugzilla tickets
func (client *bugzillaJSONRPCClient) bugsInfo(idList []int, token string) (bugInfo map[string]interface{}, err error) {
	args := make(map[string]interface{})
	args["ids"] = idList
	args["token"] = token

	err = client.call("Bug.get", args, &bugInfo)
	if err != nil {
		return nil, err
	}
	return bugInfo, nil
}


// bugsHistory returns history of selected bugzilla tickets
func (client *bugzillaJSONRPCClient) bugsHistory(idList []int, token string) (bugInfo map[string]interface{}, err error) {
	args := make(map[string]interface{})
	args["ids"] = idList
	args["token"] = token

	err = client.call("Bug.history", args, &bugInfo)
	if err != nil {
		return nil, err
	}
	return bugInfo, nil
}

// bugsHistory returns history of selected bugzilla tickets
func (client *bugzillaJSONRPCClient) addComment(id int, token string, comment string) (commentInfo map[string]interface{}, err error) {
	args := make(map[string]interface{})
	args["id"] = id
	args["token"] = token
	args["comment"] = comment

	err = client.call("Bug.add_comment", args, &commentInfo)
	if err != nil {
		return nil, err
	}
	return commentInfo, nil
}

// call performs JSON RPC call
func (client *bugzillaJSONRPCClient) call(serviceMethod string, args interface{}, reply interface{}) error {
	var params [1]interface{}
	params[0] = args

	client.m.Lock()
	seq := client.seq
	client.seq++
	client.m.Unlock()

	cr := &clientRequest{
		Method: serviceMethod,
		Params: params,
		ID:     seq,
	}

	byteData, err := json.Marshal(cr)
	if err != nil {
		return err
	}

	req, err := newHTTPRequest("POST", client.jsonRPCAddr, bytes.NewReader(byteData))
	if err != nil {
		return err
	}
	req.Header.Set("Content-Type", "application/json-rpc")

	res, err := client.httpClient.Do(req)
	defer func() {
		if res != nil && res.Body != nil {
			res.Body.Close()
		}
	}()
	if err != nil {
		return err
	}

	//body, _ := ioutil.ReadAll(res.Body)
	//fmt.Printf("%v", string(body))
	if res.StatusCode < 200 || res.StatusCode >= 300 {
		return fmt.Errorf("Response status: %v", res.StatusCode)
	}

	v := &clientResponse{}

	err = json.NewDecoder(res.Body).Decode(v)
	if err != nil {
		return err
	}

	if v.Error != nil {
		return errors.New(v.Error.Message)
	}

	//fmt.Println(string(*v.Result))
	return json.Unmarshal(*v.Result, reply)
}
