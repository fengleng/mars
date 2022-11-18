package utils

import (
	"bytes"
	"hash/fnv"
	"strings"
	"unicode"

	"github.com/satori/go.uuid"
)

func RemoveSuffixIfMatched(s string, suffix string) string {
	if strings.HasSuffix(s, suffix) {
		return s[0 : len(s)-len(suffix)]
	}
	return s
}
func HashStr(s string) uint32 {
	f := fnv.New32a()
	_, _ = f.Write([]byte(s))
	return f.Sum32()
}
func GenUUID() (string, error) {
	u := uuid.NewV4()
	return strings.Replace(u.String(), "-", "", -1), nil
}
func UnderScore2Camel(name string) string {
	var buf []byte
	toggleUpper := true
	for i := 0; i < len(name); i++ {
		if name[i] == '_' {
			toggleUpper = true
		} else {
			c := name[i]
			if toggleUpper {
				toggleUpper = false
				if c >= 'a' && c <= 'z' {
					c = c - 'a' + 'A'
				}
			}
			if c >= '0' && c <= '9' {
				toggleUpper = true
			}
			buf = append(buf, c)
		}
	}
	return string(buf)
}

// Camel2UnderScore 驼峰转下划线
func Camel2UnderScore(name string) string {
	var posList []int
	i := 1
	for i < len(name) {
		if name[i] >= 'A' && name[i] <= 'Z' {
			posList = append(posList, i)
			i++
			for i < len(name) && name[i] >= 'A' && name[i] <= 'Z' {
				i++
			}
		} else {
			i++
		}
	}
	lower := strings.ToLower(name)
	if len(posList) == 0 {
		return lower
	}
	b := strings.Builder{}
	left := 0
	for _, right := range posList {
		b.WriteString(lower[left:right])
		b.WriteByte('_')
		left = right
	}
	b.WriteString(lower[left:])
	return b.String()
}
func ShortStr(s string, max int) string {
	sr := []rune(s)
	if len(sr) > max {
		sr = sr[:max]
	}
	return string(sr)
}
func ShortStrConvertLineEnding(s string, max int) string {
	x := ShortStr(s, max)
	if len(s) > max {
		x += "..."
	} else {
		x = s
	}
	x = strings.Replace(x, "\n", "\\n", -1)
	x = strings.Replace(x, "\r", "\\r", -1)
	return x
}

// StrEncode 字符串混淆
func StrEncode(s string) string {
	var isDouble bool
	b := []byte(s)
	bLen := len(b)
	rang := bLen / 2
	if bLen%2 == 0 {
		isDouble = true
	}
	var buffer bytes.Buffer
	for i := 0; i < rang; i++ {
		j := bLen - 1 - i
		buffer.Write(b[i : i+1])
		if isDouble || j != rang {
			buffer.Write(b[j : j+1])
		}
	}
	return buffer.String()
}

// StrDecode 字符串反混淆
func StrDecode(s string) string {
	b := []byte(s)
	bLen := len(b)
	var buffer bytes.Buffer
	for i := 0; i < bLen; i = i + 2 {
		buffer.Write(b[i : i+1])
	}
	for i := bLen - 1; i >= 0; i = i - 2 {
		buffer.Write(b[i : i+1])
	}
	return buffer.String()
}
func StrEmpty(s string) bool {
	return s == ""
}
func StrEqWith(s string) func(o string) bool {
	return func(o string) bool {
		return s == o
	}
}

const insertCost = 1
const deleteCost = 1
const editCost = 2

func MinEditDistance(target string, source string) (distance int) {
	targetR := []rune(target)
	sourceR := []rune(source)
	n := len(targetR)
	m := len(sourceR)
	// create distance matrix
	matrix := make([][]int, n+1)
	for i := 0; i <= n; i++ {
		matrix[i] = make([]int, m+1)
	}
	for i := 0; i <= m; i++ {
		matrix[0][i] = i
	}
	for j := 0; j <= n; j++ {
		matrix[j][0] = j
	}
	for i := 1; i <= n; i++ {
		for j := 1; j <= m; j++ {
			var insertDistance int
			var substituteDistance int
			insertDistance = matrix[i-1][j] + insertCost
			if sourceR[j-1] == targetR[i-1] {
				substituteDistance = matrix[i-1][j-1]
			} else {
				substituteDistance = matrix[i-1][j-1] + editCost
			}
			deleteDistance := matrix[i][j-1] + deleteCost
			matrix[i][j] = Min(insertDistance, substituteDistance, deleteDistance).(int)
		}
	}
	distance = matrix[n][m]
	return
}

// UpperFirst 首字母大写
func UpperFirst(str string) string {
	for i, v := range str {
		return string(unicode.ToUpper(v)) + str[i+1:]
	}
	return ""
}

// LowerFirst 首字母小写
func LowerFirst(str string) string {
	for i, v := range str {
		return string(unicode.ToLower(v)) + str[i+1:]
	}
	return ""
}
