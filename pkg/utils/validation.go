package utils

import (
	"github.com/track/blogserver/pkg/common"
	"regexp"
	"time"
	"unicode/utf8"
)

//常用的验证工具类
type Validation struct {
}

//仅用于int string
func (that *Validation) Required(v interface{}) bool {
	switch v.(type) {
	case int:
		return v.(int) <= 0
	case string:
		return len(v.(string)) <= 0
	}
	return false
}

//仅用于int  判断不在 min-max 范围
func (that *Validation) Range(value, min, max int) bool {
	return value < min || value > max
}

//时间格式判断
func (that *Validation) Time(v string) bool {
	if _, err := time.Parse(common.Layout, v); err != nil {
		return true
	}
	return false
}

//判断数字字母  必须字母开头
func (that *Validation) NumberAndLetter(v string) bool {
	r := regexp.MustCompile("^[A-Za-z]+\\w+$")
	return r.Match([]byte(v))
}

//判断字符串长度 在某个范围
func (that *Validation) Length(v string, l int) bool {
	return utf8.RuneCountInString(v) <= l
}
func (that *Validation) Email(v string) bool {
	r := regexp.MustCompile("\\w+([-+.]\\w+)*@\\w+([-.]\\w+)*\\.\\w+([-.]\\w+)*")
	return r.Match([]byte(v))
}
