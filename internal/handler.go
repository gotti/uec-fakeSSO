package internal
import (
	"crypto/rand"
    "net/http"
	"fmt"
	"gotti/smtpMail"
)
var Tokens = make(map[string]string,5)

type APIRegisterHandler struct{
    token string
}

func (a APIRegisterHandler)ServeHTTP(w http.ResponseWriter, r *http.Request) {
    u := r.URL.Query().Get("username")
    t := r.URL.Query().Get("appToken")
    fmt.Println(a.token)
    if t!=a.token{
        fmt.Println(t)
        return
    }
    if u=="" {
        return
    }
    //TODO: database access
    /* 疑似コードは下
    if (db.user.isRegistered){
        return errors.New("user already registered")
    }
    if (db.user.isExists){
        return errors.New("this user don't have an uec account")
    }
    db.user.isExists=true
    */
    addr := u+"@edu.cc.uec.ac.jp"
    buf := make([]byte,4)
    _, err := rand.Read(buf)
    if err != nil{
        fmt.Println(err)
    }
    fmt.Println(addr)
    err = smtpMail.Send(addr,fmt.Sprintf("%8x",buf))
    if err != nil{
        fmt.Println(err)
        return
    }
    fmt.Println("sent")
}
