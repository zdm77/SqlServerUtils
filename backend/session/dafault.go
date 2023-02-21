package session

import (
	"encoding/gob"
	"fmt"
	"github.com/gorilla/sessions"
	"log"
	"net/http"
	"os"
	"path"
	"sqlutils/backend/model"
)

const (
	cookieName = "utils-sql-server-sessions"
	defaultKey = "default"
)

var (
	key         = []byte("f9abac83a0db40d2238d09fc22d0fce4")
	sessionPath = path.Join(".", "utils-sql-server-sessions")
	Store       = sessions.NewFilesystemStore(sessionPath, key)
)

func init() {
	err := os.MkdirAll(sessionPath, 0755)
	if err != nil {
		log.Panicf("unable to create session directory: %s\n", err)
	}
}
func GetSessionData(request *http.Request) *model.User {
	session, err := Store.Get(request, cookieName)
	if err != nil {

	}
	result := session.Values[defaultKey]
	if result != nil {
		return result.(*model.User)
	}
	return nil
}

func Save(val interface{}, writer http.ResponseWriter, request *http.Request) error {
	gob.Register(&model.User{})
	sess, err := Store.Get(request, cookieName)

	sess.Values[defaultKey] = val
	err = sess.Save(request, writer)
	if err != nil {
		return fmt.Errorf("unable to store the session with %s ssid: %s", sess.ID, err)
	}

	return nil
}
