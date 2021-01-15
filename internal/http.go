package internal

import "net/http"

func APIServer(token string){
    http.Handle("/register", APIRegisterHandler{token})
    http.Handle("/verify", APIVerifyHandler{token})
    http.ListenAndServe(":8084",nil)
}
