package capable

import "fmt"

// Req is a merge of route, url, form
type Req struct{}

// GetAll return a merged map of route, url, form
func (c Req) GetAll(reqValues map[string][]string) {
	fmt.Printf(`reqValues %q
    `, reqValues)
}

// GetAll2 return a merged map of route, url, form
func (c Req) GetAll2(reqValues map[string]string) {
	fmt.Printf(`reqValues %q
    `, reqValues)
}

// GetOne return the first value in route, url, form
func (c Req) GetOne(reqArg1 string) {
	fmt.Printf(`reqArg1 %q
    `, reqArg1)
}

// GetMany return the first value of each parameter in route, url, form
func (c Req) GetMany(reqArg1, reqArg2 string) {
	fmt.Printf(`reqArg1 %q
    `, reqArg1)
	fmt.Printf(`reqArg2 %q
    `, reqArg2)
}

// ConvertToInt an arg
func (c Req) ConvertToInt(reqArg1 int) {
	fmt.Printf(`reqArg1 %q
    `, reqArg1)
}

// ConvertToBool an arg
func (c Req) ConvertToBool(reqArg1 bool) {
	fmt.Printf(`reqArg1 %v
    `, reqArg1)
}

// ConvertToSlice is impossible as route can not accept []string
// func (c Req) ConvertToSlice(reqArg1 []bool) {
// 	fmt.Printf(`reqArg1 %v
//     `, reqArg1)
// }

// MaybeGet an arg if it exists.
func (c Req) MaybeGet(reqArg1 *string) {
	fmt.Printf(`reqArg1 %q
    `, reqArg1)
}
