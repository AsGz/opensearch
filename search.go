package opensearch

import (
	"fmt"
)

//对单个索引进行搜索
func (o *OpenSearchClient) SingleSearch(indexName, word string, start, hit int) *AliResult {
	var params ParamsList
	params = append(params, Param{"index_name", o.cf.OS_APPNAME})
	query := fmt.Sprintf("config=start:%d,hit:%d,format=fulljson&&query=%s:'%s'", start, hit, indexName, word)
	params = append(params, Param{"query", query})
	sign, queryString := o.getAliSign(params, "GET")
	url := fmt.Sprintf("%s/search?%s&Signature=%s", o.cf.OS_HOST, queryString, sign)
	return doHttpRequest(url, "GET", "")
}
