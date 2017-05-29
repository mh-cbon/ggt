package capable

import (
	"context"
	"fmt"
)

// Ctx ...
type Ctx struct{}

// Get the context
//@route get
func (c Ctx) Get(whatever context.Context) {
	fmt.Printf(`whatever %q
    `, whatever)
}

// GetOne arg from context.
//@route getone
func (c Ctx) GetOne(ctxArg1 Whatever) {
	fmt.Printf(`ctxArg1 %q
    `, ctxArg1)
}

// MaybeGetOne arg from the context
//@route maybegetone
func (c Ctx) MaybeGetOne(ctxArg1 *Whatever) {
	fmt.Printf(`ctxArg1 %q
    `, ctxArg1)
}

// SetOne arg on the context.
//@route setone
func (c Ctx) SetOne() (ctxArg1 Whatever) {
	fmt.Printf(`ctxArg1 %q
    `, ctxArg1)
	return
}

// MaybeSetOne arg on the context
//@route maybesetone
func (c Ctx) MaybeSetOne() (ctxArg1 *Whatever) {
	fmt.Printf(`ctxArg1 %q
    `, ctxArg1)
	return
}
