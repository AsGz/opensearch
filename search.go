package opensearch

import (
	"fmt"
)

//fetch_fields:多个字段使用英文分号（;）分隔
func (o *OpenSearchClient) Search(result interface{}, fetch_fields string, query string) error {
	var params ParamsList
	params = append(params, Param{"index_name", o.cf.OS_APPNAME})
	params = append(params, Param{"query", query})
	params = append(params, Param{"fetch_fields", fetch_fields})
	sign, queryString := o.getAliSign(params, "GET")
	url := fmt.Sprintf("%s/search?%s&Signature=%s", o.cf.OS_HOST, queryString, sign)
	return doHttpRequest(result, url, "GET", "")
}
