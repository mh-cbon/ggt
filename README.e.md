# {{.Name}}

{{pkgdoc}}

Check the demo [here](https://github.com/mh-cbon/ggt/tree/master/ademo)

# {{toc 5}}

# DDP 2 3R

Domain Driven Programming 2 Remove Redundant Repetition

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

Transform a business controller into an http end point.

#### $ {{exec "ggt" "-help" "http-provider" | color "sh"}}

[here](https://github.com/mh-cbon/ggt/tree/master/http-provider).

## http-consumer

Transform a business controller into an http client.

#### $ {{exec "ggt" "-help" "http-consumer" | color "sh"}}

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
