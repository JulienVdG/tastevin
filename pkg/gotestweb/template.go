// Copyright 2019 Splitted-Desktop Systems. All rights reserved
// Copyright 2019 Julien Viard de Galbert
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package gotestweb

import (
	"html/template"
	"io"
	"log"
)

// IndexData group the options for rendering the Index template
type IndexData struct {
	File      string
	Asciicast string
	Summary   bool
	AppPrefix string
	UseCDN    bool
}

// WriteIndex customize the index with given parameters and outputs it
func WriteIndex(w io.Writer, data IndexData) {
	fmap := template.FuncMap{
		"htmlSafe": htmlSafe,
	}
	tmpl, err := template.New("index").Funcs(fmap).Parse(indexTemplate)
	if err != nil {
		log.Fatal(err)
	}
	tmpl.Execute(w, data)
}

func htmlSafe(text string) template.HTML {
	return template.HTML(text)
}

var indexTemplate = `<!DOCTYPE html>
<html lang="en">

<head>

    <meta charset="utf-8">
    <meta http-equiv="X-UA-Compatible" content="IE=edge">
    <meta name="viewport" content="width=device-width, initial-scale=1">

    <title>GoTestWeb</title>
{{if .UseCDN}}
    <link rel="stylesheet" href="https://stackpath.bootstrapcdn.com/bootstrap/4.3.1/css/bootstrap.min.css" integrity="sha384-ggOyR0iXCbMQv3Xipma34MD+dH/1fQ784/j6cY/iJTQUOhcWr7x9JvoRxT2MZw1T" crossorigin="anonymous">
{{else}}
    <link href="{{.AppPrefix}}vendor/bootstrap/css/bootstrap.min.css" rel="stylesheet">
{{end}}
    <link href="{{.AppPrefix}}vendor/font-awesome/css/all.min.css" rel="stylesheet" type="text/css">
    <link href="{{.AppPrefix}}css/main.min.css" rel="stylesheet">
    <link rel="stylesheet" type="text/css" href="{{.AppPrefix}}vendor/asciinema/asciinema-player.css" />

    {{htmlSafe "<!-- HTML5 Shim and Respond.js IE8 support of HTML5 elements and media queries -->"}}
    {{htmlSafe "<!-- WARNING: Respond.js doesn't work if you view the page via file:// -->"}}
    {{htmlSafe "<!--[if lt IE 9]>"}}
        <script src="https://oss.maxcdn.com/libs/html5shiv/3.7.0/html5shiv.js"></script>
        <script src="https://oss.maxcdn.com/libs/respond.js/1.4.2/respond.min.js"></script>
    {{htmlSafe "<![endif]-->"}}

</head>

<body>

    <div id="wrapper">
    </div>

{{if .UseCDN}}
    <script src="https://code.jquery.com/jquery-3.2.1.min.js" integrity="sha256-hwg4gsxgFZhOsEEamdOYGBf13FyQuiTwlAQgxVSNgt4=" crossorigin="anonymous"></script>
    <script src="https://stackpath.bootstrapcdn.com/bootstrap/4.3.1/js/bootstrap.min.js" integrity="sha384-JjSmVgyd0p3pXB1rRibZUAYoIIy6OrQ6VrjIEaFf/nJGzIxFDsf4x0xIM+B07jRM" crossorigin="anonymous"></script>
{{else}}
    <script src="{{.AppPrefix}}vendor/jquery/jquery.min.js"></script>
    <script src="{{.AppPrefix}}vendor/bootstrap/js/bootstrap.min.js"></script>
{{end}}
    <script src="{{.AppPrefix}}vendor/asciinema/asciinema-player.js"></script>

    <script src="{{.AppPrefix}}js/app.min.js"{{if .File}} data-file="{{.File}}"{{end}}{{if .Summary}} data-summary="true"{{end}}{{if .Asciicast}} data-asciicast="{{.Asciicast}}"{{end}}></script>
</body>

</html>
`
