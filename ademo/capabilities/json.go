package capable

import "fmt"

// JSON ...
type JSON struct{}

// ReadJSONBody ...
func (c JSON) ReadJSONBody(jsonReqBody Whatever) {
	fmt.Printf(`jsonReqBody %v
    `, jsonReqBody)
}

// WriteJSONBody ...
func (c JSON) WriteJSONBody() (jsonResBody Whatever) {
	fmt.Printf(`jsonResBody %q
    `, jsonResBody)
	return
}

// Whatever is your model type
type Whatever struct{}
