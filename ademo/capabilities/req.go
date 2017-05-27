package capable

import "fmt"

// Req ...
type Req struct{}

// GetAll ...
func (c Req) GetAll(reqValues map[string][]string) {
	fmt.Printf(`reqValues %q
    `, reqValues)
}

// GetAll2 ...
func (c Req) GetAll2(reqValues map[string]string) {
	fmt.Printf(`reqValues %q
    `, reqValues)
}

// GetOne ...
func (c Req) GetOne(reqArg1 string) {
	fmt.Printf(`reqArg1 %q
    `, reqArg1)
}

// GetMany ...
func (c Req) GetMany(reqArg1, reqArg2 string) {
	fmt.Printf(`reqArg1 %q
    `, reqArg1)
	fmt.Printf(`reqArg2 %q
    `, reqArg2)
}

// GetConvertedToInt ...
func (c Req) GetConvertedToInt(reqArg1 int) {
	fmt.Printf(`reqArg1 %q
    `, reqArg1)
}

// GetConvertedToBool ...
func (c Req) GetConvertedToBool(reqArg1 bool) {
	fmt.Printf(`reqArg1 %v
    `, reqArg1)
}

// GetConvertedToSlice is impossible as route can not accept []string
// func (c Req) GetConvertedToSlice(reqArg1 []bool) {
// 	fmt.Printf(`reqArg1 %v
//     `, reqArg1)
// }

// GetMaybe ...
func (c Req) GetMaybe(reqArg1 *string) {
	fmt.Printf(`reqArg1 %q
    `, reqArg1)
}
