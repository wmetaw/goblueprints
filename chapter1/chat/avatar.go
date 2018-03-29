package main

import "errors"

// ErrNoAvatarURLはAvatarインスタンスがアバターのURLを返すことができない場合に発生するエラー
var ErrNoAvatarURL = errors.New("chat: アバターURLを取得できません。")

// GetAvatarURLは指定されたクライアントのアバターのURLを返す。
// 問題が発生した場合にはエラーを返す。特にURLを取得できなかった場合にはErrNoAvatarURLを返す
type Avatar interface {
	GetAvatarURL(c *client) (string, error)
}

type AuthAvatar struct{}

var UseAvatar AuthAvatar

func (_ AuthAvatar) GetAvatarURL(c *client) (string, error) {
	if url, ok := c.userData["avatar_url"]; ok {
		if urlStr, ok := url.(string); ok {
			return urlStr, nil
		}
	}
	return "", ErrNoAvatarURL
}
