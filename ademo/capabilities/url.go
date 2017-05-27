package capable

import "fmt"

// URL is a merge of route, url
type URL struct{}

// GetAll  return a merged map of route, url
func (c URL) GetAll(urlValues map[string][]string) {
	fmt.Printf(`urlValues %q
    `, urlValues)
}

// GetAll2 return a merged map of route, url
func (c URL) GetAll2(urlValues map[string]string) {
	fmt.Printf(`urlValues %q
    `, urlValues)
}

// GetOne return the first value in route, url
func (c URL) GetOne(urlArg1 string) {
	fmt.Printf(`urlArg1 %q
    `, urlArg1)
}

// GetMany return the first value of each parameter in route, url
func (c URL) GetMany(urlArg1, urlArg2 string) {
	fmt.Printf(`urlArg1 %q
    `, urlArg1)
	fmt.Printf(`urlArg2 %q
    `, urlArg2)
}

// ConvertToInt an arg
func (c URL) ConvertToInt(urlArg1 int) {
	fmt.Printf(`urlArg1 %q
    `, urlArg1)
}

// ConvertToBool an arg
func (c URL) ConvertToBool(urlArg1 bool) {
	fmt.Printf(`urlArg1 %v
    `, urlArg1)
}

// ConvertToSlice is impossible as route can not accept []string
// func (c URL) ConvertToSlice(urlArg1 []bool) {
// 	fmt.Printf(`urlArg1 %v
//     `, urlArg1)
// }

// MaybeGet an arg if it exists.
func (c URL) MaybeGet(urlArg1 *string) {
	fmt.Printf(`urlArg1 %q
    `, urlArg1)
}
