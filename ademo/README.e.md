# A demo

A demo of ggt capabilities to create a service to read/create `tomatoes`.

# The main

{{cat "main.go" | color "go"}}

# The controller

{{cat "controller/tomate.go" | color "go"}}

# The test

```sh
[mh-cbon@pc4 ademo] $ go generate *go
2017/05/24 15:43:01 no initial packages were loaded
2017/05/24 15:43:01 no initial packages were loaded
2017/05/24 15:43:01 no initial packages were loaded
[mh-cbon@pc4 ademo] $ go run *go &
[1] 5833
[mh-cbon@pc4 ademo] $ curl http://localhost:8080/GetById?id=0
{"ID":"0","Color":"Red"}
[mh-cbon@pc4 ademo] $ curl --data "color=blue" http://localhost:8080/Create
{"ID":"1","Color":"blue"}
[mh-cbon@pc4 ademo] $ curl --data "color=blue" http://localhost:8080/Create
2017/05/24 15:43:14 http: multiple response.WriteHeader calls
color must be unique
null
[mh-cbon@pc4 ademo] $ curl --data "color=" http://localhost:8080/Create
2017/05/24 15:43:17 http: multiple response.WriteHeader calls
color must not be empty
null
[mh-cbon@pc4 ademo] $ curl --data "color=green" http://localhost:8080/Create
{"ID":"2","Color":"green"}
[mh-cbon@pc4 ademo] $ curl http://localhost:8080/GetById?id=1
{"ID":"1","Color":"blue"}
[mh-cbon@pc4 ademo] $ curl http://localhost:8080/GetById?id=2
{"ID":"2","Color":"green"}
[mh-cbon@pc4 ademo] $ fg
go run *go
^Csignal: interrupt
[mh-cbon@pc4 ademo] $
```
