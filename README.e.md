# {{.Name}}

{{pkgdoc}}

Check the demo [here](https://github.com/mh-cbon/ggt/tree/master/ademo)

# {{toc 5}}

# Install

#### go
{{template "go/install" .}}

# Usage

#### $ {{exec "ggt" "-help" | color "sh"}}

# Toolbox

## slicer

Create a typed slice of a struct, [here](https://github.com/mh-cbon/ggt/tree/master/slicer).

#### $ {{exec "ggt" "-help" "slicer" | color "sh"}}

## chaner / mutexer

Mutex a type so its access are sync and protected of race conditions.

#### $ {{exec "ggt" "-help" "chaner" | color "sh"}}
[chaner](https://github.com/mh-cbon/ggt/tree/master/chaner)

#### $ {{exec "ggt" "-help" "mutexer" | color "sh"}}
[mutexer](https://github.com/mh-cbon/ggt/tree/master/mutexer)

## http-provider

Transform a business controller into an http end point, [here](https://github.com/mh-cbon/ggt/tree/master/http-provider).

## http-clienter

Transform a business controller into an http client, tbd.


# notes

This repository is the refactoring and rationalization of those repos,

- https://github.com/mh-cbon/http-clienter
- https://github.com/mh-cbon/goriller
- https://github.com/mh-cbon/httper
- https://github.com/mh-cbon/jsoner
- https://github.com/mh-cbon/channeler
- https://github.com/mh-cbon/mutexer
- https://github.com/mh-cbon/lister
- https://github.com/mh-cbon/astutil

Please come back later for complete port, watch it if you are interested.

# dev

```sh
go install && go generate ./ademo/*go
go run ademo/*go &
curl http://localhost:8080/GetById?id=2
curl http://localhost:8080/GetById?id=0
curl --data "color=blue" http://localhost:8080/Create
curl -H "Content-Type: application/json" -X POST -d '{"color":""}' http://localhost:8080/Create
curl -H "Content-Type: application/json" -X POST -d '{"color":"blue"}' http://localhost:8080/Create
```
