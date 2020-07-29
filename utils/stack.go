package utils

import (
	"fmt"
	"runtime"
)

func PrintStack(skip int) string {
	stackSlice := make([]byte, 512)
	s := runtime.Stack(stackSlice, false)
	return fmt.Sprintf("\n%s", stackSlice[skip:s])
}
