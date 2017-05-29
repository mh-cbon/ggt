package capable

import (
	"fmt"

	finder "github.com/mh-cbon/service-finder"
)

// Svc ...
type Svc struct{}

// Get the services provider
//@route get
func (c Svc) Get(provider finder.ServiceFinder) {
	fmt.Printf(`provider %q
    `, provider)
}

// GetOne service
//@route getone
func (c Svc) GetOne(svcMail WhateverMailService) {
	fmt.Printf(`svcMail %q
    `, svcMail)
}

// WhateverMailService ...
type WhateverMailService interface {
	Send(from, to, title, body string) error
}
