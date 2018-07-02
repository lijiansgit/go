package libs

import (
	"time"
	"strings"
)

func TimestampToStr(timestamp int) (str string) {
	str = time.Unix(int64(timestamp), 0).Format("2006-01-02 15:04:05")
	return str
}

func N2br(str string) string {
	replacer := strings.NewReplacer("\n", "<br/>")
	return replacer.Replace(str)
}

func Add(a int) int {
	return a + 1
}

func Reduce(a int) int {
	return a - 1
}