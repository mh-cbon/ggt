package tester

import "testing"

//go:generate ggt http-provider GetTest:RestGetTest

// GetTest ...
type GetTest struct {
	T              *testing.T
	WantValueParam string
	WantRefParam   *string
}

//ValueParam ...
func (g GetTest) ValueParam(getParam string) {
	wanted := g.WantValueParam
	if getParam != wanted {
		g.T.Errorf("want=%q got=%q", wanted, getParam)
	}
}

//RefParam ...
func (g GetTest) RefParam(getParam *string) {
	wanted := g.WantRefParam
	if wanted == nil && getParam != nil {
		g.T.Errorf("want=%q got=%q", wanted, getParam)
	} else if wanted != nil && getParam == nil {
		g.T.Errorf("want=%q got=%q", wanted, getParam)
	} else if getParam != nil && wanted != nil && *getParam != *wanted {
		g.T.Errorf("want=%q got=%q", *wanted, *getParam)
	}
}

//MultipleParam ...
func (g GetTest) MultipleParam(getParam *string, getID string) {

}
