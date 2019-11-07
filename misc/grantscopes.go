package misc

// GrantScope defines the structure of a grant scope
type GrantScope struct {
	Description string
	Details     string
	Scope       string
}

// GrantScopes defines a map of grant scopes
type GrantScopes map[string]GrantScope

// GetScopeListFromGrantScopeMap builds a list of scopes from a grant scope map
func (gss GrantScopes) GetScopeListFromGrantScopeMap() []string {
	var toReturn []string
	for _, v := range gss {
		toReturn = append(toReturn, v.Scope)
	}
	return toReturn
}
