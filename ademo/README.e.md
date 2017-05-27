# A demo

{{template "badge/godoc" .}}

{{pkgdoc}}

# {{toc 5}}

# The main

#### $ {{cat "main.go" | color "go"}}

# The model

#### $ {{cat "tomate/model.go" | color "go"}}

# The controller

#### $ {{cat "tomate/controller.go" | color "go"}}

# The gen

#### $ {{cat "tomate/gen.go" | color "go"}}

# The code for free

## a backend in-memory

- {{link "tomate/zz_tomatessync.go" "tomate/zz_tomatessync.go"}}
- {{link "tomate/zz_tomates.go" "tomate/zz_tomates.go"}}

## an http rpc implementation

- {{link "tomate/zz_rpccontroller.go" "tomate/zz_rpccontroller.go"}}
- {{link "tomate/zz_rpcclient.go" "tomate/zz_rpcclient.go"}}

## an http rest implementation

- {{link "tomate/zz_restclient.go" "tomate/zz_restclient.go"}}
- {{link "tomate/zz_restcontroller.go" "tomate/zz_restcontroller.go"}}

# The test

#### $ {{shell "sh test.sh" | color "sh" }}
