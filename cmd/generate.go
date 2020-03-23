package main

import (
	"encoding/json"
	"flag"
	"fmt"
	"io/ioutil"
	"os"
	"text/template"
)

type pluginDto struct {
	Videos map[string]string `json:"videos"`
}

//Generate
func main() {

	flag.Parse()
	args := flag.Args()

	if len(args) < 1 {
		fmt.Println("argument error(plugin.json)")
		os.Exit(1)
	}

	out := "plugins_gen.go"

	err := generate(args[0], out)
	if err != nil {
		fmt.Println(err)
		os.Exit(1)
	}

	fmt.Println("Success.")

}

func generate(name string, out string) error {

	//指定されたJsonファイルを読み込み
	b, err := ioutil.ReadFile(name)
	if err != nil {
		return err
	}

	dto := pluginDto{}
	err = json.Unmarshal(b, &dto)
	if err != nil {
		return err
	}

	f, err := os.Create(out)
	if err != nil {
		return err
	}
	defer f.Close()

	tmpl := template.New("")

	tmpl, err = tmpl.Parse(pluginTmpl)
	if err != nil {
		return err
	}

	err = tmpl.Execute(f, dto)
	if err != nil {
		return err
	}

	return nil
}

const pluginTmpl = `
package main

import (

	"fmt"
	"github.com/ikascrew/core"

{{ range $key,$val := .Videos }}
	{{$key}} "{{$val}}"
{{ end }}
)

var NotFoundError = fmt.Errorf("NotFound Video Type")

func GetVideo(t string, params ...string) (core.Video, error) {

	var v core.Video
	var err error

	switch t {
{{ range $key,$val := .Videos }}
	case "{{ $key }}":
		v, err = {{$key}}.New(params...)
{{ end }}
	}

	if v == nil {
		err = NotFoundError
	}

	if err != nil {
		return nil, err
	}
	return v, nil
}
`
