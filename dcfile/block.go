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
    //written    int
}

const fileMode fs.FileMode = 0644
var ErrorNilFile = errors.New("block file ref is nil")

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

    filePath := block.filePath()
    openMode := os.O_APPEND|os.O_CREATE|os.O_RDWR
    file, err := os.OpenFile(filePath, openMode, fileMode)
    if err != nil {
            return err
    }
    block.file = file
    return err
}

func (block *Block) Size() (int64, error) {
    var err error
    var size int64
    if block.file == nil {
        return size, ErrorNilFile
    }
    stat, err := block.file.Stat()
    if err != nil {
            return size, err
    }
    size = stat.Size()
    return size, err
}


func (block *Block) Write(reader io.Reader) (int64, error) {
    var err error
    var written int64
    if block.file == nil {
        return written, ErrorNilFile
    }
    size, err := block.Size()
    if err != nil {
            return written, err
    }
    remains := block.capacity - size
    written, err = copyBytes(reader, block.file, remains)
    if err != nil {
            return written, err
    }
    return written, err
}

func (block *Block) Read(writer io.Writer) (int64, error) {
    var err error
    var read int64
    if block.file == nil {
        return read, ErrorNilFile
    }
    size, err := block.Size()
    if err != nil {
            return read, err
    }
    read, err = copyBytes(block.file, writer, size)
    if err != nil {
            return read, err
    }
    return read, err
}

func (block *Block) Close() error {
    var err error
    if block.file == nil {
        return err
    }
    err = block.file.Close()
    return err
}

func (block *Block) Purge() error {
    var err error
    block.Close()
    filePath := block.filePath()
    err = os.Remove(filePath)
    if err != nil {
        return err
    }

    return err
}

func (block *Block) Truncate() error {
    var err error
    if block.file == nil {
        return ErrorNilFile
    }
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
        written, err := writer.Write(buffer[0:received])
        if err != nil {
            return total, err
        }
        if written != received {
            return total, errors.New("write error")
        }
        total += int64(written)
        remains -= int64(written)
    }
    return total, err
}

func (block *Block) filePath() string {
    fileName := fmt.Sprintf("%02d-%02d-%02d.bin", block.fileId, block.batchId, block.blockId)
    filePath := filepath.Join(block.baseDir, fileName)
    return filePath
}
