package opensearch

import (
	"encoding/json"
	"errors"
	"fmt"
	"io/ioutil"
	"net/http"
	"strings"
)

const (
	S_SUCCESS = "OK"
	S_FAILED  = "FAIL"
)

type Config struct {
	OS_ACCESS_KEY string
	OS_SECRET_KEY string
	OS_HOST       string
	OS_APPNAME    string
}

type OpenSearchClient struct {
	cf	Config
	seq	int64
}

type AliResult struct {
	Status     string
	Request_id string
	Result     interface{}
	Errors     interface{}
}

//创建一个client
func NewOpenSearchClient(cf Config) (*OpenSearchClient, error) {
	o := new(OpenSearchClient)
	o.cf = cf
	result := o.ListApp()
	if result.Status != S_SUCCESS {
		return nil, errors.New(fmt.Sprintf("%s", result.Errors))
	}
	return o, nil
}

func doHttpRequest(url, method, params string) *AliResult {
	var result = &AliResult{}
	var err error
	var request *http.Request
	request, err = http.NewRequest(method, url, strings.NewReader(params))
	if err != nil {
		result.Status = S_FAILED
		result.Errors = err
		return result
	}
	request.Header.Set("Content-Type", "application/x-www-form-urlencoded")
	request.Header.Set("Connection", "Keep-Alive")
	var resp *http.Response
	resp, err = http.DefaultClient.Do(request)
	if err != nil {
		result.Status = S_FAILED
		result.Errors = errors.New(fmt.Sprintf("resp:%d, err=%s", resp, err))
		return result
	}
	defer resp.Body.Close()
	var b []byte
	b, err = ioutil.ReadAll(resp.Body)
	if err != nil || resp.StatusCode != 200 {
		result.Status = S_FAILED
		result.Errors = errors.New(fmt.Sprintf("[httpcode:%d][err:%s]", resp.StatusCode, err))
		return result
	} else {
		err = json.Unmarshal(b, &result)
		if err != nil {
			result.Status = S_FAILED
			result.Errors = err
		}
	}
	return result
}

//列出所有应用
func (o *OpenSearchClient) ListApp() *AliResult {
	var params ParamsList
	sign, queryString := o.getAliSign(params, "GET")
	url := fmt.Sprintf("%s/index?%s&Signature=%s", o.cf.OS_HOST, queryString, sign)
	return doHttpRequest(url, "GET", "")
}
