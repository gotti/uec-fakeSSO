package internal
import (
	"crypto/rand"
	"fmt"
	"gotti/smtpMail"
	"net/http"
	"regexp"
	"sync"
	"time"
)
type SafeTokens struct{
    mu sync.Mutex
    Tokens map[string](userToken)
}
var safeTokens SafeTokens

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
    if t!=a.token || isProperUsername(u){
        fmt.Println(t,u)
        w.WriteHeader(401)
        return
    }
    //TODO: database access
    /* 疑似コードは下
    if (db.user.isRegistered){
        w.WriteHeader(401)
        return errors.New("user already registered")
    }
    if (db.user.isExists){
        w.WriteHeader(401)
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

    //Handlerはgoroutineで呼び出されるのでスレッドセーフでないスライスはロック
    safeTokens.mu.Lock()
    safeTokens.Tokens[u] = userToken{username: u, ott:ott, registered: time.Now()}
    safeTokens.mu.Unlock()

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
    if t!=a.token || !isProperUsername(u){
        fmt.Println(t,u)
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
    safeTokens.mu.Lock()
    v,ok := safeTokens.Tokens[u]
    if !ok || v.ott!=o{
        w.WriteHeader(401)
        return
    } else {
        w.WriteHeader(200)
        delete(safeTokens.Tokens,u)
        return
    }
    safeTokens.mu.Unlock()

}

//弊学の学籍番号かどうか確認 ok: a2010123, ng: abc2010123
func isProperUsername(u string) bool{
    usernameValidater := regexp.MustCompile(`[a-z]\d{7}`)
    return usernameValidater.MatchString(u)
}
