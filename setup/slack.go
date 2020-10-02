package setup

import (
	"fmt"
	"github.com/Alvarios/go-slack"
	soscrud "github.com/Alvarios/s-os/crud"
	"github.com/gin-gonic/gin"
	"log"
	"net/http"
	"os"
	"runtime/debug"
	"time"
)

type SlackConfig struct {
	WebHook     string `json:"webhook"`
	Application string `json:"application"`
}

type extract struct {
	Watcher SlackConfig `json:"watcher"`
}

func (sc *SlackConfig) LoadFrom(filePath string) error {
	var data extract
	err := soscrud.Get(filePath, &data)

	sc.WebHook = data.Watcher.WebHook
	sc.Application = data.Watcher.Application

	return err
}

func (sc *SlackConfig) Print(m string) string {
	return fmt.Sprintf(
		"*Message*\n%s\n\n*Stack*\n```%s```\n\n*Time*\n%s",
		m,
		string(debug.Stack()),
		time.Now().Format("2006-01-02 3:4:5"),
	)
}

func (sc *SlackConfig) Error(m string) (*http.Response, error) {
	env := os.Getenv("ENV")
	// If no ENV is specified, assume we are in development mode, so we don't want to flood Slack uselessly.
	if env == "" || env == "development" {
		return nil, nil
	}

	return slack.Send(
		sc.WebHook,
		fmt.Sprintf("Unexpected error in %s (%s)", sc.Application, env),
		nil,
		[]map[string]interface{}{
			{
				"fallback": fmt.Sprintf("Unexpected error in %s (%s)", sc.Application, env),
				"color":    "#FF9300",
				"text":     sc.Print(m),
			},
		},
	)
}

func (sc *SlackConfig) Fatal(m string) {
	fm := sc.Print(m)

	env := os.Getenv("ENV")
	// If no ENV is specified, assume we are in development mode, so we don't want to flood Slack uselessly.
	if env == "" || env == "development" {
		log.Fatalf(fm)
	}

	_, _ = slack.Send(
		sc.WebHook,
		fmt.Sprintf("Unexpected error in %s (%s)", sc.Application, env),
		nil,
		[]map[string]interface{}{
			{
				"fallback": fmt.Sprintf("Fatal error in %s (%s)", sc.Application, env),
				"color":    "#FF3232",
				"text":     fm,
			},
		},
	)

	log.Fatalf(fm)
}

func (sc *SlackConfig) GinFormatter(param gin.LogFormatterParams, m string) string {
	return fmt.Sprintf(
		"*%s* -> <%s>\n\n```%s```\n\n",
		param.Method,
		param.Path,
		m,
	)
}
