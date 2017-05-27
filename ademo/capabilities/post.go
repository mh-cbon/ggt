package capable

import "fmt"

// Post ...
type Post struct{}

// GetOne ...
func (c Post) GetOne(postArg1 string) {
	fmt.Printf(`postArg1 %q
    `, postArg1)
}

// GetMany ...
func (c Post) GetMany(postArg1, postArg2 string) {
	fmt.Printf(`postArg1 %q
    `, postArg1)
	fmt.Printf(`postArg2 %q
    `, postArg2)
}

// GetConvertedToInt ...
func (c Post) GetConvertedToInt(postArg1 int) {
	fmt.Printf(`postArg1 %q
    `, postArg1)
}

// GetConvertedToBool ...
func (c Post) GetConvertedToBool(postArg1 bool) {
	fmt.Printf(`postArg1 %v
    `, postArg1)
}

// GetConvertedToSlice ...
func (c Post) GetConvertedToSlice(postArg1 []bool) {
	fmt.Printf(`postArg1 %v
    `, postArg1)
}

// GetMaybe ...
func (c Post) GetMaybe(postArg1 *string) {
	fmt.Printf(`postArg1 %q
    `, postArg1)
}
