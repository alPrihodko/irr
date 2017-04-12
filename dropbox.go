package main

import (
	"irrigation/irr"

	"github.com/stacktic/dropbox"
)

/*
DB - dropbox database
*/
var DB *dropbox.Dropbox

func init() {

	// 1. Create a new dropbox object.
	DB = dropbox.NewDropbox()

	// 2. Provide your clientid and clientsecret (see prerequisite).
	DB.SetAppInfo(irr.DropboxClientid, irr.DropboxClientsecret)

	// 3. Provide the user token.
	DB.SetAccessToken(irr.DropboxToken)

}
