package main

import (
	"flag"
	"github.com/wmetaw/goblueprints/chapter1/trace"
	"html/template"
	"log"
	"net/http"
	"os"
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
		t.templ = template.Must(template.ParseFiles(filepath.Join("templates", t.filename)))
	})

	// コンパイルしたテンプレートをResponseWriterに出力
	t.templ.Execute(w, r)
}

func main() {

	// コマンドライン引数で受け取った値をパース
	var addr = flag.String("addr", ":8080", "アプリケーションのアドレス")
	flag.Parse()

	r := newRoom()
	r.tracer = trace.New(os.Stdout)
	http.Handle("/", MustAuth(&templateHandler{filename: "chat.html"}))
	http.Handle("/room", r)

	// chatルームを開始
	go r.run()

	// WEbサーバーの起動
	log.Println("Webサーバーを開始します。ポート ", *addr)

	// webサーバーの開始
	if err := http.ListenAndServe(*addr, nil); err != nil {
		log.Fatal("ListenAndServe", err)
	}
}
