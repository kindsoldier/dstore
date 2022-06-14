/*
 * Copyright 2022 Oleg Borodin  <borodin@unix7.org>
 */


package fsfile

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

    "ndstore/dscom"
    "ndstore/dsrpc"
    "ndstore/bstore/bsfunc"
    "ndstore/dserr"
)


type Block struct {
    file        *os.File
    baseDir     string
    fileId      int64
    batchId     int64
    blockId     int64
    blockSize   int64

    remote      bool
    local       bool
    storeId     int64

    hasher      hash.Hash
    hashSum     []byte
    hashInit    []byte
}

const fileMode fs.FileMode = 0644

var ErrorNilFile = errors.New("block file ref is nil")

func NewBlock(baseDir string, fileId, batchId, blockId int64, blockSize int64) *Block {
    var block Block
    block.baseDir   = baseDir
    block.fileId    = fileId
    block.batchId   = batchId
    block.blockId   = blockId
    block.blockSize = blockSize
    block.hashSum   = make([]byte, 0)

    hashInit := make([]byte, 32)
    rand.Read(hashInit)
    hasher, _ := highwayhash.New(hashInit)

    block.hashInit  = hashInit
    block.hasher    = hasher
    return &block
}

func (block *Block) Meta() (*dscom.BlockDescr, error) {
    var err error
    meta := dscom.NewBlockDescr()
    meta.FileId     = block.fileId
    meta.BatchId    = block.batchId
    meta.BlockId    = block.blockId
    meta.BlockSize  = block.blockSize

    meta.FilePath   = block.fileName()
    meta.DataSize, err  = block.Size()
    if err != nil {
            return meta, dserr.Err(err)
    }
    meta.HashInit   = hex.EncodeToString(block.hashInit)
    meta.HashSum    = hex.EncodeToString(block.hashSum)
    return meta, dserr.Err(err)
}

func (block *Block) Open() error {
    var err error

    filePath := block.filePath()
    openMode := os.O_APPEND|os.O_CREATE|os.O_RDWR
    file, err := os.OpenFile(filePath, openMode, fileMode)
    if err != nil {
            return dserr.Err(err)
    }
    block.file = file
    return dserr.Err(err)
}

func (block *Block) Size() (int64, error) {
    var err error
    var size int64
    if block.file == nil {
        return size, ErrorNilFile
    }
    stat, err := block.file.Stat()
    if err != nil {
            return size, dserr.Err(err)
    }
    size = stat.Size()
    return size, dserr.Err(err)
}


func (block *Block) Write(reader io.Reader) (int64, error) {
    var err error
    var written int64
    if block.file == nil {
        return written, ErrorNilFile
    }
    size, err := block.Size()
    if err != nil {
            return written, dserr.Err(err)
    }
    remains := block.blockSize - size
    multiWriter := io.MultiWriter(block.file, block.hasher)
    written, err = copyBytes(reader, multiWriter, remains)
    block.hashSum = block.hasher.Sum(nil)
    if err != nil {
            return written, dserr.Err(err)
    }
    return written, dserr.Err(err)
}

func (block *Block) Lwrite(reader io.Reader, need int64) (int64, error) {
    var err error
    var written int64
    if block.file == nil {
        return written, ErrorNilFile
    }
    if need < 1 {
        return written, io.EOF
    }

    size, err := block.Size()
    if err != nil {
            return written, dserr.Err(err)
    }
    remains := block.blockSize - size
    if need < remains {
        remains = need
    }
    multiWriter := io.MultiWriter(block.file, block.hasher)
    written, err = copyBytes(reader, multiWriter, remains)
    block.hashSum = block.hasher.Sum(nil)
    if err != nil {
            return written, dserr.Err(err)
    }
    need -= written
    if need < 1 {
            return written, io.EOF
    }
    return written, dserr.Err(err)
}

func (block *Block) Save(pool dscom.IBSPool) error {
    var err error

    err = block.ToBegin()
    if err != nil {
        return dserr.Err(err)
    }

    fileId  := block.fileId
    batchId := block.batchId
    blockId := block.blockId
    blockSize   := block.blockSize
    blockReader := block.file
    dataSize, err  := block.Size()
    if err != nil {
        return dserr.Err(err)
    }
    blockType   := dscom.BTypeData
    hashAlg     := dscom.HashTypeHW
    hashInit    := hex.EncodeToString(block.hashInit)
    hashSum     := hex.EncodeToString(block.hashSum)

    storeId, err := pool.SaveBlock(fileId, batchId, blockId, blockSize, blockReader,
                                                dataSize, blockType, hashAlg, hashInit, hashSum)

    if err != nil {
        return dserr.Err(err)
    }
    block.remote = true
    block.storeId = storeId
    return dserr.Err(err)
}

func (block *Block) Read(writer io.Writer) (int64, error) {
    var err error
    var read int64
    if block.file == nil {
        return read, ErrorNilFile
    }
    //size, err := block.Size()
    //if err != nil {
            //return read, dserr.Err(err)
    //}
    //read, err = copyBytes(block.file, writer, size)
    //if err != nil {
    //        return read, dserr.Err(err)
    //}

    fileId      := block.fileId
    batchId     := block.batchId
    blockId     := block.blockId
    blockWriter := writer //io.Discard
    blockType   := dscom.BTypeData

    uri         := "localhost:5101"
    auth        := dsrpc.CreateAuth([]byte("admin"), []byte("admin"))

    _ = bsfunc.LoadBlock(uri, auth, fileId, batchId, blockId, blockWriter, blockType)

    return read, dserr.Err(err)
}

func (block *Block) Close() error {
    var err error
    if block.file == nil {
        return dserr.Err(err)
    }
    err = block.file.Close()
    return dserr.Err(err)
}

func (block *Block) Purge() error {
    var err error
    block.Close()
    filePath := block.filePath()
    err = os.Remove(filePath)
    if err != nil {
        return dserr.Err(err)
    }
    return dserr.Err(err)
}

func (block *Block) Truncate() error {
    var err error
    if block.file == nil {
        return ErrorNilFile
    }
    err = block.file.Truncate(0)
    if err != nil {
        return dserr.Err(err)
    }
    _, err = block.file.Seek(0,0)
    if err != nil {
        return dserr.Err(err)
    }
    return dserr.Err(err)
}

func (block *Block) Seek(offset int64) error {
    _, err := block.file.Seek(offset, 0)
    return dserr.Err(err)
}

func (block *Block) ToBegin() error {
    _, err := block.file.Seek(0, 0)
    return dserr.Err(err)
}

func (block *Block) ToEnd() error {
    _, err := block.file.Seek(0, 2)
    return dserr.Err(err)
}


func copyBytes(reader io.Reader, writer io.Writer, size int64) (int64, error) {
    var err error
    var bufSize int64 = 1024 * 4
    var total   int64 = 0
    var remains int64 = size
    buffer := make([]byte, bufSize)

    for {
        if remains == 0 {
            return total, dserr.Err(err)
        }
        if remains < bufSize {
            bufSize = remains
        }
        received, err := reader.Read(buffer[0:bufSize])
        if err != nil {
            return total, dserr.Err(err)
        }
        written, err := writer.Write(buffer[0:received])
        if err != nil {
            return total, dserr.Err(err)
        }
        if written != received {
            err = errors.New("write error")
            return total, dserr.Err(err)
        }
        total += int64(written)
        remains -= int64(written)
    }
    return total, dserr.Err(err)
}

func (block *Block) fileName() string {
    fileName := fmt.Sprintf("%04d-%04d-%03d.blk", block.fileId, block.batchId, block.blockId)
    return fileName
}

func (block *Block) filePath() string {
    filePath := filepath.Join(block.baseDir, block.fileName())
    return filePath
}
