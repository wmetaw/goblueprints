package main

import "errors"

func main() {
}

// ErrNoAvatarURLはAvatarインスタンスがアバターのURLを返すことができない場合に発生するエラー
var ErrNoAvatarURL = errors.New("chat: アバターURLを取得できません。")

// GetAvatarURLは指定されたクライアントのアバターのURLを返す。
// 問題が発生した場合にはエラーを返す。特にURLを取得できなかった場合にはErrNoAvatarURLを返す
type Avatar interface {
	GetAvatarURL(c *client) (string, error)
}
