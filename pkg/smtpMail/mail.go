package smtpMail

import (
	"fmt"
	"net/smtp"
)

type Sender struct {
    smtpAddress string
    port string
    from string
    username string
    password string
    sub string
    msg string
}
var sender Sender
func Initialize(smtpAddress, port, from, username, password, sub, msg string) {
    sender.smtpAddress = smtpAddress
    sender.port = port
    sender.from = from
    sender.username = username
    sender.password = password
    sender.sub = sub
    sender.msg = msg
}
func Send(to, ott string) error{
    fmt.Println(sender,to,ott)
    auth := smtp.PlainAuth("",sender.username,sender.password,sender.smtpAddress)
    s := sender.smtpAddress+":"+sender.port
    err := smtp.SendMail(s, auth, sender.from, []string{to}, []byte(fmt.Sprintf("SubJect: %s\r\n\r\n%s",sender.sub, sender.msg+ott)))
    if err != nil{
        return err
    }
    return nil
}
