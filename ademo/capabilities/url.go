package capable

import "fmt"

// URL ...
type URL struct{}

// GetOne ...
func (c URL) GetOne(urlArg1 string) {
	fmt.Printf(`urlArg1 %q
    `, urlArg1)
}

// GetMany ...
func (c URL) GetMany(urlArg1, urlArg2 string) {
	fmt.Printf(`urlArg1 %q
    `, urlArg1)
	fmt.Printf(`urlArg2 %q
    `, urlArg2)
}

// GetConvertedToInt ...
func (c URL) GetConvertedToInt(urlArg1 int) {
	fmt.Printf(`urlArg1 %q
    `, urlArg1)
}

// GetConvertedToBool ...
func (c URL) GetConvertedToBool(urlArg1 bool) {
	fmt.Printf(`urlArg1 %v
    `, urlArg1)
}

// GetConvertedToSlice is impossible as route can not accept []string
// func (c URL) GetConvertedToSlice(urlArg1 []bool) {
// 	fmt.Printf(`urlArg1 %v
//     `, urlArg1)
// }

// GetMaybe ...
func (c URL) GetMaybe(urlArg1 *string) {
	fmt.Printf(`urlArg1 %q
    `, urlArg1)
}
