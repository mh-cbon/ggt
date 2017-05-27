package capable

import (
	"fmt"
	"net/http"
)

// Cookie ...
type Cookie struct{}

// GetAll values in cookies
func (c Cookie) GetAll(cookieValues map[string]string) {
	fmt.Printf(`cookieValues %q
    `, cookieValues)
}

// GetAllRaw  cookies
func (c Cookie) GetAllRaw(cookieValues []*http.Cookie) {
	fmt.Printf(`cookieValues %q
    `, cookieValues)
}

// GetOne value form cookies
func (c Cookie) GetOne(cookieWhatever string) {
	fmt.Printf(`cookieWhatever %q
    `, cookieWhatever)
}

// GetOneRaw cookie
func (c Cookie) GetOneRaw(cookieWhatever http.Cookie) {
	fmt.Printf(`cookieWhatever %v
    `, cookieWhatever)
}

// MaybeGetOneRaw cookie
func (c Cookie) MaybeGetOneRaw(cookieWhatever *http.Cookie) {
	fmt.Printf(`cookieWhatever %q
    `, cookieWhatever)
}

// Write a cookie
func (c Cookie) Write() (cookieWhatever http.Cookie) {
	cookieWhatever = http.Cookie{Value: "whatever"}
	fmt.Printf(`cookieWhatever %v
    `, cookieWhatever)
	return cookieWhatever
}

// MaybeDelete a cookie
func (c Cookie) MaybeDelete() (cookieWhatever *http.Cookie) {
	cookieWhatever = nil
	fmt.Printf(`cookieWhatever %q
    `, cookieWhatever)
	return cookieWhatever
}

// Delete a cookie
func (c Cookie) Delete() (cookieWhatever *http.Cookie) {
	return nil
}

// GetMany args from url query
func (c Cookie) GetMany(cookieArg1, cookieArg2 string) {
	fmt.Printf(`cookieArg1 %q
    `, cookieArg1)
	fmt.Printf(`cookieArg2 %q
    `, cookieArg2)
}

// ConvertToInt an arg from url query
func (c Cookie) ConvertToInt(cookieArg1 int) {
	fmt.Printf(`cookieArg1 %q
    `, cookieArg1)
}

// ConvertToBool an arg from url query
func (c Cookie) ConvertToBool(cookieArg1 bool) {
	fmt.Printf(`cookieArg1 %v
    `, cookieArg1)
}

// ConvertToSlice an arg from url query
// func (c Cookie) ConvertToSlice(cookieArg1 []bool) {
// 	fmt.Printf(`cookieArg1 %v
//     `, getArg1)
// }

// MaybeGet an arg if it exists in url query.
func (c Cookie) MaybeGet(cookieArg1 *string) {
	fmt.Printf(`cookieArg1 %q
    `, cookieArg1)
}
