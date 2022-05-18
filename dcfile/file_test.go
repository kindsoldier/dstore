package dcfile

import (
    "io"
    "testing"
)

func TestFile(t * testing.T) {
    var writer io.Writer
    writer = NewFile()

    writer.Write([]byte("qwerty"))

    var reader io.Reader
    reader = NewFile()

    buffer := make([]byte, 12)
    reader.Read(buffer)
}
