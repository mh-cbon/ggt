package capable

import "fmt"

// Get ...
type Get struct{}

// GetAll values in url query as a map of values
//@route getall
func (c Get) GetAll(getValues map[string][]string) {
	fmt.Printf(`getValues %q
    `, getValues)
}

// GetAll2 values in url query as a map of value
//@route getall2
func (c Get) GetAll2(getValues map[string]string) {
	fmt.Printf(`getValues %q
    `, getValues)
}

// GetOne arg from url query
//@route getone
func (c Get) GetOne(getArg1 string) {
	fmt.Printf(`getArg1 %q
    `, getArg1)
}

// GetMany args from url query
//@route getmany
func (c Get) GetMany(getArg1, getArg2 string) {
	fmt.Printf(`getArg1 %q
    `, getArg1)
	fmt.Printf(`getArg2 %q
    `, getArg2)
}

// ConvertToInt an arg from url query
//@route converttoint
func (c Get) ConvertToInt(getArg1 int) {
	fmt.Printf(`getArg1 %q
    `, getArg1)
}

// ConvertToBool an arg from url query
//@route converttobool
func (c Get) ConvertToBool(getArg1 bool) {
	fmt.Printf(`getArg1 %v
    `, getArg1)
}

// ConvertToSlice an arg from url query
//@route converttoslice
func (c Get) ConvertToSlice(getArg1 []bool) {
	fmt.Printf(`getArg1 %v
    `, getArg1)
}

// MaybeGet an arg if it exists in url query.
//@route maybeget
func (c Get) MaybeGet(getArg1 *string) {
	fmt.Printf(`getArg1 %q
    `, getArg1)
}
