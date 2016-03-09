package opensearch

import (
	"encoding/json"
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

type AliErr struct {
    Code int
    Message string
}

type AliResult struct {
	Status     string
	Errors     []*AliErr
    Result     interface{}
}

func (m *AliResult) IsOK() bool {
    return m.Status == S_SUCCESS
}

func (m *AliResult) IsFailed() bool {
    return m.Status == S_FAILED
}

func (m *AliResult) MakeError() error {
    n := len(m.Errors)
    if n == 0 {
        if m.IsOK() {
            return nil
        } else if m.IsFailed() {
            return fmt.Errorf("Status is failed, but no detail error")
        } else {
            return fmt.Errorf("Happen some error before aliyun")
        }
    } else if n==1 {
        return fmt.Errorf("code=%d,message=%s",m.Errors[0].Code,m.Errors[0].Message)
    } else {
        s := ""
        for i, e := range m.Errors {
            s += fmt.Sprintf("code[%d]=%d,message[%d]=%s\n",i,e.Code,i,e.Message)
        }
        return fmt.Errorf("%s",s)
    }
}

//创建一个client
func NewOpenSearchClient(result *AliResult, cf Config) (*OpenSearchClient, error) {
	o := new(OpenSearchClient)
	o.cf = cf
	err := o.ListApp(result)
	if err != nil {
		return nil, err
	}
	return o, nil
}

func doHttpRequest(result interface{}, url, method, params string) error {
	var err error
	var request *http.Request
	request, err = http.NewRequest(method, url, strings.NewReader(params))
	if err != nil {
		return err
	}
	request.Header.Set("Content-Type", "application/x-www-form-urlencoded")
	request.Header.Set("Connection", "Keep-Alive")
	var resp *http.Response
	resp, err = http.DefaultClient.Do(request)
	if err != nil {
		return err
	}
	defer resp.Body.Close()
	var b []byte
	b, err = ioutil.ReadAll(resp.Body)
    if err != nil {
        return err
    }
	if resp.StatusCode != 200 {
		return fmt.Errorf("HTTPCode=%d", resp.StatusCode)
	} else {
		err = json.Unmarshal(b, &result)
		if err != nil {
            fmt.Printf("json=%s\n",string(b))
			return err
		}
	}
	return nil
}

//列出所有应用
/*
{
    "result":{
        "index_name":"test_create_index"
    },
    "status":"OK", 
    "RequestId":"1422264399046999200623000"
 }
*/
func (o *OpenSearchClient) ListApp(result *AliResult) error {
	var params ParamsList
	sign, queryString := o.getAliSign(params, "GET")
	url := fmt.Sprintf("%s/index?%s&Signature=%s", o.cf.OS_HOST, queryString, sign)
	return doHttpRequest(result, url, "GET", "")
}
