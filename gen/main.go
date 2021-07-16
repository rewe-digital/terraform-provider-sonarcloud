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
	"net/http"
	"os"
	"sort"
	"strings"
	"time"
)

var (
	AllowedEndpoints = []string{"user_groups", "permissions", "user_tokens"}
)

type Api struct {
	Services []Service `json:"webServices"`
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

type ResponseExampleRequest struct {
	ID         string
	RequestID  string
	Action     string
	Controller string
}

type ResponseExampleRequestResponse struct {
	Format  string `json:"format"`
	Example string `json:"example"` // yes, it's a string...
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

func services(service Service, output string) {
	path := strings.Split(service.Path, "/")
	endpoint := path[len(path)-1]

	if !contains(endpoint, AllowedEndpoints) {
		return
	}

	fmt.Println("Generating request and response types for: " + service.Path)

	f := NewFile("api")
	f.Commentf("// AUTOMATICALLY GENERATED, DO NOT EDIT BY HAND!\n")

	for _, action := range service.Actions {
		statements := make([]Code, 0)
		for _, param := range action.Params {
			id := strcase.ToCamel(param.Key)
			statement := Id(id).String().Tag(map[string]string{"form": param.Key + ",omitempty"}).Comment(param.Description)
			statements = append(statements, statement)
		}

		id := strcase.ToCamel(fmt.Sprintf("%s_%s", endpoint, action.Key))
		f.Commentf("%s %s", id, action.Description)

		if action.DeprecatedSince != "" {
			f.Commentf("Deprecated: this action has been deprecated since version %s", action.DeprecatedSince)
		}

		f.Type().Id(id).Struct(statements...)

		if action.HasResponseExample {
			controller := fmt.Sprintf("api/%s", endpoint)
			responseId := strcase.ToCamel(fmt.Sprintf("%s_%s", id, "Response"))
			request := ResponseExampleRequest{ID: responseId, RequestID: id, Controller: controller, Action: action.Key}

			err := response(f, request)
			if err != nil {
				fmt.Printf("ERROR: %+v\n", err)
			}
		}

	}

	err := f.Save(fmt.Sprintf("%s%s.go", output, endpoint))
	if err != nil {
		fmt.Printf("ERROR: %+v\n", err)
	}
}

func response(f *File, request ResponseExampleRequest) error {
	client := &http.Client{Timeout: 10 * time.Second}

	req, err := newRequest(request)
	if err != nil {
		return fmt.Errorf("could not create request: %+v", err)
	}

	res, err := client.Do(req)
	if err != nil {
		return fmt.Errorf("error sending reqeust: %+v", err)
	}
	defer res.Body.Close()

	body, err := ioutil.ReadAll(res.Body)
	if err != nil {
		return fmt.Errorf("could not read response body: %+v", err)
	}

	var j ResponseExampleRequestResponse
	err = json.Unmarshal(body, &j)
	if err != nil {
		return fmt.Errorf("could not unmarshall body: %+v", err)
	}

	// Convert the example JSON string (!!) to a map
	var example map[string]interface{}
	err = json.Unmarshal([]byte(j.Example), &example)
	if err != nil {
		return fmt.Errorf("could not marshall example: %+v", err)
	}

	code := responseStructs(example)

	f.Commentf("%s is the response for %s", request.ID, request.RequestID)
	f.Type().Id(request.ID).Struct(code...)

	return nil
}

func responseStructs(example map[string]interface{}) []Code {
	code := make([]Code, 0)

	keys := sortedKeys(example)

	for _, k := range keys {
		v := example[k]
		statement := responseField(k, v)

		statement = statement.Tag(map[string]string{"json": k + ",omitempty"})
		code = append(code, statement)
	}
	return code
}

func responseField(k string, v interface{}) *Statement {
	id := strcase.ToCamel(k)
	// Note: an empty id is valid and simply doesn't output anything when rendered
	statement := Id(id)

	switch v.(type) {
	case string:
		statement = statement.String()
	case float64:
		statement = statement.Float64()
	case bool:
		statement = statement.Bool()
	case map[string]interface{}:
		subCode := responseStructs(v.(map[string]interface{}))
		statement = statement.Struct(subCode...)
	case []interface{}:
		// Only look at the first element of the array to determine its type and render it without identifier
		vv := v.([]interface{})[0]
		subCode := responseField("", vv)
		statement = statement.Index().Add(subCode)
	}
	return statement
}

func newRequest(responseExampleRequest ResponseExampleRequest) (*http.Request, error) {
	req, err := http.NewRequest("GET", "https://www.sonarcloud.io/api/webservices/response_example", nil)
	if err != nil {
		return nil, err
	}

	q := req.URL.Query()
	q.Add("action", responseExampleRequest.Action)
	q.Add("controller", responseExampleRequest.Controller)

	req.URL.RawQuery = q.Encode()

	req.Header.Set("Content-Type", "application/x-www-form-urlencoded")
	req.Header.Set("Accept", "application/json")

	return req, nil
}

func sortedKeys(m map[string]interface{}) []string {
	keys := make([]string, len(m))
	i := 0
	for k := range m {
		keys[i] = k
		i++
	}
	sort.Strings(keys)
	return keys
}
