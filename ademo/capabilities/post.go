package capable

import "fmt"

// Post ...
type Post struct{}

// GetAll values from the form.
func (c Post) GetAll(postValues map[string][]string) {
	fmt.Printf(`postValues %q
    `, postValues)
}

// GetAll2 values from the form.
func (c Post) GetAll2(postValues map[string]string) {
	fmt.Printf(`postValues %q
    `, postValues)
}

// GetOne arg form the form.
func (c Post) GetOne(postArg1 string) {
	fmt.Printf(`postArg1 %q
    `, postArg1)
}

// GetMany args form the form.
func (c Post) GetMany(postArg1, postArg2 string) {
	fmt.Printf(`postArg1 %q
    `, postArg1)
	fmt.Printf(`postArg2 %q
    `, postArg2)
}

// ConvertToInt an arg from the form.
func (c Post) ConvertToInt(postArg1 int) {
	fmt.Printf(`postArg1 %q
    `, postArg1)
}

// ConvertToBool an arg from the form.
func (c Post) ConvertToBool(postArg1 bool) {
	fmt.Printf(`postArg1 %v
    `, postArg1)
}

// ConvertToSlice an arg from the form.
func (c Post) ConvertToSlice(postArg1 []bool) {
	fmt.Printf(`postArg1 %v
    `, postArg1)
}

// MaybeGet an arg if it exists in the form.
func (c Post) MaybeGet(postArg1 *string) {
	fmt.Printf(`postArg1 %q
    `, postArg1)
}
