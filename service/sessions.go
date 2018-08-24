package service

import (
	"github.com/gorilla/sessions"
)

var store *sessions.FilesystemStore

func init() {
	store = sessions.NewFilesystemStore("./data", []byte("something-very-secret"))
	store.Options = &sessions.Options{
		MaxAge: 86400 * 1,
	}
}

//SessionStore this is public func provide session store for handler
func SessionStore() sessions.Store {
	return store
}
