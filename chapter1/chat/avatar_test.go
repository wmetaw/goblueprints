package main

import "testing"

func TestAuthAvatar(t *testing.T) {
	var authAvatar AuthAvatar
	client := new(client)
	url, err := authAvatar.GetAvatarURL(client)
	if err != ErrNoAvatarURL {
		t.Error("値が存在しない場合、 authAvatar.GetAvatarURLはErrNoAvatarURLを返すべきです")
	}

	// 値をセットします
	testUrl := "http://url-to-avatar/"
	client.userData = map[string]interface{}{"avatar_url": testUrl}
	url, err = authAvatar.GetAvatarURL(client)
	if err != nil {
		t.Error("値が存在する場合、 authAvatar.GetAvatarURLはエラーを返すべきではありません")
	} else {
		if url != testUrl {
			t.Error("authAvatar.GetAvatarURLは正しいURLを返すべきです")
		}
	}
}

func TestGravatarAvatar(t *testing.T) {
	var gravatarAvatar GravatarAvatar
	client := new(client)
	//client.userData = map[string]interface{}{"email": "lio114514@gmail.com"}
	client.userData = map[string]interface{}{"userid": "0bc83cb571cd1c50ba6f3e8af1346"}
	url, err := gravatarAvatar.GetAvatarURL(client)
	if err != nil {
		t.Error("GravatarAvatar.GetAvatarURLはエラーを返すべきではありません")
	}
	if url != "//www.gravatar.com/avatar/abc" {
		t.Errorf("GravatarAvatar.GetAvatarURLが%sという誤った値を返しました", url)
	}
}
