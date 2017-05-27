package capable

import "fmt"

// Route ...
type Route struct{}

// GetAll ...
func (c Route) GetAll(routeValues map[string]string) {
	fmt.Printf(`routeValues %q
    `, routeValues)
}

// GetOne ...
func (c Route) GetOne(routeArg1 string) {
	fmt.Printf(`routeArg1 %q
    `, routeArg1)
}

// GetMany ...
func (c Route) GetMany(routeArg1, routeArg2 string) {
	fmt.Printf(`routeArg1 %q
    `, routeArg1)
	fmt.Printf(`routeArg2 %q
    `, routeArg2)
}

// GetConvertedToInt ...
func (c Route) GetConvertedToInt(routeArg1 int) {
	fmt.Printf(`routeArg1 %q
    `, routeArg1)
}

// GetConvertedToBool ...
func (c Route) GetConvertedToBool(routeArg1 bool) {
	fmt.Printf(`routeArg1 %v
    `, routeArg1)
}

// GetConvertedToSlice ... is impossible
// func (c Route) GetConvertedToSlice(routeArg1 []bool) {
// 	fmt.Printf(`routeArg1 %v
//     `, routeArg1)
// }

// GetMaybe ...
func (c Route) GetMaybe(routeArg1 *string) {
	fmt.Printf(`routeArg1 %q
    `, routeArg1)
}
