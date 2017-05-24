# ggt

ggt's generator toolbox


Check the demo [here](https://github.com/mh-cbon/ggt/tree/master/ademo)

# TOC
- [Install](#install)
  - [go](#go)
- [Usage](#usage)
  - [$ ggt -help](#-ggt--help)
- [Toolbox](#toolbox)
  - [slicer](#slicer)
    - [$ ggt -help slicer](#-ggt--help-slicer)
  - [chaner / mutexer](#chaner--mutexer)
    - [$ ggt -help chaner](#-ggt--help-chaner)
    - [$ ggt -help mutexer](#-ggt--help-mutexer)
  - [http-provider](#http-provider)
  - [http-clienter](#http-clienter)
- [notes](#notes)
- [dev](#dev)

# Install

#### go
```sh
go get github.com/mh-cbon/ggt
```

# Usage

#### $ ggt -help
```sh
ggt - 0.0.0
```

# Toolbox

## slicer

Create a typed slice of a struct, [here](https://github.com/mh-cbon/ggt/tree/master/slicer).

#### $ ggt -help slicer
```sh
ggt [options] slicer ...[FromTypeName:ToTypeName]

generates typed slice

[options]
	see ggt -help

...[FromTypeName:ToTypeName]
	A list of types such as src:dst.
	A type is defined by its package path and its type name,
	[pkgpath/]name.
	If the Package path is empty, it is set to the package name being generated.
	Name can be a valid type identifier such as TypeName, *TypeName, []TypeName

Example
	ggt -c slicer MySrcType:gen/*NewGenType
	ggt -c slicer myModule/*MySrcType:gen/NewGenType
```

## chaner / mutexer

Mutex a type so its access are sync and protected of race conditions.

#### $ ggt -help chaner
```sh
sh
```
[chaner](https://github.com/mh-cbon/ggt/tree/master/chaner)

#### $ ggt -help mutexer
```sh
sh
```
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
