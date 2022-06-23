package fsfile

import (
    "errors"
    "io"
    "ndstore/dscom"
    "ndstore/dserr"
)

type File struct {
    reg         dscom.IFSReg
    baseDir     string
    fileId      int64
    batchSize   int64
    blockSize   int64
    fileSize    int64
    batchs      []*Batch

    openedWOErrors  bool
    fileIsClosed    bool
    fileIsErased    bool
}

func NewFile(reg dscom.IFSReg, baseDir string, batchSize, blockSize int64) (int64, *File, error) {
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
    file.fileId = fileId
    file.openedWOErrors = true
    return fileId, &file, dserr.Err(err)
}

func OpenFile(reg dscom.IFSReg, baseDir string, fileId int64) (*File, error) {
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
    file.openedWOErrors = true
    return &file, dserr.Err(err)
}

func (file *File) Write(reader io.Reader, need int64) (int64, error) {
    var err error
    var written int64

    updater := func() {
        file.fileSize += written
        file.updateFileDescr()
    }
    defer updater()

    for i := range file.batchs {
        if need < 1 {
            return written, dserr.Err(err)
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
    file.fileIsErased = true
    return dserr.Err(err)
}


func (file *File) BrutalErase() error {
    var err error
    fineClean := true
    blockDescrs, _ := file.reg.ListBlockDescrsByFileId(file.fileId)
    for _, descr := range blockDescrs {
        block, err := OpenBlock(file.reg, file.baseDir, descr.FileId, descr.BatchId, descr.BlockId, descr.BlockType)
        if err == nil && block != nil {
            err = block.Erase()
            if err != nil {
                fineClean = false
            }
            block.Close()
        }
    }
    batchDescrs, _ := file.reg.ListBatchDescrsByFileId(file.fileId)
    for _, descr := range batchDescrs {
        batch, err := OpenBatch(file.reg, file.baseDir, descr.FileId, descr.BatchId)
        if err == nil && batch != nil {
            err = batch.Erase()
            if err != nil {
                fineClean = false
            }
            batch.Close()
        }
        err = file.reg.EraseBatchDescr(descr.FileId, descr.BatchId)
        if err != nil {
            fineClean = false
        }
    }
    if fineClean {
        err = file.reg.EraseFileDescr(file.fileId)
        if err != nil {
            fineClean = false
        }
    }
    file.fileIsErased = true
    return dserr.Err(err)
}


func (file *File) Close() error {
    var err error
    if file.fileIsErased {
        return dserr.Err(err)
    }
    if file.fileIsClosed {
        return dserr.Err(err)
    }
    if !file.openedWOErrors {
        return dserr.Err(err)
    }
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
    descr.UCounter = 2
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
