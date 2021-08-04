package panic

import (
	"bytes"
	"fmt"
	"runtime/debug"

	"github.com/marprin/postman-lib/pkg/env"
	"github.com/sirupsen/logrus"
)

// ToPanicError log recovered panic error
func ToPanicError(r interface{}, message string) {
	var buff bytes.Buffer
	errStr := fmt.Sprint(r)
	errStack := string(debug.Stack())

	buff.WriteString(fmt.Sprintf("[PANIC] (%s) %s", env.Get(), errStr))
	if message != "" {
		buff.WriteString(fmt.Sprintf(" ```%s``` ", message))
	}

	logrus.Error(buff.String())
	logrus.Error(errStack)
}

// HandlePanic handle when panic ocured
func HandlePanic(handler func(interface{})) {
	if r := recover(); r != nil {
		handler(r)
	}
}
