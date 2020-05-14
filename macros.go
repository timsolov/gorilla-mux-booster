package gmb

import "strings"

var macros = map[string]string{
	"@uuid@": "[0-9a-f]{8}-[0-9a-f]{4}-[0-9a-f]{4}-[0-9a-f]{4}-[0-9a-f]{12}",
	"@num@":  "[0-9]+",
}

// Macros adds new macros with name and regex.
// Then you can use it in your routes as @name@
//   router.GET("/{param:@name@}/ok", handler)
func Macros(name, regex string) {
	macros["@"+name+"@"] = regex
}

// m replaces macroses in uri string
func m(uri string) string {
	const (
		colon = ":"
		slash = "/"
		left  = "{"
		right = "}"
	)

	uriParts := strings.Split(uri, slash)

	for i := 0; i < len(uriParts); i++ {
		part := uriParts[i]

		if !strings.Contains(part, colon) || !strings.HasPrefix(part, left) || !strings.HasSuffix(part, right) {
			continue
		}

		s := strings.TrimSuffix(strings.TrimPrefix(part, left), right)

		paramParts := strings.SplitN(s, colon, 2)
		tpl := paramParts[1]

		if v, yes := macros[tpl]; yes {
			uriParts[i] = left + strings.Replace(s, tpl, v, 1) + right
		}
	}

	return strings.Join(uriParts, slash)
}
