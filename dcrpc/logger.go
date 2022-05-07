/*
 * Copyright 2022 Oleg Borodin  <borodin@unix7.org>
 */

package dcrpc

import (
    "fmt"
    "io"
    "os"
    "time"
)

var messageWriter io.Writer = os.Stdout
var accessWriter io.Writer = os.Stdout

func logDebug(messages ...interface{}) {
    stamp := time.Now().Format(time.RFC3339Nano)
    fmt.Fprintln(messageWriter, stamp, "debug", messages)
}

func logInfo(messages ...interface{}) {
    stamp := time.Now().Format(time.RFC3339Nano)
    fmt.Fprintln(messageWriter, stamp, "info", messages)
}

func logError(messages ...interface{}) {
    stamp := time.Now().Format(time.RFC3339Nano)
    fmt.Fprintln(messageWriter, stamp, "error", messages)
}

func logAccess(messages ...interface{}) {
    stamp := time.Now().Format(time.RFC3339Nano)
    fmt.Fprintln(messageWriter, stamp, messages)
}

func SetAccessWriter(writer io.Writer) {
    accessWriter = writer
}

func SetMessageWriter(writer io.Writer) {
    messageWriter = writer
}
