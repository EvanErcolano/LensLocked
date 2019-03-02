package views

import (
	"bytes"
	"fmt"
	"html/template"
	"io"
	"net/http"
	"path/filepath"

	"lenslocked.com/context"
)

var (
	LayoutDir   string = "views/layouts/"
	TemplateDir string = "views/"
	TemplateExt string = ".gohtml"
)

func NewView(layout string, files ...string) *View {
	addTemplatePath(files)
	addTemplateExt(files)
	files = append(files, layoutFiles()...) // ... unpacks the slice into multiple strings

	t, err := template.ParseFiles(files...)
	if err != nil {
		panic(err)
	}
	return &View{
		Template: t,
		Layout:   layout,
	}
}

type View struct {
	Template *template.Template
	Layout   string
}

func (v *View) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	v.Render(w, r, nil)
}

func (v *View) Render(w http.ResponseWriter, r *http.Request, data interface{}) error {
	w.Header().Set("Content-Type", "text/html")
	var vd Data
	switch d := data.(type) {
	case Data:
		vd = d
	default:
		vd = Data{
			Yield: data,
		}
	}
	vd.User = context.User(r.Context())
	var buf bytes.Buffer
	if err := v.Template.ExecuteTemplate(&buf, v.Layout, vd); err != nil {
		fmt.Println(err)
		http.Error(w, "Something went wrong. If the problems persists, please email support@lenslocked.com",
			http.StatusInternalServerError)
		return nil
	}
	io.Copy(w, &buf)
	return nil
}

// layoutFiles returhns a slice of strfings representing
// the layouut files used in our app
func layoutFiles() []string {
	files, err := filepath.Glob(LayoutDir + "*" + TemplateExt)
	if err != nil {
		panic(err)
	}
	return files
}

// addTemplatePath takes in a slice of strings repr file paths for templates
// prepends tempalte dir directory to each string in slice
// e.g. the input "home" would resolve to "views/home" if templateDirs == "views/"
func addTemplatePath(files []string) {
	for i, f := range files {
		files[i] = TemplateDir + f
	}
}

// addTemplateExt takes in a slice of strings representing filepaths for templates
// it appends the template extension to each slice
// "home" would resolve to "home.gothml" if templatExt == .gohtml
func addTemplateExt(files []string) {
	for i, f := range files {
		files[i] = f + TemplateExt
	}
}
