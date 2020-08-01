# Watcher

A simple and efficient package to handle staging and production errors.

```cgo
go get github.com/Alvarios/watcher
```

## Prerequisite

You need a [Slack application](https://api.slack.com/apps) that supports [incoming webhooks](https://api.slack.com/messaging/webhooks#:~:text=Incoming%20Webhooks%20are%20a%20simple,make%20the%20messages%20stand%20out.).

## Configuration

Create a SlackConfig in your go code.

```go
package myPackage

import "github.com/Alvarios/watcher"

var Watcher watcher.SlackConfig

func main() {
    Watcher := watcher.SlackConfig{
        WebHook: "https://hooks.slack.com/services/XXXXX/XXXXX/XXXXX",
        Application: "Your application name",
        Environment: "Your environment",
    }
}
```

## Methods

### Error

Send a basic error message to Slack.

```go
package myPackage

func MyFunction() {
	_, err := Watcher.Error("some error message")
}
```

### Fatal

Same as Error method, but kills the running server and doesn't return any status.

```go
package myPackage

func MyFunction() {
	Watcher.Fatal("some fatal message")
}
```

### GinAbort

A special message formatter for Gin, that also interrupts the request.

```go
package myPackage

import "github.com/gin-gonic/gin"

func MyMiddleware(c *gin.Context) {
	Watcher.GinAbort(c, "some error message")
}
```

## Copyright
2020 Kushuh - [MIT license](https://github.com/Alvarios/watcher/blob/master/LICENSE)