package setup

import (
	"fmt"
	"github.com/Alvarios/go-slack"
	"github.com/Alvarios/watcher/utils"
	"github.com/gin-gonic/gin"
	"log"
	"net/http"
	"time"
)

type SlackConfig struct {
	WebHook string
	Application string
	Environment string
}

func (sc *SlackConfig) Print(m string) string {
	stack := utils.PrintStack(2)

	return fmt.Sprintf(
		"%s\n\n*File*\n%s\n\n*Environment*\n%s\n\n*Noticed*\n%s",
		m,
		stack,
		sc.Environment,
		time.Now().Format("2006-01-02 3:4:5"),
	)
}

func (sc *SlackConfig) Error(m string) (*http.Response, error) {
	return slack.Send(
		sc.WebHook,
		fmt.Sprintf("Unexpected error in %s", sc.Application),
		nil,
		[]map[string]interface{}{
			{
				"fallback" : fmt.Sprintf("Unexpected error in %s (%s)", sc.Application, sc.Environment),
				"color" : "#FF9300",
				"text" : sc.Print(m),
			},
		},
	)
}

func (sc *SlackConfig) Fatal(m string) {
	fm := sc.Print(m)

	_, _ = slack.Send(
		sc.WebHook,
		fmt.Sprintf("Unexpected error in %s", sc.Application),
		nil,
		[]map[string]interface{}{
			{
				"fallback" : fmt.Sprintf("Fatal error in %s (%s)", sc.Application, sc.Environment),
				"color" : "#FF9300",
				"text" : fm,
			},
		},
	)

	log.Fatalf(fm)
}

func (sc *SlackConfig) GinAbort(c *gin.Context, m string) {
	fm := fmt.Sprintf(
		"*%s* -> <http://%s%s>\n\n```%s```\n\n",
		c.Request.Method,
		c.Request.Host,
		c.Request.URL.Path,
		m,
	)

	_, _ = sc.Error(fm)

	c.AbortWithStatusJSON(
		http.StatusInternalServerError,
		gin.H{"message" : fm},
	)
}