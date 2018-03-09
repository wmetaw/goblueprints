package main

import (
	"html/template"
	"log"
	"net/http"
	"path/filepath"
	"sync"
)

// templは１つのテンプレートを表す
type templateHandler struct {
	once     sync.Once
	filename string
	templ    *template.Template
}

// ServeHTTPはHTTPリクエストを処理
// ServeHTTPの中でテンプレートをコンパイルすると、本当に必要になるまで処理を後回しにできる。これを遅延初期化(lagy initilization)という
func (t *templateHandler) ServeHTTP(w http.ResponseWriter, r *http.Request) {

	// 一度だけ実行
	t.once.Do(func() {
		// テンプレートをコンパイル
		t.templ = template.Must(template.ParseFiles(filepath.Join("chapter1/chat/templates", t.filename)))
	})

	// コンパイルしたテンプレートをResponseWriterに出力
	t.templ.Execute(w, nil)
}

func main() {

	http.Handle("/", &templateHandler{filename: "chat.html"})

	// webサーバーの開始
	if err := http.ListenAndServe(":8080", nil); err != nil {
		log.Fatal("ListenAndServe", err)
	}
}
