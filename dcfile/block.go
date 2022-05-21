package dcfile

import (
    "errors"
    "encoding/hex"
    "hash"
    "fmt"
    "io"
    "io/fs"
    "os"
    "path/filepath"
    "math/rand"

    "github.com/minio/highwayhash"
)

type BlockMeta struct {
    FileName      string
}

func NewBlockMeta() *BlockMeta {
    var blockMeta BlockMeta
    return &blockMeta
}


type Block struct {
    file        *os.File
    baseDir     string
    fileId      int64
    batchId     int64
    blockId     int64
    blockSize    int64
    hasher      hash.Hash
    hexHash     []byte
}

const fileMode fs.FileMode = 0644
var ErrorNilFile = errors.New("block file ref is nil")

func NewBlock(baseDir string, fileId, batchId, blockId int64, blockSize int64) *Block {
    var block Block
    block.baseDir   = baseDir
    block.fileId    = fileId
    block.batchId   = batchId
    block.blockId   = blockId
    block.blockSize  = blockSize
    block.hexHash   = make([]byte, 0)

    initBytes := make([]byte, 32)
    rand.Read(initBytes)

    hasher, _ := highwayhash.New(initBytes)
    block.hasher = hasher

    return &block
}

func (block *Block) Meta() *BlockMeta {
    meta := NewBlockMeta()
    meta.FileName = block.filePath()
    return meta
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
    remains := block.blockSize - size
    multiWriter := io.MultiWriter(block.file, block.hasher)
    written, err = copyBytes(reader, multiWriter, remains)
    if err != nil {
            return written, err
    }

    hashBytes := block.hasher.Sum(nil)
    hexBytes := make([]byte, hex.EncodedLen(len(hashBytes)))
    hex.Encode(hexBytes, hashBytes)
    block.hexHash = hexBytes
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

func (block *Block) Seek(offset int64) error {
    _, err := block.file.Seek(offset, 0)
    return err
}

func (block *Block) ToBegin() error {
    _, err := block.file.Seek(0, 0)
    return err
}

func (block *Block) ToEnd() error {
    _, err := block.file.Seek(0, 2)
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
    fileName := fmt.Sprintf("%020d-%00d-%020d.blk", block.fileId, block.batchId, block.blockId)
    filePath := filepath.Join(block.baseDir, fileName)
    return filePath
}
