# alibaba opensearch go client

## 对阿里巴巴开放搜索的请求封装

- [官方文档](https://help.aliyun.com/document_detail/opensearch/api-reference/api-interface/data-manager.html)
- 目前只简单封装了ListApp\AddDoc\SingleSearch三个接口
- [当前API文档](https://godoc.org/github.com/AsGz/opensearch)

##示例

```go
//列出当前账户所有的APP
var cf opensearch.Config
cf.OS_ACCESS_KEY = "xxxxxxxxxx"
//这里在私钥后面要加上&符号
cf.OS_SECRET_KEY = "xxxxxxxxxxxx&"
cf.OS_HOST = "http://opensearch-cn-hangzhou.aliyuncs.com"
cf.OS_APPNAME = "xxxxxxxxx"
o, err := opensearch.NewOpenSearchClient(cf)
if err != nil {
	fmt.Println(err)
} else {
	fmt.Println(o.ListApp())
}
```


```go
//添加文档,建立索引
ar := o.PushDoc("table_name", formatJson(items))
fmt.Printf("%#v\n", ar)

//进行单索引库的检索
indexName := "default"
r := o.SingleSearch(indexName, "abc", 0, 10)
fmt.Printf("%#v\n", r)
```
