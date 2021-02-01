// This small program is used to generate request structs for the SonarCloud API.
// It expects a JSON file with the same structure as returned by `https://sonarcloud.io/api/webservices/list`.
// See AllowedEndpoints for the list of endpoints that are considered during generation.
package main

import (
	"encoding/json"
	"flag"
	"fmt"
	. "github.com/dave/jennifer/jen"
	"github.com/iancoleman/strcase"
	"io/ioutil"
	"os"
	"strings"
)

var (
	AllowedEndpoints = []string{"user_groups", "permissions", "user_tokens"}
)

type Api struct {
	Services []Service `json:"webServices"`
}

func main() {
	var filename string
	var output string
	flag.StringVar(&filename, "filename", "gen/services.json", "name of the file which contains the api definition")
	flag.StringVar(&output, "output", "pkg/api/", "directory where the generated files will be stored")

	file, err := ioutil.ReadFile(filename)
	guard(err)

	var api Api
	err = json.Unmarshal(file, &api)
	guard(err)

	for _, service := range api.Services {
		services(service, output)
	}
}

func exit(code int, s interface{}) {
	fmt.Println(s)
	os.Exit(code)
}

func guard(err error) {
	if err != nil {
		exit(1, err)
	}
}

func contains(needle string, haystack []string) bool {
	found := false
	for _, hay := range haystack {
		if hay == needle {
			found = true
			break
		}
	}
	return found
}

type Service struct {
	Path        string   `json:"path"`
	Description string   `json:"description"`
	Actions     []Action `json:"actions"`
}

type Action struct {
	Key                string  `json:"key"`
	Description        string  `json:"description"`
	Internal           bool    `json:"internal"`
	Post               bool    `json:"post"`
	HasResponseExample bool    `json:"hasResponseExample"`
	Params             []Param `json:"params"`
	DeprecatedSince    string  `json:"deprecatedSince"`
}

type Param struct {
	Key             string `json:"key"`
	Description     string `json:"description"`
	Internal        bool   `json:"internal"`
	Required        bool   `json:"required"`
	DeprecatedSince string `json:"deprecatedSince"`
}

func services(service Service, output string) {
	path := strings.Split(service.Path, "/")
	endpoint := path[len(path)-1]

	if !contains(endpoint, AllowedEndpoints) {
		return
	}

	fmt.Println("Generating request types for: " + service.Path)

	f := NewFile("api")
	f.Commentf("// AUTOMATICALLY GENERATED, DO NOT EDIT BY HAND!\n")

	for _, action := range service.Actions {
		if action.HasResponseExample {
			// TODO: generate response type
		}

		statements := make([]Code, 0)
		for _, param := range action.Params {
			id := strcase.ToCamel(param.Key)
			statement := Id(id).String().Tag(map[string]string{"form": param.Key + ",omitempty"}).Comment(param.Description)
			statements = append(statements, statement)
		}

		id := strcase.ToCamel(fmt.Sprintf("%s_%s", endpoint, action.Key))
		f.Commentf("%s: %s", id, action.Description)

		if action.DeprecatedSince != "" {
			f.Commentf("Deprecated: this action has been deprecated since version %s", action.DeprecatedSince)
		}

		f.Type().Id(id).Struct(statements...)
	}

	err := f.Save(fmt.Sprintf("%s%s.go", output, endpoint))
	if err != nil {
		fmt.Printf("ERROR: %+v\n", err)
	}
}
