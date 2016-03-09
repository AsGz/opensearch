package opensearch

import (
	"crypto/hmac"
	"crypto/sha1"
	"encoding/base64"
	"fmt"
	"math/rand"
	"net/url"
	"sort"
	"strings"
	"time"
)

const (
	TIME_FORMAT = "2006-01-02T15:04:05Z"
)

type Param struct {
	Key   string
	Value string
}

type ParamsList []Param

func (info ParamsList) Less(i, j int) bool { //升序
	return info[i].Key < info[j].Key
}

func (info ParamsList) Swap(i, j int) {
	info[i], info[j] = info[j], info[i]
}

func (info ParamsList) Len() int {
	return len(info)
}

//阿里的参数编码规则
func aliEncode(str string) string {
	encode := url.QueryEscape(str)
	encode = strings.Replace(encode, "+", "%20", -1)
	encode = strings.Replace(encode, "*", "%2A", -1)
	encode = strings.Replace(encode, "%7E", "~", -1)
	return encode
}

//组织公共参数
func (o *OpenSearchClient) getPublicParams() ParamsList {
	rand.Seed(time.Now().UnixNano())
	o.seq++
	signNonce := fmt.Sprintf("%d%d%d", time.Now().Unix(), rand.Intn(1000), o.seq)
	var params ParamsList
	params = append(params, Param{"AccessKeyId", o.cf.OS_ACCESS_KEY})
	params = append(params, Param{"SignatureMethod", "HMAC-SHA1"})
	params = append(params, Param{"SignatureVersion", "1.0"})
	params = append(params, Param{"SignatureNonce", signNonce})
	params = append(params, Param{"Timestamp", time.Now().UTC().Format(TIME_FORMAT)})
	params = append(params, Param{"Version", "v2"})
	//添加文档时, items参数不纳入签名
	// params = append(params, Param{"sign_mode", "1"})
	return params
}

//根据参数、方法获取签名
func (o *OpenSearchClient) getAliSign(params ParamsList, method string) (string, string) {
	publicParams := o.getPublicParams()
	params = append(params, publicParams...)
	//参数按字典序进行排序
	sort.Sort(params)
	var queryParams []string
	for _, v := range params {
		s := fmt.Sprintf("%s=%s", v.Key, aliEncode(v.Value))
		queryParams = append(queryParams, s)
	}
	queryString := strings.Join(queryParams, "&")
	stringToSign := fmt.Sprintf("%s&%s&%s", method, aliEncode("/"), aliEncode(queryString))
	hash := hmac.New(sha1.New, []byte(o.cf.OS_SECRET_KEY))
	hash.Write([]byte(stringToSign))
	sign := base64.StdEncoding.EncodeToString(hash.Sum(nil))
	sign = aliEncode(sign)
	return sign, queryString
}
