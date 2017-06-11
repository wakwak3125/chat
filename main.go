package main

import (
	"flag"
	"html/template"
	"log"
	"net/http"
	"path/filepath"
	"sync"
	"os"
	"github.com/wakwa3125/trace"
)

const ROOT string = "/"

type templateHandler struct {
	once     sync.Once
	filename string
	template *template.Template
}

func (t *templateHandler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	t.once.Do(func() {
		t.template = template.Must(template.ParseFiles(filepath.Join("templates", t.filename)))
	})
	t.template.Execute(w, r)
}

func main() {
	addr := flag.String("addr", ":8080", "アプリケーションのアドレス")
	flag.Parse()
	
	r := newRoom()
	r.tracer = trace.New(os.Stdout)
	http.Handle(ROOT, &templateHandler{filename: "chat.html"})
	http.Handle("/room", r)
	
	go r.run()
	
	log.Println("Webサーバーを開始します。ポート: ", *addr)
	
	if err := http.ListenAndServe(*addr, nil); err != nil {
		log.Fatal("ListenAndServe", err)
	}
}
