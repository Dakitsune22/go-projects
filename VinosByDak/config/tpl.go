package config

import (
	"html/template"
	"strings"
)

var fm = template.FuncMap{
	"iOne":  incrementOne,
	"hello": hello,
	"tu":    strings.ToUpper,
}

var Tpl *template.Template

func init() {
	//Tpl = template.Must(template.ParseGlob("templates/*.gohtml"))
	Tpl = template.Must(template.New("").Funcs(fm).ParseGlob("templates/*.gohtml"))
}

func incrementOne(n int) int {
	return n + 1
}

func hello(s string) string {
	return "Hola... " + s
}
