package main

import (
	"flag"
	"html/template"
	"log"
	"net/http"
	"path/filepath"
	"sync"
	_ "os"
	_ "github.com/wakwa3125/chat/trace"
	"github.com/stretchr/gomniauth"
	"github.com/stretchr/gomniauth/providers/google"
	"github.com/stretchr/gomniauth/providers/facebook"
	"github.com/stretchr/gomniauth/providers/github"
)

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
	
	gomniauth.SetSecurityKey(SEC_KEY)
	gomniauth.WithProviders(
		google.New(GOOGLE_CLIENT_ID, GOOGLE_SECRET, GOOGLE_CALLBACK_URL),
		facebook.New(FB_CLIENT_ID, FB_SECRET, FB_CALLBACK_URL),
		github.New(GITHUB_CLIENT_ID, GITHUB_SECRET, GITHUB_CALLBACK_URL),
	)
	
	r := newRoom()
	// r.tracer = trace.New(os.Stdout)
	http.Handle("/chat", MustAuth(&templateHandler{filename: "chat.html"}))
	http.Handle("/login", &templateHandler{filename: "login.html"})
	http.HandleFunc("/auth/", loginHandler)
	http.Handle("/room", r)
	
	go r.run()
	
	log.Println("Webサーバーを開始します。ポート: ", *addr)
	
	if err := http.ListenAndServe(*addr, nil); err != nil {
		log.Fatal("ListenAndServe", err)
	}
}
