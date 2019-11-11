package main

import (
	"html/template"
	"io/ioutil"
	"log"
	"net/http"
)

// Page Structure
type Page struct {
	Title   string
	Content []byte // a byte slice
}

// For persistent storage - add a save() method to Page
func (p *Page) save() error {
	filename := p.Title + ".txt"
	return ioutil.WriteFile(filename, p.Content, 0600) // 0600 is the permission -> read-write permission for current user only
}

func loadPage(title string) (*Page, error) {
	filename := title + ".txt"
	content, err := ioutil.ReadFile(filename)
	if err != nil {
		return nil, err
	}
	return &Page{Title: title, Content: content}, nil
}

func renderTemplate(templateName string, w http.ResponseWriter, p *Page) {
	t, _ := template.ParseFiles(templateName + ".html")
	t.Execute(w, p)
}

func viewHandler(w http.ResponseWriter, r *http.Request) {
	title := r.URL.Path[len("/view/"):]
	p, err := loadPage(title)
	if err != nil {
		http.Redirect(w, r, "/edit/"+title, http.StatusFound)
		return
	}
	renderTemplate("view", w, p)
}

func editHandler(w http.ResponseWriter, r *http.Request) {
	title := r.URL.Path[len("/edit/"):]
	p, err := loadPage(title)
	if err != nil {
		p = &Page{Title: title}
	}
	renderTemplate("edit", w, p)
}

func saveHandler(w http.ResponseWriter, r *http.Request) {
	title := r.URL.Path[len("/save/"):]
	content := r.FormValue("body")
	p := &Page{Title: title, Content: []byte(content)}
	p.save()
	http.Redirect(w, r, "/view/"+title, http.StatusFound)
}

func main() {
	http.HandleFunc("/view/", viewHandler)
	http.HandleFunc("/edit/", editHandler)
	http.HandleFunc("/save/", saveHandler)
	log.Fatal(http.ListenAndServe(":8080", nil))
}
