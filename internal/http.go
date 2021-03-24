package internal

import (
	"gotti/userdb"
	"log"
	"net/http"
	"os"
)

func APIServer(token string, users userdb.UserDatabase){
    http.Handle("/register", APIRegisterHandler{token,users})
    http.Handle("/verify", APIVerifyHandler{token,users})
    err := http.ListenAndServe(":8084",nil)
    if err != nil{
        log.Fatal(err)
        os.Exit(1)
    }
}
