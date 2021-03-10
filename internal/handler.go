package internal
import (
	"crypto/rand"
	"fmt"
	"gotti/smtpMail"
    "gotti/utils"
    "gotti/userdb"
	"net/http"
	"sync"
	"time"
)
type typeSafeTokens struct{
    mu sync.Mutex
    Tokens map[string](userToken)
}
var SafeTokens = typeSafeTokens{Tokens:make(map[string](userToken),10)}

type userToken struct{
    username string
    ott string
    registered time.Time
}

type APIRegisterHandler struct{
    token string
    db userdb.UserDatabase
}

func (a APIRegisterHandler)ServeHTTP(w http.ResponseWriter, r *http.Request) {
    u := r.URL.Query().Get("username")
    t := r.URL.Query().Get("appToken")
    if t!=a.token || !utils.IsProperUsername(u){
        fmt.Println(t,u,utils.IsProperUsername(u))
        w.WriteHeader(401)
        w.Write([]byte("Improper username format or invalid appToken"))
        return
    }
    b,e := a.db.IsRegisteredUser(u)
    if e != nil{
        w.WriteHeader(401)
        w.Write([]byte("This user does not exist on our database"))
        return
    }
    if b{
        w.WriteHeader(401)
        w.Write([]byte("This user is already registered"))
        return
    }
    addr := u+"@edu.cc.uec.ac.jp"
    buf := make([]byte,4)
    _, err := rand.Read(buf)
    ott := fmt.Sprintf("%8x",buf)
    if err != nil{
        fmt.Println(err)
    }

    //Handlerはgoroutineで呼び出されるのでスレッドセーフでないスライスはロック
    SafeTokens.mu.Lock()
    SafeTokens.Tokens[u] = userToken{username: u, ott:ott, registered: time.Now()}
    SafeTokens.mu.Unlock()

    fmt.Println(addr)
    err = smtpMail.Send(addr,ott)
    if err != nil{
        fmt.Println(err)
        return
    }
    fmt.Println("sent")
    w.WriteHeader(200)
    w.Write([]byte("success"))
}

type APIVerifyHandler struct {
    token string
    db userdb.UserDatabase
}

func (a APIVerifyHandler)ServeHTTP(w http.ResponseWriter, r *http.Request) {
    u := r.URL.Query().Get("username")
    t := r.URL.Query().Get("appToken")
    o := r.URL.Query().Get("ott")
    if t!=a.token || !utils.IsProperUsername(u){
        fmt.Println(t,u)
        w.WriteHeader(401)
        w.Write([]byte("Improper username format or invalid appToken"))
        return
    }
    SafeTokens.mu.Lock()
    v,ok := SafeTokens.Tokens[u]
    if !ok || v.ott!=o{
        w.WriteHeader(401)
        w.Write([]byte("Invalid one time token"))
    } else {
        w.WriteHeader(200)
        w.Write([]byte("success"))
        fmt.Println(a.db.RegisterUser(u))
        fmt.Println(a.db.IsRegisteredUser(u))
        delete(SafeTokens.Tokens,u)
    }
    SafeTokens.mu.Unlock()
    return
}
