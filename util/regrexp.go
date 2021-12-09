package util

import (
	"regexp"
	"strconv"
)

func CheckMessageCode(str string, len int) string {
	reg := regexp.MustCompile("\\d{" + strconv.Itoa(len) + "}")
	return  reg.FindString(str)
}

