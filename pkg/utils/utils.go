package utils

import (
	"regexp"
)

//弊学の学籍番号かどうか確認 ok: a2010123, ng: abc2010123
func IsProperUsername(u string) bool{
    usernameValidater := regexp.MustCompile(`^[a-z](\d){7}$`)
    return usernameValidater.MatchString(u)
}
