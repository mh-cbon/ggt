package capable

// file generated by
// github.com/mh-cbon/ggt
// do not edit

import (
	ggt "github.com/mh-cbon/ggt/lib"
	finder "github.com/mh-cbon/service-finder"
	"io"
	"net/http"
	"strconv"
)

var xx6b84bdd3e515a722878248845f8083b0d9390955 = strconv.Atoi
var xx3efbdeb36d52176f0553035ccfbb8819cdd065e0 = io.Copy
var xx9a840474508691f902be6ac58b01856c4ce5e4af = http.StatusOK

// RestSession is an httper of Session.
// Session provide access to the session
type RestSession struct {
	embed    Session
	Services finder.ServiceFinder
	Log      ggt.HTTPLogger
	Session  ggt.SessionStoreProvider
	Upload   ggt.Uploader
}

// NewRestSession constructs an httper of Session
func NewRestSession(embed Session) *RestSession {
	ret := &RestSession{
		embed:    embed,
		Services: finder.New(),
		Log:      &ggt.VoidLog{},
		Session:  &ggt.VoidSession{},
		Upload:   &ggt.FileProvider{},
	}
	ret.Log.Handle(nil, nil, nil, "constructor", "RestSession")
	return ret
}

// GetAll invoke Session.GetAll using the request body as a json payload.
// GetAll return a map
func (t *RestSession) GetAll(w http.ResponseWriter, r *http.Request) {
	t.Log.Handle(w, r, nil, "begin", "RestSession", "GetAll")
	var sessionName map[interface{}]interface{}
	{

		storesessionName, storeErr := t.Session.Get(r, "name")

		if storeErr != nil {

			t.Log.Handle(w, r, storeErr, "session", "store", "get", "error", "sessionName", "error", "RestSession", "GetAll")
			http.Error(w, storeErr.Error(), http.StatusInternalServerError)

			return
		}

		defer func() {
			saveErr := storesessionName.Save(r, w)

			if saveErr != nil {

				t.Log.Handle(w, r, saveErr, "session", "save", "error", "sessionName", "error", "RestSession", "GetAll")
				http.Error(w, saveErr.Error(), http.StatusInternalServerError)

				return
			}

		}()

		valsessionName, getErr := storesessionName.Get()

		if getErr != nil {

			t.Log.Handle(w, r, getErr, "session", "read", "error", "sessionName", "error", "RestSession", "GetAll")
			http.Error(w, getErr.Error(), http.StatusInternalServerError)

			return
		}

		sessionName = valsessionName

	}

	t.embed.GetAll(sessionName)

	w.WriteHeader(200)

	t.Log.Handle(w, r, nil, "end", "RestSession", "GetAll")

}

// RestSessionDescriptor describe a *RestSession
type RestSessionDescriptor struct {
	ggt.TypeDescriptor
	about        *RestSession
	methodGetAll *ggt.MethodDescriptor
}

// NewRestSessionDescriptor describe a *RestSession
func NewRestSessionDescriptor(about *RestSession) *RestSessionDescriptor {
	ret := &RestSessionDescriptor{about: about}
	ret.methodGetAll = &ggt.MethodDescriptor{
		Name:    "GetAll",
		Handler: about.GetAll,
		Route:   "GetAll",
		Methods: []string{},
	}
	ret.TypeDescriptor.Register(ret.methodGetAll)
	return ret
}

// GetAll returns a MethodDescriptor
func (t *RestSessionDescriptor) GetAll() *ggt.MethodDescriptor { return t.methodGetAll }
