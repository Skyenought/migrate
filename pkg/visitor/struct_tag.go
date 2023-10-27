package visitor

import (
	"go/ast"
	"go/token"
	"golang.org/x/tools/go/ast/astutil"
	"regexp"
	"strings"
)

func (v *Visitor) ReplaceGoTag(c *astutil.Cursor, old string, new string) {
	if field, ok := c.Node().(*ast.Field); ok && field.Tag != nil {

		if field.Tag.Kind == token.STRING && strings.Contains(field.Tag.Value, old+":") {
			regex := regexp.MustCompile(`uri:`)
			replaced := regex.ReplaceAllString(field.Tag.Value, "path:")
			field.Tag.Value = replaced
		}
	}
}
