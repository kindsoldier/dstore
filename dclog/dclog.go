package dclog

import (
    "errors"
    "fmt"
    "io"
    "os"
    "time"

    "github.com/sirupsen/logrus"
)

type logFormatter struct {
}

func (f *logFormatter) Format(entry *logrus.Entry) ([]byte, error) {
    var err error
    timeStamp := time.Now().Format(time.RFC3339Nano)
    levelString := entry.Level.String()
    message := fmt.Sprintf("%s %s %s\n", timeStamp, levelString, entry.Message)
    return []byte(message), err
}

const DebugLevel int = 1
const ErrorLevel int = 2

func init() {
    logrus.SetOutput(os.Stdout)
    logrus.SetLevel(logrus.DebugLevel)
    logrus.SetFormatter(new(logFormatter))
}

func SetLevel(level int) error {
    var err error
    switch level {
        case DebugLevel:
            logrus.SetLevel(logrus.DebugLevel)
        case ErrorLevel:
            logrus.SetLevel(logrus.DebugLevel)
        default:
            return errors.New("wrong log level")
    }
    return err
}

func SetOutput(writer io.Writer) error {
    var err error
    logrus.SetOutput(writer)
    return err
}

func LogDebug(message ...interface{}) {
        logrus.Debug(message)
}

func LogError(message ...interface{}) {
        logrus.Error(message)
}

func LogWarning(message ...interface{}) {
        logrus.Warning(message)
}

func LogInfo(message ...interface{}) {
        logrus.Info(message)
}
