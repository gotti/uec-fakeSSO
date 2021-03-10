package smtpMail

import (
    "encoding/base64"
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

func composeEmailMessage(sender Sender, to,  ott string) string{
    header := make(map[string]string)
    header["Return-Path"] = sender.from
    header["From"] = sender.from
    header["To"] = to
    header["Subject"] = fmt.Sprintf("%s%s%s","=?utf-8?b?",base64.StdEncoding.EncodeToString([]byte(sender.sub)),"?=")
    header["MIME-Version"] = "1.0"
    header["Content-Type"] = "text/plain; charset=\"utf-8\""
    header["Content-Transfer-Encoding"] = "base64"
    message := ""
    for k, v := range header {
        message += fmt.Sprintf("%s: %s\r\n", k, v)
    }
    message += "\r\n" + base64.StdEncoding.EncodeToString([]byte(sender.msg+ott))
    return message
}

func Send(to, ott string) error{
    fmt.Println(sender,to,ott)
    auth := smtp.PlainAuth("",sender.username,sender.password,sender.smtpAddress)
    s := sender.smtpAddress+":"+sender.port
    err := smtp.SendMail(s, auth, sender.from, []string{to}, []byte(composeEmailMessage(sender,to,ott)))
    if err != nil{
        return err
    }
    return nil
}
