package capable

import "fmt"

// Get ...
type Get struct{}

// GetAll ...
func (c Get) GetAll(getValues map[string][]string) {
	fmt.Printf(`getValues %q
    `, getValues)
}

// GetAll2 ...
func (c Get) GetAll2(getValues map[string]string) {
	fmt.Printf(`getValues %q
    `, getValues)
}

// GetOne ...
func (c Get) GetOne(getArg1 string) {
	fmt.Printf(`getArg1 %q
    `, getArg1)
}

// GetMany ...
func (c Get) GetMany(getArg1, getArg2 string) {
	fmt.Printf(`getArg1 %q
    `, getArg1)
	fmt.Printf(`getArg2 %q
    `, getArg2)
}

// GetConvertedToInt ...
func (c Get) GetConvertedToInt(getArg1 int) {
	fmt.Printf(`getArg1 %q
    `, getArg1)
}

// GetConvertedToBool ...
func (c Get) GetConvertedToBool(getArg1 bool) {
	fmt.Printf(`getArg1 %v
    `, getArg1)
}

// GetConvertedToSlice ...
func (c Get) GetConvertedToSlice(getArg1 []bool) {
	fmt.Printf(`getArg1 %v
    `, getArg1)
}

// GetMaybe ...
func (c Get) GetMaybe(getArg1 *string) {
	fmt.Printf(`getArg1 %q
    `, getArg1)
}
