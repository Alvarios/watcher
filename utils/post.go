package utils

import (
	"github.com/Alvarios/kushuh-go-utils/ext-apis/slack"
	"github.com/Alvarios/watcher/setup"
	"net/http"
)

func (sc *setup.SlackConfig) Error(m string) (*http.Response, error) {
	stack := PrintStack(2)

	return slack.Send(
		sc.WebHook,
	)
}