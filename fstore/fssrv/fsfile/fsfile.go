package fsfile

import (
    "io"
    "time"
    "dstore/dscomm/dsdescr"
    "dstore/dscomm/dsinter"
    "dstore/dscomm/dserr"
)

type File struct {
    reg             dsinter.FStoreReg
    baseDir         string

    login           string
    filePath        string
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

func NewFile(reg dsinter.FStoreReg, baseDir, login, filePath string, fileId, batchSize, blockSize int64) (*File, error) {
    var file File
    var err error
    file.reg        = reg
    file.baseDir    = baseDir

    file.login      = login
    file.filePath   = filePath

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

func OpenFile(reg dsinter.FStoreReg, baseDir, login, filePath string) (*File, error) {
    return openFile(false, reg, baseDir, login, filePath)
}

func ForceOpenFile(reg dsinter.FStoreReg, baseDir, login, filePath string) (*File, error) {
    return openFile(true, reg, baseDir, login, filePath)
}

func openFile(force bool, reg dsinter.FStoreReg, baseDir, login, filePath string) (*File, error) {
    var err error
    var file File
    file.reg        = reg
    file.baseDir    = baseDir

    descr, err := reg.GetFile(login, filePath)
    if err != nil {
        return &file, dserr.Err(err)
    }

    file.login      = descr.Login
    file.filePath   = descr.FilePath

    file.fileId     = descr.FileId
    file.batchSize  = descr.BatchSize
    file.blockSize  = descr.BlockSize
    file.dataSize   = descr.DataSize
    file.createdAt  = descr.CreatedAt
    file.updatedAt  = descr.UpdatedAt
    file.batchCount = descr.BatchCount

    file.batchs = make([]*Batch, file.batchCount)
    for i := int64(0); i < file.batchCount; i++ {
        switch {
            case force == false:
                batch, err := OpenBatch(reg, baseDir, i, file.fileId)
                if err != nil {
                    return &file, dserr.Err(err)
                }
                file.batchs[i] = batch
            default:
                batch, batchErr := ForceOpenBatch(reg, baseDir, i, file.fileId)
                if batchErr == nil {
                    file.batchs[i] = batch
                }
        }
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
        batch, err := NewBatch(file.reg, file.baseDir, batchNumber, file.fileId, file.batchSize, file.blockSize)
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
        file.dataSize += batchWritten
        if err != nil {
            return written, dserr.Err(err)
        }
        dataSize -= batchWritten
    }
    return written, dserr.Err(err)
}

func (file *File) Read(writer io.Writer) (int64, error) {
    var err error
    var readSize int64
    dataSize := file.dataSize
    for i := int64(0); i < file.batchCount; i++ {
        batchRead, err := file.batchs[i].Read(writer, dataSize)
        readSize += batchRead
        dataSize -= batchRead
        if err == io.EOF {
            err = nil
            return readSize, dserr.Err(err)
        }
        if err != nil {
            return readSize, dserr.Err(err)
        }
    }
    return readSize, dserr.Err(err)
}


func (file *File) Clean() error {
    var err error
    for i := file.batchCount - 1; i > -1; i-- {
        if file.batchs[i] != nil {
            err := file.batchs[i].Clean()
            if err != nil {
                return dserr.Err(err)
            }
            err = file.reg.DeleteBatch(i, file.fileId)
            if err != nil {
                return dserr.Err(err)
            }
            file.batchCount = i
            file.batchs[i] = nil
        }
    }
    file.dataSize = 0
    err = file.reg.DeleteFile(file.login, file.filePath)
    if err != nil {
        return dserr.Err(err)
    }
    return dserr.Err(err)
}

func (file *File) FileId() int64 {
    return file.fileId
}

func (file *File) DataSize() int64 {
    return file.dataSize
}


func (file *File) toDescr() *dsdescr.File {
    descr := dsdescr.NewFile()
    descr.Login         = file.login
    descr.FilePath      = file.filePath
    descr.FileId        = file.fileId
    descr.BatchSize     = file.batchSize
    descr.BlockSize     = file.blockSize
    descr.DataSize      = file.dataSize
    descr.BatchCount    = file.batchCount
    descr.CreatedAt     = file.createdAt
    descr.UpdatedAt     = file.updatedAt
    return descr
}
