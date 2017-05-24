# A demo

{{pkgdoc}}

# {{toc 5}}

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
...
< HTTP/1.1 200 OK
< Content-Type: application/json
< Date: Wed, 24 May 2017 13:50:28 GMT
< Content-Length: 25
<
{"ID":"0","Color":"Red"}

[mh-cbon@pc4 ademo] $ curl http://localhost:8080/GetById?id=10
...
< HTTP/1.1 404 Not Found
< Content-Type: text/plain; charset=utf-8
< X-Content-Type-Options: nosniff
< Date: Wed, 24 May 2017 13:52:26 GMT
< Content-Length: 17
<
Tomate not found

[mh-cbon@pc4 ademo] $ curl -v --data "color=blue" http://localhost:8080/Create
...
< HTTP/1.1 200 OK
< Content-Type: application/json
< Date: Wed, 24 May 2017 13:49:58 GMT
< Content-Length: ...
<
{"ID":"1","Color":"blue"}

[mh-cbon@pc4 ademo] $ curl --data "color=blue" http://localhost:8080/Create
...
< HTTP/1.1 400 Bad Request
< Content-Type: text/plain; charset=utf-8
< X-Content-Type-Options: nosniff
< Date: Wed, 24 May 2017 13:49:15 GMT
< Content-Length: 21
<
color must be unique

[mh-cbon@pc4 ademo] $ curl --data "color=" http://localhost:8080/Create
...
< HTTP/1.1 400 Bad Request
< Content-Type: text/plain; charset=utf-8
< X-Content-Type-Options: nosniff
< Date: Wed, 24 May 2017 13:48:46 GMT
< Content-Length: 24
<
color must not be empty

[mh-cbon@pc4 ademo] $ curl --data "color=green" http://localhost:8080/Create
{"ID":"2","Color":"green"}

[mh-cbon@pc4 ademo] $ curl http://localhost:8080/GetById?id=1
{"ID":"1","Color":"blue"}

[mh-cbon@pc4 ademo] $ curl http://localhost:8080/GetById?id=2
{"ID":"2","Color":"green"}

[mh-cbon@pc4 ademo] $ fg
go run *go
^Csignal: interrupt
```
