package main

import (
	"crypto/md5"
	"fmt"
	"github.com/stretchr/gomniauth"
	"github.com/stretchr/objx"
	"io"
	"log"
	"net/http"
	"strings"
)

type authHandler struct {
	next http.Handler
}

func (h *authHandler) ServeHTTP(w http.ResponseWriter, r *http.Request) {

	// 未認証
	if cookie, err := r.Cookie("auth"); err == http.ErrNoCookie || cookie.Value == "" {
		w.Header().Set("Location", "/login")
		w.WriteHeader(http.StatusTemporaryRedirect) // 307 一時的リダイレクト
	} else if err != nil {
		// 何らかのエラーが発生
		panic(err.Error())
	} else {
		// 成功。 ラップされたハンドラを呼び出す
		h.next.ServeHTTP(w, r)
	}
}

func MustAuth(handler http.Handler) http.Handler {
	return &authHandler{next: handler}
}

/*
// loginHandlerはサードパーティへのログインの処理を受け持つ
// パスの形式: /auth/{action}/{provider}
この関数は内部状態を保持する必要がない為、Handler interfaceを実装していない。
HnadleFunc関数でhttp.Handleと同様のパスの関連付けができる

type Handler interface {
	ServeHTTP(ResponseWriter, *Request)
}
*/
func loginHandler(w http.ResponseWriter, r *http.Request) {
	segs := strings.Split(r.URL.Path, "/")
	action := segs[2]
	provider := segs[3]
	switch action {

	case "login":
		provider, err := gomniauth.Provider(provider)
		if err != nil {
			log.Fatalln("プロバイダーの取得に失敗しました:", provider, "-", err)
		}
		loginUrl, err := provider.GetBeginAuthURL(nil, nil)
		if err != nil {
			log.Fatalln("エラー:GetBeginAuthURL", provider, "-", err)
		}
		w.Header().Set("Location", loginUrl)
		w.WriteHeader(http.StatusTemporaryRedirect)

	case "callback":
		provider, err := gomniauth.Provider(provider)
		if err != nil {
			log.Fatalln("プロバイダーの取得に失敗しました:", provider, "-", err)
		}

		creds, err := provider.CompleteAuth(objx.MustFromURLQuery(r.URL.RawQuery))
		if err != nil {
			log.Fatalln("認証を完了できませんでした", provider, "-", err)
		}

		user, err := provider.GetUser(creds)
		if err != nil {
			log.Fatalln("ユーザーの取得に失敗しました", provider, "-", err)
		}

		m := md5.New()
		io.WriteString(m, strings.ToLower(user.Name()))
		userID := fmt.Sprintf("%x", m.Sum(nil))

		// データを保存
		authCookieValue := objx.New(
			map[string]interface{}{
				"userid":     userID,
				"name":       user.Name(),
				"avatar_url": user.AvatarURL(),
				"email":      user.Email(),
			}).MustBase64()

		http.SetCookie(w, &http.Cookie{Name: "auth", Value: authCookieValue, Path: "/"})
		w.Header()["Location"] = []string{"/chat"}
		w.WriteHeader(http.StatusTemporaryRedirect)

	default:
		w.WriteHeader(http.StatusNotFound)
		fmt.Fprintf(w, "アクション%sには非対応です", action)
	}
}
