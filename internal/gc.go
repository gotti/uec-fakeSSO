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
                gerbageCollect(e)
            }
        }
    }()
}

func gerbageCollect(e int){
    now := time.Now()
    SafeTokens.mu.Lock()
    for k,v := range SafeTokens.Tokens {
        if now.After(v.registered.Add(time.Duration(e)*time.Minute)) {
            delete(SafeTokens.Tokens, k)
        }
    }
    SafeTokens.mu.Unlock()
}
