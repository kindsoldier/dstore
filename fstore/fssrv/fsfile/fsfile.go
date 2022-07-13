package fsfile

import (
    "io"
    "time"
    "dstore/dsdescr"
    "dstore/dsinter"
    "dstore/dserr"
)

type File struct {
    reg             dsinter.StoreReg
    baseDir         string

    fileId          int64
    fileVer         int64
    batchSize       int64
    blockSize       int64

    dataSize        int64
    createdAt       int64
    updatedAt       int64
    batchCount      int64
    batchs          []*Batch
}

func NewFile(baseDir string, reg dsinter.StoreReg, fileId, batchSize, blockSize int64) (*File, error) {
    var file File
    var err error
    file.reg        = reg
    file.baseDir    = baseDir

    file.fileId     = fileId
    file.batchSize  = batchSize
    file.blockSize  = blockSize
    file.dataSize   = 0
    file.createdAt  = time.Now().Unix()
    file.updatedAt  = file.createdAt
    file.batchs     = make([]*Batch, 0)

    descr := file.toDescr()
    err = reg.PutFile(descr)
    if err != nil {
        return &file, dserr.Err(err)
    }
    return &file, dserr.Err(err)
}

func OpenFile(baseDir string, reg dsinter.StoreReg, fileId int64) (*File, error) {
    var err error
    var file File
    file.reg        = reg
    file.baseDir    = baseDir

    descr, err := reg.GetFile(fileId)
    if err != nil {
        return &file, dserr.Err(err)
    }
    file.fileId     = descr.FileId
    file.batchSize  = descr.BatchSize
    file.blockSize  = descr.BlockSize
    file.dataSize   = descr.DataSize
    file.createdAt  = descr.CreatedAt
    file.updatedAt  = descr.UpdatedAt
    file.batchCount = descr.BatchCount

    file.batchs = make([]*Batch, file.batchCount)
    for i := int64(0); i < file.batchCount; i++ {
        batch, err := OpenBatch(baseDir, reg, i, file.fileId)
        if err != nil {
            return &file, dserr.Err(err)
        }
        file.batchs[i] = batch
    }
    return &file, dserr.Err(err)
}

func (file *File) Write(reader io.Reader, dataSize int64) (int64, error) {
    var err error
    var written int64

    existFunc := func() {
        if written > 0 {
            descr := file.toDescr()
            err = file.reg.PutFile(descr)
        }
    }
    defer existFunc()

    for i := range file.batchs {
        if dataSize < 1 {
            return written, dserr.Err(err)
        }
        batchWritten, err := file.batchs[i].Write(reader, dataSize)
        written += batchWritten
        if err != nil {
            return written, dserr.Err(err)
        }
        dataSize -= batchWritten
    }
    for {
        if dataSize < 1 {
            return written, dserr.Err(err)
        }

        batchNumber := file.batchCount
        batch, err := NewBatch(file.baseDir, file.reg, batchNumber, file.fileId, file.batchSize, file.blockSize)
        if err != nil {
            return written, dserr.Err(err)
        }
        file.batchs = append(file.batchs, batch)
        file.batchCount++
        if err != nil {
            return written, dserr.Err(err)
        }
        batchWritten, err := batch.Write(reader, dataSize)
        written += batchWritten
        if err != nil {
            return written, dserr.Err(err)
        }
        dataSize -= batchWritten
    }
    return written, dserr.Err(err)
}

func (file *File) Read(writer io.Writer) (int64, error) {
    var err error
    var read int64
    for i := int64(0); i < file.batchCount; i++ {
        blockRead, err := file.batchs[i].Read(writer)
        read += blockRead
        if err != nil {
            return read, dserr.Err(err)
        }
    }
    return read, dserr.Err(err)
}


func (file *File) Clean() error {
    var err error
    for i := int64(0); i < file.batchCount; i++ {
        batch := file.batchs[i]
        if batch != nil {
            err := batch.Clean()
            if err != nil {
                return dserr.Err(err)
            }
        }
    }
    descr := file.toDescr()
    err = file.reg.PutFile(descr)
    if err != nil {
        return dserr.Err(err)
    }
    return dserr.Err(err)
}

func (file *File) toDescr() *dsdescr.File {
    descr := dsdescr.NewFile()
    descr.FileId        = file.fileId
    descr.BatchSize     = file.batchSize
    descr.BlockSize     = file.blockSize
    descr.DataSize      = file.dataSize
    descr.BatchCount    = file.batchCount
    descr.CreatedAt     = file.createdAt
    descr.UpdatedAt     = file.updatedAt
    return descr
}
