package userdb

import (
	"bufio"
	"database/sql"
	"errors"
	"gotti/utils"
	"os"

	_ "github.com/mattn/go-sqlite3"
)

type UserDatabase struct{
    DbPath string
    DB *sql.DB
}

type UserRecord struct {
    username string
    registered int
}

func (users *UserDatabase)InitializeDB(file string) error{
    var err error
    users.DB, err = sql.Open("sqlite3", users.DbPath)
    if err != nil{
        return err
    }
    _, err = users.DB.Exec(`CREATE TABLE IF NOT EXISTS "REGISTRATION" ("ID" TEXT PRIMARY KEY, "REGISTERED" INTEGER)`)
    if err != nil{
        return err
    }
    err = users.copyUsersFromFile(file)
    if err != nil{
        return err
    }
    return nil
}

func (users *UserDatabase)existUser(username string) (bool, error){
    res, err := users.DB.Query(`SELECT COUNT(*) FROM REGISTRATION WHERE ID = ?`,username)
    if err!=nil{
        return false,err
    }
    if res == nil{
        return false,errors.New("result not found")
    }
    return true,nil
}

func (users *UserDatabase)RegisterUser(username string) error{
    res, err := users.DB.Exec(`UPDATE "REGISTRATION" SET REGISTERED = 1 WHERE ID = ?`, username)
    if err != nil{
        return err
    }
    if res == nil{
        return errors.New("result not found")
    }
    return nil
}

func (users *UserDatabase)IsRegisteredUser(username string)(bool, error){
    res, err := users.DB.Query(`SELECT * FROM "REGISTRATION" WHERE ID = ?`, username)
    defer res.Close()
    if err != nil{
        return true,err
    }
    var buf UserRecord
    res.Next()
    res.Scan(&buf.username, &buf.registered)
    if buf.username==""{
        return true,errors.New("result not found")
    }
    if buf.registered==1{
        return true,nil
    }
    return false,nil
}

//DBの中身が既に存在する場合は動作を保証しません．
func (users *UserDatabase)copyUsersFromFile(file string) error{
    fp, err := os.Open(file)
    if err != nil{
        return err
    }
    defer fp.Close()
    s := bufio.NewScanner(fp)
    for s.Scan(){
        u := s.Text()
        if !utils.IsProperUsername(u){
            return errors.New("Invalid username:"+u)
        }
        _, err := users.DB.Exec(`INSERT INTO "REGISTRATION" ("ID", "REGISTERED") VALUES (?, ?)`, u, 0)
        if err != nil{
            return err
        }
    }
    return nil
}
