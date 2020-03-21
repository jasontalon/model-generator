package main

import (
	"errors"
	"fmt"
	"io/ioutil"
	"log"
	"os"
	"path/filepath"
	"reflect"
	"strings"
	"text/template"
	"time"

	"gopkg.in/yaml.v2"
)

var dateTimeUtc string

func main() {
	generate()
}

func killIf(err error) {
	if err != nil {
		log.Fatal(err)
		os.Exit(1)
	}
}

func generate() {
	b, err := getYaml()

	killIf(err)

	m := make(map[interface{}]interface{})

	err = yamlToMap(&b, &m)

	killIf(err)

	if len(*&m) == 0 {
		killIf(errors.New("No models found"))
	}
	pkgName := ""
	if _, err = os.Stat("main.go"); err == nil {
		pkgName = "main"
	} else {
		pkgName, err = currentFolder()
	}

	if err != nil {
		return
	}

	file, err := os.Create("models.go")

	killIf(err)

	mf := ModelFl{PackageName: pkgName, DateGenerated: time.Now().UTC().Format(time.RFC3339), Models: []Model{}}

	for key, val := range *&m {
		if reflect.ValueOf(val).Kind() == reflect.Map {
			k := val.(map[interface{}]interface{})
			makeModel(fmt.Sprintf("%s", key), &k, &mf)
		}
	}

	pack, err := template.New("").Parse(modelTextTemplate)

	err = pack.Execute(file, &mf)

}

func (*ModelFl) JsonFormat(s string) string {
	return fmt.Sprintf("`json:\"%s\"`", strings.ToLower(s))
}

func makeModel(name string, props *map[interface{}]interface{}, mf *ModelFl) (err error) {
	s := Model{ModelName: name, Props: make(map[string]string)}

	for k, v := range *props {
		s.Props[k.(string)] = v.(string)
	}

	*&mf.Models = append(*&mf.Models, s)

	return
}

type ModelFl struct {
	DateGenerated string
	PackageName   string
	Models        []Model
}

type Model struct {
	ModelName string
	Props     map[string]string
}

var modelTextTemplate string = `//DO NOT MODIFY!
//This is an auto generated code!.
//DateTime Generated {{ .DateGenerated }} UTC
package {{ .PackageName }}
{{ range .Models }}
	type {{ .ModelName }} struct {
		{{ range $key, $value := .Props }} {{$key}} {{$value}} {{ $.JsonFormat $key }}
		{{ end }}
	}
{{ end }}`

func yamlToMap(b *[]byte, m *map[interface{}]interface{}) (err error) {
	err = yaml.Unmarshal(*b, &m)

	if err != nil {
		log.Fatal(err)
	}

	return
}

func getYaml() ([]byte, error) {
	c, err := ioutil.ReadFile("models.yaml")
	if err != nil {
		return nil, err
	}

	return c, nil
}

func currentFolder() (c string, err error) {
	c, err = os.Getwd()
	if err != nil {
		return
	}
	c = filepath.Base(c)
	return
}
