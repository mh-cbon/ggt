package capable

import "fmt"

// Session provide access to the session
type Session struct{}

// GetAll return a map
func (c Session) GetAll(sessionName map[interface{}]interface{}) {
	fmt.Printf(`sessionName %q
    `, sessionName)
}
