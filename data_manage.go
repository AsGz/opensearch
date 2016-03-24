package opensearch

import (
	"fmt"
)

//添加文档到服务器
func (o *OpenSearchClient) PushDoc(aliResult *AliResult, tableName, itemsJson string) error {
	var params ParamsList
	params = append(params, Param{"action", "push"})
	params = append(params, Param{"table_name", tableName})
	sign, queryString := o.getAliSign(params, "POST")
	url := fmt.Sprintf("%s/index/doc/%s?%s&Signature=%s", o.cf.OS_HOST, o.cf.OS_APPNAME, queryString, sign)
	//fmt.Println(url)
	//itemsData := fmt.Sprintf("items=%s", aliEncode(itemsJson))
	itemsData := "items=" + itemsJson
	return doHttpRequest(aliResult, url, "POST", itemsData)
}
