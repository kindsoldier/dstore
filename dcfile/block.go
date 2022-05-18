package dcfile

import (
    "errors"
    "fmt"
    "io"
    "io/fs"
    "os"
    "path/filepath"
)


type Block struct {
    file        *os.File
    baseDir     string
    fileId      int64
    batchId     int64
    blockId     int64
    capacity    int64
    //recorded    int
}

const fileMode fs.FileMode = 0644

func NewBlock(baseDir string, fileId, batchId, blockId int64, capacity int64) *Block {
    var block Block
    block.baseDir   = baseDir
    block.fileId    = fileId
    block.batchId   = batchId
    block.blockId   = blockId
    block.capacity  = capacity

    return &block
}

func (block *Block) Open() error {
    var err error

    fileName := fmt.Sprintf("%02d-%02d-%02d.bin", block.fileId, block.batchId, block.blockId)
    filePath := filepath.Join(block.baseDir, fileName)

    openMode := os.O_APPEND | os.O_CREATE | os.O_RDWR
    file, err := os.OpenFile(filePath, openMode, fileMode)
    if err != nil {
            return err
    }
    block.file = file

    return err
}

func (block *Block) Write(reader io.Reader) (int64, error) {
    var err error
    var recorded int64
    stat, err := block.file.Stat()
    if err != nil {
            return recorded, err
    }
    size := stat.Size()
    remains := block.capacity - size
    recorded, err = copyBytes(reader, block.file, remains)
    if err != nil {
            return recorded, err
    }
    return recorded, err
}

func (block *Block) Read(writer io.Writer) (int64, error) {
    var err error
    var read int64
    stat, err := block.file.Stat()
    if err != nil {
            return read, err
    }
    size := stat.Size()
    read, err = copyBytes(block.file, writer, size)
    if err != nil {
            return read, err
    }
    return read, err
}

func (block *Block) Close() error {
    return block.file.Close()
}


func (block *Block) Truncate() error {
    var err error
    err = block.file.Truncate(0)
    if err != nil {
            return err
    }
    _, err = block.file.Seek(0,0)
    if err != nil {
            return err
    }
    return err
}

func copyBytes(reader io.Reader, writer io.Writer, size int64) (int64, error) {
    var err error
    var bufSize int64 = 1024 * 4
    var total   int64 = 0
    var remains int64 = size
    buffer := make([]byte, bufSize)

    for {
        if remains == 0 {
            return total, err
        }
        if remains < bufSize {
            bufSize = remains
        }
        received, err := reader.Read(buffer[0:bufSize])
        if err != nil {
            return total, err
        }
        recorded, err := writer.Write(buffer[0:received])
        if err != nil {
            return total, err
        }
        if recorded != received {
            return total, errors.New("write error")
        }
        total += int64(recorded)
        remains -= int64(recorded)
    }
    return total, err
}
