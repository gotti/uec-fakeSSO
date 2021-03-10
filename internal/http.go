package internal

import (
	"gotti/userdb"
	"net/http"
)

func APIServer(token string, users userdb.UserDatabase){
    http.Handle("/register", APIRegisterHandler{token,users})
    http.Handle("/verify", APIVerifyHandler{token,users})
    http.ListenAndServe(":8084",nil)
}
