package capable

import "fmt"

// Route ...
type Route struct{}

// GetAll values from the route.
func (c Route) GetAll(routeValues map[string]string) {
	fmt.Printf(`routeValues %q
    `, routeValues)
}

// GetOne value from the route.
func (c Route) GetOne(routeArg1 string) {
	fmt.Printf(`routeArg1 %q
    `, routeArg1)
}

// GetMany values from the route.
func (c Route) GetMany(routeArg1, routeArg2 string) {
	fmt.Printf(`routeArg1 %q
    `, routeArg1)
	fmt.Printf(`routeArg2 %q
    `, routeArg2)
}

// ConvertToInt an arg from the route.
func (c Route) ConvertToInt(routeArg1 int) {
	fmt.Printf(`routeArg1 %q
    `, routeArg1)
}

// ConvertToBool an arg from the route.
func (c Route) ConvertToBool(routeArg1 bool) {
	fmt.Printf(`routeArg1 %v
    `, routeArg1)
}

// ConvertToSlice ... is impossible
// func (c Route) ConvertToSlice(routeArg1 []bool) {
// 	fmt.Printf(`routeArg1 %v
//     `, routeArg1)
// }

// MaybeGet an arg from the route if it exists.
func (c Route) MaybeGet(routeArg1 *string) {
	fmt.Printf(`routeArg1 %q
    `, routeArg1)
}
