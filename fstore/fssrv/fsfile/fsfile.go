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

func NewFile(baseDir string, reg dsinter.FStoreReg, login, filePath string, fileId, batchSize, blockSize int64) (*File, error) {
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

    return &file, dserr.Err(err)
}

func OpenFile(baseDir string, reg dsinter.FStoreReg, descr *dsdescr.File) (*File, error) {
    return openFile(false, baseDir, reg, descr)
}

func ForceOpenFile(baseDir string, reg dsinter.FStoreReg, descr *dsdescr.File) (*File, error) {
    return openFile(true, baseDir, reg, descr)
}

func openFile(force bool, baseDir string, reg dsinter.FStoreReg, descr *dsdescr.File) (*File, error) {
    var err error
    var file File
    file.reg        = reg
    file.baseDir    = baseDir

    file.login      = descr.Login
    file.filePath   = descr.FilePath

    file.fileId     = descr.FileId
    file.batchSize  = descr.BatchSize
    file.blockSize  = descr.BlockSize
    file.dataSize   = descr.DataSize
    file.createdAt  = descr.CreatedAt
    file.updatedAt  = descr.UpdatedAt
    file.batchCount = descr.BatchCount

    file.batchs = make([]*Batch, file.batchCount + 1)
    for i := int64(0); i < file.batchCount; i++ {
        switch {
            case force == false:
                batchDescr, err := file.reg.GetBatch(file.fileId, i)
                if err != nil {
                    return &file, dserr.Err(err)
                }
                batch, err := OpenBatch(baseDir, reg, batchDescr)
                if err != nil {
                    return &file, dserr.Err(err)
                }
                file.batchs[i] = batch
            default:
                batchDescr, getErr := file.reg.GetBatch(file.fileId, i)
                if getErr != nil {
                    continue
                }
                batch, batchErr := ForceOpenBatch(baseDir, reg, batchDescr)
                if batchErr != nil {
                    continue
                }
                file.batchs[i] = batch
        }
    }
    return &file, dserr.Err(err)
}

func (file *File) Write(reader io.Reader, dataSize int64) (int64, bool, error) {
    var err error
    var written int64
    var eof bool
    for i := range file.batchs {
        if dataSize < 1 {
            return written, eof, dserr.Err(err)
        }
        batchWritten, eof, err := file.batchs[i].Write(reader, dataSize)
        written += batchWritten

        if batchWritten > 0 {
            batchDescr := file.batchs[i].Descr()

            err = file.reg.PutBatch(batchDescr)
            if err != nil {
                return written, eof, dserr.Err(err)
            }
        }

        if err != nil {
            return written, eof, dserr.Err(err)
        }
        dataSize -= batchWritten
    }

    for {
        if dataSize < 1 {
            return written, eof, dserr.Err(err)
        }
        if eof {
            return written, eof, dserr.Err(err)
        }
        batchNumber := file.batchCount

        batch, err := NewBatch(file.baseDir, file.reg, file.fileId, batchNumber, file.batchSize, file.blockSize)
        if err != nil {
            return written, eof, dserr.Err(err)
        }
        batchDescr := batch.Descr()
        err = file.reg.PutBatch(batchDescr)
        if err != nil {
            return written, eof, dserr.Err(err)
        }

        batchWritten, eof, err := batch.Write(reader, dataSize)
        if err == io.EOF {
            err = nil
            eof = true
        }
        written += batchWritten
        file.dataSize += batchWritten
        batchDescr = batch.Descr()
        err = file.reg.PutBatch(batchDescr)
        if err != nil {
            return written, eof, dserr.Err(err)
        }
        dataSize -= batchWritten
        file.batchCount++
        if eof {
            return written, eof, dserr.Err(err)
        }
    }
    return written, eof, dserr.Err(err)
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
    return dserr.Err(err)
}

func (file *File) FileId() int64 {
    return file.fileId
}

func (file *File) DataSize() int64 {
    return file.dataSize
}

func (file *File) SetFilePath(filePath string) {
    file.filePath = filePath
}


func (file *File) Descr() *dsdescr.File {
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
