package fsfile

import (
    "errors"
    "io"
    "ndstore/fstore/fssrv/fsreg"
    "ndstore/dscom"
    "ndstore/dserr"
    "ndstore/dslog"
)

type File struct {
    reg         *fsreg.Reg
    baseDir     string
    fileId      int64
    batchSize   int64
    blockSize   int64
    fileSize    int64
    batchs      []*Batch
}

func NewFile(reg *fsreg.Reg, baseDir string, batchSize, blockSize int64) (int64, *File, error) {
    var fileId int64
    var file File
    var err error
    file.reg        = reg
    file.batchSize  = batchSize
    file.blockSize  = blockSize
    file.baseDir    = baseDir
    file.batchs     = make([]*Batch, 0)

    fileId, err = file.addFileDescr()
    if err != nil {
        return fileId, &file, dserr.Err(err)
    }
    file.fileId     = fileId
    return fileId, &file, dserr.Err(err)
}

func OpenFile(reg *fsreg.Reg, baseDir string, fileId int64) (*File, error) {
    var err error
    var file File
    exists, descr, err := reg.GetFileDescr(fileId)
    if err != nil {
        return &file, dserr.Err(err)
    }
    if !exists {
        err = errors.New("file not exists")
        return &file, dserr.Err(err)
    }
    file.reg       = reg
    file.baseDir   = baseDir
    file.fileId    = descr.FileId
    file.batchSize = descr.BatchSize
    file.blockSize = descr.BlockSize
    file.fileSize  = descr.FileSize

    batchCount := descr.BatchCount

    file.batchs = make([]*Batch, batchCount)
    for i := int64(0); i < batchCount; i++ {
        batchId := i
        batch, err := OpenBatch(reg, baseDir, fileId, batchId)
        if err != nil {
            return &file, dserr.Err(err)
        }
        file.batchs[i] = batch
    }

    err = file.reg.IncFileDescrUC(file.fileId)
    if err != nil {
        return &file, dserr.Err(err)
    }
    return &file, dserr.Err(err)
}


func (file *File) Write(reader io.Reader, need int64) (int64, error) {
    var err error
    var written int64

    updater := func() {
        file.fileSize += written
        dslog.LogDebug("file size", file.fileSize)
        file.updateFileDescr()
    }
    defer updater()

    for i := range file.batchs {
        if need < 1 {
            return written, io.EOF
        }
        batchWritten, err := file.batchs[i].Write(reader, need)
        written += batchWritten
        if err != nil {
            return written, dserr.Err(err)
        }
        need -= batchWritten
    }
    for {
        if need < 1 {
            return written, dserr.Err(err)
        }
        batchNumber := file.batchCount()
        batch, err := NewBatch(file.reg, file.baseDir, file.fileId, batchNumber, file.batchSize, file.blockSize)
        if err != nil {
            return written, dserr.Err(err)
        }
        file.batchs = append(file.batchs, batch)
        if err != nil {
            return written, dserr.Err(err)
        }
        batchWritten, err := batch.Write(reader, need)
        written += batchWritten
        if err != nil {
            return written, dserr.Err(err)
        }
        need -= batchWritten
    }
    return written, dserr.Err(err)
}

func (file *File) Read(writer io.Writer) (int64, error) {
    var err error
    var read int64
    for i := int64(0); i < file.batchCount(); i++ {
        blockRead, err := file.batchs[i].Read(writer)
        read += blockRead
        if err != nil {
            return read, dserr.Err(err)
        }
    }
    return read, dserr.Err(err)
}

func (file *File) Delete() error {
    var err error
    err = file.reg.DecFileDescrUC(file.fileId)
    if err != nil {
        return dserr.Err(err)
    }
    return dserr.Err(err)
}

func (file *File) Erase() error {
    var err error
    for i := int64(0); i < file.batchCount(); i++ {
        err := file.batchs[i].Erase()
        if err != nil {
            return dserr.Err(err)
        }
    }
    file.batchs = make([]*Batch, 0)
    err = file.reg.EraseFileDescr(file.fileId)
    if err != nil {
        return dserr.Err(err)
    }
    return dserr.Err(err)
}

func (file *File) Close() error {
    var err error
    err = file.reg.DecFileDescrUC(file.fileId)
    if err != nil {
        return dserr.Err(err)
    }
    for i := int64(0); i < file.batchCount(); i++ {
        err := file.batchs[i].Close()
        if err != nil {
            return dserr.Err(err)
        }
    }
    return dserr.Err(err)
}

func (file *File)  batchCount() int64 {
    return int64(len(file.batchs))
}

func (batch *File) addFileDescr() (int64, error) {
    descr := batch.toDescr()
    descr.UCounter = 1
    return batch.reg.AddFileDescr(descr)
}

func (batch *File) updateFileDescr() error {
    descr := batch.toDescr()
    return batch.reg.UpdateFileDescr(descr)
}


func (file *File) toDescr() *dscom.FileDescr {
    descr := dscom.NewFileDescr()
    descr.FileId    = file.fileId
    descr.BatchSize = file.batchSize
    descr.BlockSize = file.blockSize
    descr.FileSize  = file.fileSize
    descr.BatchCount = file.batchCount()
    return descr
}
