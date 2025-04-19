package util

import (
	"strings"
)

func MaskEmail(email string) string {
	strings.LastIndex(email, "@")
	id := email[:strings.LastIndex(email, "@")]
	domain := email[strings.LastIndex(email, "@"):]

	switch len(id) {
	case 2:
		id = id[:1] + "*"
	case 3:
		id = id[:1] + "*" + id[2:]
	case 4:
		id = id[:1] + "**" + id[3:]
	default:
		masks := strings.Repeat("*", len(id)-4)
		id = id[0:2] + masks + id[len(id)-2:]
	}

	return id + domain
}

func MaskPhone(phone string) string {
	return phone[:3] + "****" + phone[7:]
}

func MaskRealName(realName string) string {
	runeRealName := []rune(realName)
	if len(runeRealName) == 2 {
		return string(runeRealName[:1]) + "*"
	}

	return string(runeRealName[:1]) + strings.Repeat("*", len(runeRealName)-2) + string(runeRealName[len(runeRealName)-1])
}

func MaskLoginName(loginName string) string {
	if strings.Contains(loginName, "@") {
		return MaskEmail(loginName)
	}

	return MaskPhone(loginName)
}
