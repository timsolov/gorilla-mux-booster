package gmb

import "strings"

const (
	colon  = ":"
	slash  = "/"
	left   = "{"
	right  = "}"
	breaks = left + right
)

var contractions = map[string]string{
	"@uuid@":     "[0-9a-f]{8}-[0-9a-f]{4}-[0-9a-f]{4}-[0-9a-f]{4}-[0-9a-f]{12}",
	"@num@":      "[0-9]+",
	"@alpha@":    "[a-zA-Z]+",
	"@alphanum@": "[a-zA-Z0-9]+",
}

// RegisterRegex adds new contraction with name and regex.
// Then you can use it in your routes as @name@
//   gmb.RegisterRegex("name", "[A-Z][a-z]+")
//   router.GET("/{param:@name@}/ok", handler)
func RegisterRegex(name, regex string) {
	contractions["@"+name+"@"] = regex
}

var aliases = map[string]string{}

// RegisterAlias adds new alias for REST variable.
// It's usefull when you want to replace all similar variables in your routes.
func RegisterAlias(what, to string) {
	what = strings.Trim(what, breaks)
	to = strings.Trim(to, breaks)
	aliases[what] = to
}

// RegisterAliases adds multiple aliases for REST variables.
// For more information please look at RegisterAlias.
func RegisterAliases(whatTo map[string]string) {
	for what, to := range whatTo {
		RegisterAlias(what, to)
	}
}

// c replaces contractions in uri string
func c(uri string) string {
	uriParts := strings.Split(uri, slash)

	for i := 0; i < len(uriParts); i++ {
		part := uriParts[i]

		if !strings.HasPrefix(part, left) || !strings.HasSuffix(part, right) {
			continue
		}

		s := strings.Trim(part, breaks)

		for what, to := range aliases {
			if s == what {
				s = to
			}
		}

		if !strings.Contains(s, colon) {
			continue
		}

		paramParts := strings.SplitN(s, colon, 2)
		tpl := paramParts[1]

		if v, yes := contractions[tpl]; yes {
			uriParts[i] = left + strings.Replace(s, tpl, v, 1) + right
		}
	}

	return strings.Join(uriParts, slash)
}
