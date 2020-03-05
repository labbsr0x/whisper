package misc

import (
	"encoding/json"
	"io/ioutil"
	"os"
)

var AppScopes = []string{"openid", "offline"}  // openid and offline will be there by default

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
	for _, v := range gss {
		AppScopes = append(AppScopes, v.Scope)
	}
	return AppScopes
}

// GetGrantScopesFromFile reads into memory the json scopes file
func GetGrantScopesFromFile(scopesFilePath string) GrantScopes {
	jsonFile, err := os.Open(scopesFilePath)
	if err != nil {
		panic(err)
	}

	defer jsonFile.Close()

	var grantScopes GrantScopes

	bytes, err := ioutil.ReadAll(jsonFile)
	if err != nil {
		panic(err.Error())
	}

	err = json.Unmarshal(bytes, &grantScopes)
	if err != nil {
		panic(err.Error())
	}

	return grantScopes
}

