GET http://127.0.0.1:8081/?firstQuery=123&secondQuery=hello
HTTP 200

[Asserts]
jsonpath "$.queryStringParameters.firstQuery" == "123"
jsonpath "$.queryStringParameters.secondQuery" == "hello"
