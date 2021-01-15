package internal

import (
    "time"
)

func InitializeGC(e int){
    go func() {
        t := time.NewTicker(60 * time.Second)
        for {
            select {
            case <-t.C:
                gerbaseCollect(e)
            }
        }
    }()
}

func gerbaseCollect(e int){
    now := time.Now()
    for k,v := range Tokens {
        if now.After(v.registered.Add(time.Duration(e)*time.Minute)) {
            delete(Tokens,k)
        }
    }
}
