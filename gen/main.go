package main

import (
	"context"
	"encoding/json"
	"flag"
	"fmt"
	"github.com/iancoleman/strcase"
	"io/ioutil"
	"net/http"
	"os"
	"strings"
	"time"

	. "github.com/dave/jennifer/jen"
)

var (
	AllowedEndpoints = []string{"user_groups", "permissions"}
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

	//client(output)
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

type ResponseExample struct {
	Format  string `json:"format"`
	Example string `json:"example"`
}

func (re ResponseExample) Keys() []string {
	//fmt.Printf("Example: %+v\n", re.Example)

	example := make(map[string]interface{})
	err := json.Unmarshal([]byte(re.Example), &example)
	guard(err)

	i := 0
	keys := make([]string, len(example))
	for key := range example {
		keys[i] = key
		i++
	}

	return keys
}

// TODO: correctly implement creating a response type from the example
func response(controller, action string) {
	ctx := context.Background()
	ctx, _ = context.WithTimeout(ctx, 10*time.Second)
	req, err := http.NewRequestWithContext(ctx, "GET", "http://www.sonarcloud.io/api/webservices/response_example", nil)
	guard(err)

	q := req.URL.Query()
	q.Add("action", action)
	q.Add("controller", controller)
	req.URL.RawQuery = q.Encode()

	client := &http.Client{}
	resp, err := client.Do(req)
	guard(err)

	defer resp.Body.Close()

	body, err := ioutil.ReadAll(resp.Body)
	guard(err)

	re := ResponseExample{}
	err = json.Unmarshal(body, &re)
	guard(err)

	fmt.Printf("%+v\n", re.Keys())
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
		//fmt.Println("Action: " + action.Key)
		if action.HasResponseExample {
			//fmt.Println("Response Example:")
			// TODO: generate response type
			//response(service.Path, action.Key)
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

// WIP
func client(output string) {
	fmt.Println("Generating client")

	f := NewFile("api")
	f.Commentf("// AUTOMATICALLY GENERATED, DO NOT EDIT BY HAND!\n")

	f.Func().Id("New").Params().Block(
		Id("client").Op(":=").Qual("http", "Client"),
		Return(Id("client")),
	)

	//client := &http.Client{Timeout: 10 * time.Second}
	//
	//var diags diag.Diagnostics
	//
	//req, err := http.NewRequest("GET", fmt.Sprintf("%s/user_groups/search", API), nil)
	//if err != nil {
	//	return diag.FromErr(err)
	//}
	//
	//c := m.(*Config)
	//q := req.URL.Query()
	//q.Add("organization", c.Organization)
	//req.URL.RawQuery = q.Encode()
	//
	//req.SetBasicAuth(c.Token, "")

	err := f.Save(fmt.Sprintf("%s/client.go", output))
	if err != nil {
		fmt.Printf("ERROR: %+v\n", err)
	}
}
