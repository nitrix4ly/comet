package gen

import (
	"strings"
	"text/template"
	"unicode"

	"github.com/nitrix4ly/comet/core"
)

func init() {
	template.Must(template.New("").Funcs(templateFuncs).Parse(""))
}

var templateFuncs = template.FuncMap{
	"ToSnakeCase":  core.ToSnakeCase,
	"ToPascalCase": core.ToPascalCase,
	"ToCamelCase":  core.ToCamelCase,
	"ToPlural":     core.ToPlural,
	"ToLower":      strings.ToLower,
	"ToUpper":      strings.ToUpper,
	"Title":        strings.Title,
	"HasPrefix":    strings.HasPrefix,
	"HasSuffix":    strings.HasSuffix,
	"TrimSpace":    strings.TrimSpace,
	"Join":         strings.Join,
	"Split":        strings.Split,
	"Replace":      strings.Replace,
	"Contains":     strings.Contains,
	"IsEmpty": func(s string) bool {
		return strings.TrimSpace(s) == ""
	},
	"IsNotEmpty": func(s string) bool {
		return strings.TrimSpace(s) != ""
	},
	"FirstLower": func(s string) string {
		if len(s) == 0 {
			return s
		}
		r := []rune(s)
		r[0] = unicode.ToLower(r[0])
		return string(r)
	},
	"FirstUpper": func(s string) string {
		if len(s) == 0 {
			return s
		}
		r := []rune(s)
		r[0] = unicode.ToUpper(r[0])
		return string(r)
	},
	"Add": func(a, b int) int {
		return a + b
	},
	"Sub": func(a, b int) int {
		return a - b
	},
	"Mul": func(a, b int) int {
		return a * b
	},
	"Div": func(a, b int) int {
		if b == 0 {
			return 0
		}
		return a / b
	},
	"Mod": func(a, b int) int {
		if b == 0 {
			return 0
		}
		return a % b
	},
	"Eq": func(a, b interface{}) bool {
		return a == b
	},
	"Ne": func(a, b interface{}) bool {
		return a != b
	},
	"Lt": func(a, b int) bool {
		return a < b
	},
	"Le": func(a, b int) bool {
		return a <= b
	},
	"Gt": func(a, b int) bool {
		return a > b
	},
	"Ge": func(a, b int) bool {
		return a >= b
	},
	"And": func(a, b bool) bool {
		return a && b
	},
	"Or": func(a, b bool) bool {
		return a || b
	},
	"Not": func(a bool) bool {
		return !a
	},
}
