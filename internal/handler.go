package internal
import (
	"crypto/rand"
    "net/http"
	"fmt"
    "time"
	"gotti/smtpMail"
)
var Tokens = make(map[string]userToken,5)

type userToken struct{
    username string
    ott string
    registered time.Time
}

type APIRegisterHandler struct{
    token string
}

func (a APIRegisterHandler)ServeHTTP(w http.ResponseWriter, r *http.Request) {
    u := r.URL.Query().Get("username")
    t := r.URL.Query().Get("appToken")
    if t!=a.token{
        fmt.Println(t)
        w.WriteHeader(401)
        return
    }
    if u=="" {
        w.WriteHeader(401)
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
    */
    addr := u+"@edu.cc.uec.ac.jp"
    buf := make([]byte,4)
    _, err := rand.Read(buf)
    ott := fmt.Sprintf("%8x",buf)
    if err != nil{
        fmt.Println(err)
    }

    Tokens[u] = userToken{username: u, ott:ott, registered: time.Now()}
    fmt.Println(addr)
    err = smtpMail.Send(addr,ott)
    if err != nil{
        fmt.Println(err)
        return
    }
    fmt.Println("sent")
}

type APIVerifyHandler struct {
    token string
}

func (a APIVerifyHandler)ServeHTTP(w http.ResponseWriter, r *http.Request) {
    u := r.URL.Query().Get("username")
    t := r.URL.Query().Get("appToken")
    o := r.URL.Query().Get("ott")
    if t!=a.token{
        fmt.Println(t)
        w.WriteHeader(401)
        return
    }
    if u=="" {
        w.WriteHeader(401)
        return
    }
    //TODO: database access
    /*
    if (db.user.isRegistered){
        return errors.New("user already registered")
    }
    if (db.user.isExists){
        return errors.New("this user don't have an uec account")
    }
    db.user.isRegisterd=true
    */
    v,ok := Tokens[u]
    if !ok{
        w.WriteHeader(401)
        return
    }
    if v.ott!=o {
        w.WriteHeader(401)
        return
    }
    delete(Tokens,u)
    w.WriteHeader(200)
}
