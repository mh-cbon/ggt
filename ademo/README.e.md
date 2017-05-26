# A demo

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
#### $ {{cat "tomate/zz_tomatessync.go" | color "go"}}
#### $ {{cat "tomate/zz_tomates.go" | color "go"}}

## an http rpc implementation
#### $ {{cat "tomate/zz_rpccontroller.go" | color "go"}}
#### $ {{cat "tomate/zz_rpcclient.go" | color "go"}}

## an http rest implementation
#### $ {{cat "tomate/zz_restcontroller.go" | color "go"}}
#### $ {{cat "tomate/zz_restclient.go" | color "go"}}

# The test

#### $ {{shell "sh test.sh" | color "sh" }}
