GET http://127.0.0.1:8081/
HTTP 200

[Asserts]
jsonpath "$.env.TOML_VAR" == "hello toml"
