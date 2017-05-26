# ggt

ggt's generator toolbox


Check the demo [here](https://github.com/mh-cbon/ggt/tree/master/ademo)

# TOC
- [3R](#3r)
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
    - [$ ggt -help http-provider](#-ggt--help-http-provider)
  - [http-consumer](#http-consumer)
    - [$ ggt -help http-consumer](#-ggt--help-http-consumer)
- [notes](#notes)

# 3R

Remove Redundant Repetition

# Install

#### go
```sh
go get github.com/mh-cbon/ggt
```

# Usage

#### $ ggt -help
```sh
ggt - 0.0.0
    ggt [options] [generator] [...types]

ggt's generator toolbox

[options]
    -help        Show help
    -version     Show version
    -vv          More verbose
    -mode        Generator mode when suitable (rpc|route).

[generator]

    One of slicer, chaner, mutexer, http-provider, http-consumer.

[...types]
    A list of types such as src:dst.
    A type is defined by its package path and its type name,
    [pkgpath/]name.
    If the Package path is empty, it is set to the package name being generated.
    Name can be a valid type identifier such as TypeName, *TypeName, []TypeName

Example

    ggt -c slicer MySrcType:gen/*NewGenType
    ggt -c slicer myModule/*MySrcType:gen/NewGenType
```

# Toolbox

## slicer

Create a typed slice of a struct, [here](https://github.com/mh-cbon/ggt/tree/master/slicer).

#### $ ggt -help slicer
```sh
ggt [options] slicer ...[FromTypeName:ToTypeName]

generates typed slice

[options]

    -c        Create a contract of the generated type.
    -p        Force out package name

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
ggt [options] chaner ...[FromTypeName:ToTypeName]

generates race protected type

[options]

    -p        Force out package name

...[FromTypeName:ToTypeName]

    A list of types such as src:dst.
    A type is defined by its package path and its type name,
    [pkgpath/]name.
    If the Package path is empty, it is set to the package name being generated.
    Name can be a valid type identifier such as TypeName, *TypeName, []TypeName

Example

    ggt -c chaner MySrcType:gen/*NewGenType
    ggt -c chaner myModule/*MySrcType:gen/NewGenType
```

[chaner](https://github.com/mh-cbon/ggt/tree/master/chaner)

#### $ ggt -help mutexer
```sh
ggt [options] mutexer ...[FromTypeName:ToTypeName]

generates race protected type

[options]

    -p        Force out package name

...[FromTypeName:ToTypeName]

    A list of types such as src:dst.
    A type is defined by its package path and its type name,
    [pkgpath/]name.
    If the Package path is empty, it is set to the package name being generated.
    Name can be a valid type identifier such as TypeName, *TypeName, []TypeName

Example

    ggt -c mutexer MySrcType:gen/*NewGenType
    ggt -c mutexer myModule/*MySrcType:gen/NewGenType
```

[mutexer](https://github.com/mh-cbon/ggt/tree/master/mutexer)

## http-provider

Transform a business controller into an http end point.

#### $ ggt -help http-provider
```sh
ggt [options] http-provider ...[FromTypeName:ToTypeName]

generates http oriented implementation of given type.

[options]

    -p        Force out package name
    -mode     Generation mode (rpc|route).

...[FromTypeName:ToTypeName]

    A list of types such as src:dst.
    A type is defined by its package path and its type name,
    [pkgpath/]name.
    If the Package path is empty, it is set to the package name being generated.
    Name can be a valid type identifier such as TypeName, *TypeName, []TypeName

Example

    ggt -c http-provider MySrcType:gen/*NewGenType
    ggt -c http-provider myModule/*MySrcType:gen/NewGenType
```

[here](https://github.com/mh-cbon/ggt/tree/master/http-provider).

## http-consumer

Transform a business controller into an http client.

#### $ ggt -help http-consumer
```sh
ggt [options] http-consumer ...[FromTypeName:ToTypeName]

generates http client implementation of given type.

[options]

    -p        Force out package name
    -mode     Thep referred generation mode (rpc|route)

...[FromTypeName:ToTypeName]

    A list of types such as src:dst.
    A type is defined by its package path and its type name,
    [pkgpath/]name.
    If the Package path is empty, it is set to the package name being generated.
    Name can be a valid type identifier such as TypeName, *TypeName, []TypeName

Example

    ggt http-consumer MySrcType:gen/*NewGenType
    ggt http-consumer myModule/*MySrcType:gen/NewGenType
    ggt -mode rpc http-consumer myModule/*MySrcType:gen/NewGenType
    ggt -mode route http-consumer myModule/*MySrcType:gen/NewGenType
```

[here](https://github.com/mh-cbon/ggt/tree/master/http-consumer).


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
