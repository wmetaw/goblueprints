package main

import (
	"flag"
	"github.com/stretchr/gomniauth"
	"github.com/stretchr/gomniauth/providers/google"
	"github.com/stretchr/objx"
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

	data := map[string]interface{}{
		"Host": r.Host,
	}

	if authCookie, err := r.Cookie("auth"); err == nil {
		data["UserData"] = objx.MustFromBase64(authCookie.Value)
	}

	// コンパイルしたテンプレートをResponseWriterに出力
	t.templ.Execute(w, data)
}

func main() {

	// コマンドライン引数で受け取った値をパース
	var addr = flag.String("addr", ":8080", "アプリケーションのアドレス")
	flag.Parse()

	// Gomuniauthのセットアップ
	gomniauth.SetSecurityKey("hoge")
	gomniauth.WithProviders(
		google.New(
			"",
			"",
			"http://localhost:8080/auth/callback/google"),
	)

	// AuthAvatarのインスタンスを生成していないため、メモリ使用量が増えることはない
	// 大量のチャットルームを生成する状況では大幅なメモリの節約が期待できる
	r := newRoom(UseFileSystemAvatar)

	r.tracer = trace.New(os.Stdout)
	http.Handle("/", MustAuth(&templateHandler{filename: "chat.html"}))
	http.Handle("/login", &templateHandler{filename: "login.html"})
	http.Handle("/upload", &templateHandler{filename: "upload.html"})
	http.Handle("/room", r)

	http.Handle("/avatars/",
		http.StripPrefix("/avatars/", http.FileServer(http.Dir("./avatars"))))

	http.HandleFunc("/uploader", uploaderHandler)
	http.HandleFunc("/auth/", loginHandler)
	http.HandleFunc("/logout", func(w http.ResponseWriter, r *http.Request) {
		http.SetCookie(w, &http.Cookie{
			Name:   "auth",
			Value:  "",
			Path:   "/",
			MaxAge: -1,
		})
		w.Header()["Location"] = []string{"/chat"}
		w.WriteHeader(http.StatusTemporaryRedirect)
	})

	// chatルームを開始
	go r.run()

	// WEbサーバーの起動
	log.Println("Webサーバーを開始します。ポート ", *addr)

	// webサーバーの開始
	if err := http.ListenAndServe(*addr, nil); err != nil {
		log.Fatal("ListenAndServe", err)
	}
}
