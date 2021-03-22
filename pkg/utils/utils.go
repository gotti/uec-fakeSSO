package utils

import (
	"regexp"
)

//弊学の学籍番号かどうか確認 ok: a2010123, ng: abc2010123
func IsProperUsername(u string) bool{
    usernameValidater := regexp.MustCompile(`^[a-z](\d){7}$`)
    extraValidater := regexp.MustCompile(`^[a-z](\d){3}h(\d){3}$`)
    return usernameValidater.MatchString(u) || extraValidater.MatchString(u)
}
