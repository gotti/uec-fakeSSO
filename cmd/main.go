package main

import (
	"fmt"
	"gotti/internal"
	"gotti/smtpMail"
    "gotti/userdb"

	"log"
	"github.com/BurntSushi/toml"
)
type Config struct{
    APIServerToken string `toml:"APIServerToken"`
    OTTExpire int `toml:"OTTExpire"`
    UserListDatabase string `toml:"UserListDatabase"`
    UserListFile string `toml:"UserListFile"`
    MailConfig Mail `toml:"Mail"`
}
type Mail struct{
        SmtpAddress string `toml:"SmtpAddress"`
        Port string `toml:"Port"`
        From string `toml:"From"`
        Username string `toml:"Username"`
        Password string `toml:"Password"`
        Sub string `toml:"Sub"`
        Msg string `toml:"Msg"`
}

var users userdb.UserDatabase

var Token *string
var ServerConfig Config
func main(){
    _,err := toml.DecodeFile("./config.toml",&ServerConfig)
    if err != nil{
        log.Fatal(err)
    }
    fmt.Println(ServerConfig)
    internal.InitializeGC(ServerConfig.OTTExpire)
    Token = &ServerConfig.APIServerToken
    users.DbPath = ServerConfig.UserListDatabase
    users.InitializeDB(ServerConfig.UserListFile)
    smtpMail.Initialize(ServerConfig.MailConfig.SmtpAddress, ServerConfig.MailConfig.Port, ServerConfig.MailConfig.From, ServerConfig.MailConfig.Username, ServerConfig.MailConfig.Password, ServerConfig.MailConfig.Sub, ServerConfig.MailConfig.Msg)
    internal.APIServer(ServerConfig.APIServerToken, users)
}
