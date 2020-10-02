package setup

import (
	"fmt"
	"github.com/Alvarios/go-slack"
	"github.com/gin-gonic/gin"
	"log"
	"net/http"
	"os"
	"runtime/debug"
	"time"
)

type SlackConfig struct {
	WebHook string
	Application string
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
				"fallback" : fmt.Sprintf("Unexpected error in %s (%s)", sc.Application, env),
				"color" : "#FF9300",
				"text" : sc.Print(m),
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
				"fallback" : fmt.Sprintf("Fatal error in %s (%s)", sc.Application, env),
				"color" : "#FF3232",
				"text" : fm,
			},
		},
	)

	log.Fatalf(fm)
}

func (sc *SlackConfig) FormatForGin(c *gin.Context, m string) string {
	return fmt.Sprintf(
		"*%s* -> <http://%s%s>\n\n```%s```\n\n",
		c.Request.Method,
		c.Request.Host,
		c.Request.URL.Path,
		m,
	)
}

func (sc *SlackConfig) GinAbort(c *gin.Context, m string) {
	_, _ = sc.Error(sc.FormatForGin(c, m))

	c.AbortWithStatusJSON(
		http.StatusInternalServerError,
		gin.H{"message" : m},
	)
}