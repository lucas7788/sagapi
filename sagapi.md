- 接口规范
  - 通用接口
  - API接口

#### 接口规范

​	响应规范中的结构体属于以下结构体的result域

```
{
    "result":""
    "error": 20，
    “description":""，
    "version":"v1"，
}
```

| Field_Name  | Type        | Description    |
| ----------- | ----------- | -------------- |
| result      | interface{} | 见具体接口描述 |
| error       | int         | 错误类型       |
| description | string      | 错误描述       |
| version     | string      | 版本号         |

| Error              | Num   | Description    |
| ------------------ | ----- | -------------- |
| SUCCESS            | 1     | 成功           |
| PARA_ERROR         | 40000 | 参数错误       |
| INTER_ERROR        | 40001 | 内部错误       |
| SQL_ERROR          | 40002 | 数据库操作失败 |
| VERIFY_TOKEN_ERROR | 40003 | 授权失败       |

##### 通用接口

##### 下单

```
url:/api/v1/order/takeOrder
method: post
```

- 请求

  ```
  {
      "productName": "xxx",
      "ontId": "xxx",
      "userName": "steven",
      "apiId": "9",
      "specifications":10,
  }
  ```

| Field_Name     | Type   | Description       |
| -------------- | ------ | ----------------- |
| productName    | string | nasa              |
| ontId          | string | ontId             |
| userName       | string | 用户名            |
| apiId          | int    | Id of request Api |
| specifications | int    | times of api      |

- 响应

  ```
  {
  	"qrCode": {		         "ONTAuthScanProtocol":"...api/v1/order/getQrCodeDataByQrCodeId“ (test)
  	}
  	"id":"4c9e3211-3059-4de1-ab9a-d0fc82733c7"
  }
  ```

  | Field_Name          | Type   | Description       |
  | ------------------- | ------ | ----------------- |
  | ONTAuthScanProtocol | string | 获取二维码参数url |
  | id                  | string | 二维码id          |

##### 查询二维码

```
url:/api/v1/order/getQrCodeByOrderId/{orderId}
metod：get 
```

- ##### 请求

| Field_Name | Type   | Description |
| ---------- | ------ | ----------- |
| oderId     | string | 订单ID      |

- 响应

```
{
	"qrCode": {		         "ONTAuthScanProtocol":"...api/v1/order/getQrCodeDataByQrCodeId“ (test)
	}
	"id":"4c9e3211-3059-4de1-ab9a-d0fc82733c7"
}
```

| Field_Name          | Type   | Description       |
| ------------------- | ------ | ----------------- |
| ONTAuthScanProtocol | string | 获取二维码参数url |
| id                  | string | 二维码id          |

##### 获取二维码参数数据

```
url:/api/v1/order/getQrCodeDataByQrCodeId/{qrCodeId}
metod：get
```

- 请求

| Field_Name | Type   | Description |
| ---------- | ------ | ----------- |
| qrCodeID   | string | 二维码ID    |

- 响应

```
{"requester":"did:ont:AYCcjQuB6xgXm2vKku9Vb6bdTcEguXqbt1","ver":"v2.0.0","chain":"Testnet","data":"{\"action\":\"signMessage\",\"params\":{\"invokeConfig\":{\"gasLimit\":20000,\"contractHash\":\"0000000000000000000000000000000000000000\",\"functions\":[{\"args\":[{\"name\":\"register\",\"value\":\"String:ontId\"}],\"operation\":\"signMessage\"}],\"payer\":\"\",\"gasPrice\":500}}}","signature":"01db6129d50852da913ba1bfabc2ab6e81741ac30cd2a097a7ca763722ef96d5cce3f20d85998b873986a90e1ff23ea6b5a478725d5b6593c0e534d051e3678bb3","callback":"http://a9279cdf5639211ea83090a4ed185dbd-544314116.ap-southeast-1.elb.amazonaws.com/api/v1/account/register/callback","id":"4c9e3211-3059-4de1-ab9a-d0fc82733c78","exp":1585292006,"signer":"","desc":{}}
```

##### 发送交易

```
url:/api/v1/order/sendTx
metod：post
```

- 请求

```
{
	”signer": "xxx",
	"signedTx": "xxxx",
	"extraData": {
        "id": "4c9e3211-3059-4de1-ab9a-d0fc82733c78"，
        “publickey": "xxx",
        "ontId": "xxx"，
	}
}
```

| Field_Name | Type   | Description |
| ---------- | ------ | ----------- |
| signer     | string | 签名用户    |
| signedTx   | string | 签名        |
| id         | string | 二维码Id    |
| publickey  | string | 公钥        |
| ontid      | string | ontid       |

- ##### 响应

string:"SUCCESS"

##### 取消订单

```
{
	/api/v1/order/cancelOrder
	method:post
}
```

- 请求

```
{
    "OrderId":"4c9e3211-3059-4de1-ab9a-d0fc82733c7"
}
```

| Field_Name | Type   | Description |
| ---------- | ------ | ----------- |
| OderId     | string | 订单ID      |

##### 删除订单

```
{
	/api/v1/order/deleteOrder
	method:post
}
```

- 请求

```
{
    "OrderId":"4c9e3211-3059-4de1-ab9a-d0fc82733c7"
}
```

| Field_Name | Type   | Description |
| ---------- | ------ | ----------- |
| OderId     | string | 订单ID      |

##### 查询订单状态

```
{
	/api/v1/order/queryOrderStatu/{orderId}
	method:get
}
```

- 请求

| Field_Name | Type   | Description |
| ---------- | ------ | ----------- |
| oderId     | string | 订单ID      |

- ##### 响应

```
{
    result: "1"，
    userName: "yy"，
    ontId：”yyy"，
}
```

| Field_Name | Type   | Description                        |
| ---------- | ------ | ---------------------------------- |
| result     | string | "1": 完成， ""：处理中， ”0“：失败 |

##### 生成测试键

```
{
	/api/v1/order/generateTestKey
	method:post
}
```

- 请求

```
{
	"apiId":"xxx"
}
```

| Field_Name | Type | Description |
| ---------- | ---- | ----------- |
| apiId      | int  | Api ID      |

##### 响应

```
{
    "Id": 20,
    "ApiKey": "xxx",
    "OrderId" "xxx",
    "ApiId": 30,
    "RequestLimit": 40,
    "UseNum": 50,
    "OntId": "xxxx",
    "OderStatus": 1
}
```

##### 查询API页

```
{
	/getBasicApiInfoByPage/{pageNum}/{pageSize}
	method:get
}
```

- 请求

| Field_Name | Type   | Description |
| ---------- | ------ | ----------- |
| pageNum    | string | 页数量      |
| pageSize   | strnig | 页大小      |

- ##### 响应

```
[
	{
        "apiId": 20，
        "coin": "xxx"，
        "type":"xxx",
        "icon":"xxx",
        "title":"xxx",
        "provider":"yyy"，
        “apiUrl":"",
        "price":"",
        "description":"xxx",
        "specifications": 30,
        "popularity":20,
        "delay":20,
        "successRate":20,
        "invokeFrequency": 100,
        "createTime": "yyy",
	}，
	...
	{
        "apiId": 20，
        "coin": "xxx"，
        "type":"xxx",
        "icon":"xxx",
        "title":"xxx",
        "provider":"yyy"，
        “apiUrl":"",
        "price":"",
        "description":"xxx",
        "specifications": 30,
        "popularity":20,
        "delay":20,
        "successRate":20,
        "invokeFrequency": 100,
        "createTime": "yyy",
	}
]
```

##### 查询api详细信息

```
{
	/getApiDetailByApiId/{apiId}
	method:get
}
```

- 请求

| Field_Name | Type   | Description |
| ---------- | ------ | ----------- |
| apiId      | string | api编码     |

- ##### 响应

```
{
	"ApiId":20,
	"Mark":"",
	"ResponseParam":"",
	"ResponseExample":"",
	"DataDesc":"",
	"DataSource:"",
	"ApplicationScenario":"",
	[
		{
        "ApiDetailInfoId":20,
        "ParamName":"xxx",
        "Required":1,
        "ParamType": "yyy",
        "Note":"xxx"
        },
        ...
	]，
	...
	[
	 {
        "ApiDetailInfoId":20，
        "ErrorCode":20，
        "Errordesc":"xxx"，
      },
      ...
	],
	[
	 {
        "Id":20，
        "ApiDetailInfoId":20，
        "Price":"xxx"，
        "Amount":30,
      },
      ...
	],
	{
        "apiId": 20，
        "coin": "xxx"，
        "type":"xxx",
        "icon":"xxx",
        "title":"xxx",
        "provider":"yyy"，
        “apiUrl":"",
        "price":"",
        "description":"xxx",
        "specifications": 30,
        "popularity":20,
        "delay":20,
        "successRate":20,
        "invokeFrequency": 100,
        "createTime": "yyy",
	}
}
```

| Field_Name | Type   | Description |
| ---------- | ------ | ----------- |
| datadesc   | string | 数据描述    |

##### 根据apikey查询api信息

```
{
	/SearchApiByKey/{apikey}
	method:get
}
```

- 请求

| Field_Name | Type   | Description |
| ---------- | ------ | ----------- |
| apikey     | string | api授权key  |

- ##### 响应

```
[
	{
        "apiId": 20，
        "coin": "xxx"，
        "type":"xxx",
        "icon":"xxx",
        "title":"xxx",
        "provider":"yyy"，
        “apiUrl":"",
        "price":"",
        "description":"xxx",
        "specifications": 30,
        "popularity":20,
        "delay":20,
        "successRate":20,
        "invokeFrequency": 100,
        "createTime": "yyy",
	},
  ...
]
```

##### 查询所有hot， new, free api.

- 请求

  ```
  {
  	/searchApi
  	method:get
  }
  ```

- 响应

```
[
	"newest": [
		{
        "apiId": 20，
        "coin": "xxx"，
        "type":"xxx",
        "icon":"xxx",
        "title":"xxx",
        "provider":"yyy"，
        “apiUrl":"",
        "price":"",
        "description":"xxx",
        "specifications": 30,
        "popularity":20,
        "delay":20,
        "successRate":20,
        "invokeFrequency": 100,
        "createTime": "yyy",
		}，
		....
	]
	"hottest":[
		{
        "apiId": 20，
        "coin": "xxx"，
        "type":"xxx",
        "icon":"xxx",
        "title":"xxx",
        "provider":"yyy"，
        “apiUrl":"",
        "price":"",
        "description":"xxx",
        "specifications": 30,
        "popularity":20,
        "delay":20,
        "successRate":20,
        "invokeFrequency": 100,
        "createTime": "yyy",
		},
		...
	]
	"free":[
        {
        "apiId": 20，
        "coin": "xxx"，
        "type":"xxx",
        "icon":"xxx",
        "title":"xxx",
        "provider":"yyy"，
        “apiUrl":"",
        "price":"",
        "description":"xxx",
        "specifications": 30,
        "popularity":20,
        "delay":20,
        "successRate":20,
        "invokeFrequency": 100,
        "createTime": "yyy",
		},
	]
]
```

##### 根据分类ID查API

##### 请求

```
{
	/searchApiByCategory/:categoryId
	method:get
}
```

##### 响应

```
[
	{
        "apiId": 20，
        "coin": "xxx"，
        "type":"xxx",
        "icon":"xxx",
        "title":"xxx",
        "provider":"yyy"，
        “apiUrl":"",
        "price":"",
        "description":"xxx",
        "specifications": 30,
        "popularity":20,
        "delay":20,
        "successRate":20,
        "invokeFrequency": 100,
        "createTime": "yyy",
	}，
	....
]
```

##### API接口

##### apod

```
{
	/api/v1/apod/{apikey}
	method:get
}
```

| Field_Name | Type   | Description         |
| ---------- | ------ | ------------------- |
| apikey     | string | 下单后的api访问权限 |

- 响应、

| Field_Name | Type   | Description |
| ---------- | ------ | ----------- |
| data       | string | 返回数据    |

##### feed

```
{
	/api/v1/feed/{startdate}/{enddate}/{apikey}
	method:get
}
```

| Field_Name | Type   | Description         |
| ---------- | ------ | ------------------- |
| startdate  | string | 开始日期            |
| enddate    | string | 结束日期            |
| apikey     | string | 下单后的api访问权限 |

- 响应、

| Field_Name | Type   | Description |
| ---------- | ------ | ----------- |
| data       | string | 返回数据    |
