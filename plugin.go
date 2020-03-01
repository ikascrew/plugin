package main

import (
	"encoding/json"
	"io/ioutil"
	"os"
	"text/template"
)

type pluginDto struct {
	Package string            `json:"package"`
	Plugins map[string]string `json:"plugins"`
}

//Generate
func main() {

	name := "plugin.json"

	//指定されたJsonファイルを読み込み
	b, err := ioutil.ReadFile(name)
	if err != nil {
		panic(err)
	}

	dto := pluginDto{}
	err = json.Unmarshal(b, &dto)
	if err != nil {
		panic(err)
	}

	out := "video_plugin_gen.go"

	f, err := os.Create(out)
	if err != nil {
		panic(err)
	}
	defer f.Close()

	tmpl := template.New("")

	tmpl, err = tmpl.Parse(pluginTmpl)
	if err != nil {
		panic(err)
	}

	err = tmpl.Execute(f, dto)
	if err != nil {
		panic(err)
	}

}

const pluginTmpl = `
package {{ .Package }}

import (

	"fmt"
	"github.com/ikascrew/core"

{{ range $key,$val := .Plugins }}
	{{$key}} "{{$val}}"
{{ end }}
)

var NotFoundError = fmt.Errorf("NotFound Video Type")

func Get(t string, params ...string) (core.Video, error) {

	var v video.Video
	var err error

	switch t {
{{ range $key,$val := .Plugins }}
	case "{{ $key }}":
		v, err = {{$key}}.New(params)
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
