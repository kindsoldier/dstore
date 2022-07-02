package fsfile

import (
    "errors"
    "io"
    "time"
    "ndstore/dscom"
    "ndstore/dserr"
)

type File struct {
    reg             dscom.IFSReg
    baseDir         string

    fileId          int64
    fileVer         int64
    batchSize       int64
    blockSize       int64

    fileSize        int64
    createdAt       dscom.UnixTime
    updatedAt       dscom.UnixTime

    batchs          []*Batch

    fileIsDistr     bool
    fileIsOpen      bool
    fileIsDeleted   bool
}

func NewFile(reg dscom.IFSReg, baseDir string, batchSize, blockSize int64) (int64, *File, error) {
    var file File
    var err error

    file.reg        = reg
    file.baseDir    = baseDir

    fileId, err := reg.GetNewFileId()
    if err != nil {
        return fileId, &file, dserr.Err(err)
    }

    file.fileId     = fileId
    file.fileVer    = time.Now().UnixNano()
    file.batchSize  = batchSize
    file.blockSize  = blockSize

    file.batchs     = make([]*Batch, 0)
    file.fileSize   = 0
    file.createdAt  = dscom.UnixTime(time.Now().Unix())
    file.updatedAt  = dscom.UnixTime(time.Now().Unix())

    descr := file.toDescr()
    descr.UCounter = 2
    err = reg.AddNewFileDescr(descr)
    if err != nil {
        return fileId, &file, dserr.Err(err)
    }
    file.fileIsOpen     = true
    file.fileIsDeleted  = false
    file.fileIsDistr    = false

    return fileId, &file, dserr.Err(err)
}



func OpenSpecUnusedFile(reg dscom.IFSReg, baseDir string, fileId, fileVer int64) (*File, error) {
    var err error
    var file File

    file.reg        = reg
    file.baseDir    = baseDir

    exists, descr, err := reg.GetSpecUnusedFileDescr(fileId, fileVer)
    if err != nil {
        return &file, dserr.Err(err)
    }
    if !exists {
        err = errors.New("file not exists")
        return &file, dserr.Err(err)
    }
    file.fileId     = descr.FileId
    file.fileVer    = descr.FileVer

    file.batchSize  = descr.BatchSize
    file.blockSize  = descr.BlockSize
    file.fileSize   = descr.FileSize

    file.createdAt  = descr.CreatedAt
    file.updatedAt  = descr.UpdatedAt

    file.fileIsDistr  = false

    batchCount := descr.BatchCount

    file.batchs = make([]*Batch, batchCount)
    //for i := int64(0); i < batchCount; i++ {
    //    batchId := i
    //    batch, err := OpenBatch(reg, baseDir, fileId, batchId)
    //    if err != nil {
    //        return &file, dserr.Err(err)
    //    }
    //    file.batchs[i] = batch
    //}
    //err = file.reg.IncSpecFileDescrUC(file.fileId, file.fileVer)
    //if err != nil {
    //    return &file, dserr.Err(err)
    //}
    file.fileIsOpen     = true
    file.fileIsDeleted  = false
    file.fileIsDistr    = false

    return &file, dserr.Err(err)
}

func OpenFile(reg dscom.IFSReg, baseDir string, fileId int64) (*File, error) {
    var err error
    var file File

    file.reg        = reg
    file.baseDir    = baseDir

    exists, descr, err := reg.GetNewestFileDescr(fileId)
    if err != nil {
        return &file, dserr.Err(err)
    }
    if !exists {
        err = errors.New("file not exists")
        return &file, dserr.Err(err)
    }
    file.fileId     = descr.FileId
    file.fileVer    = descr.FileVer

    file.batchSize  = descr.BatchSize
    file.blockSize  = descr.BlockSize
    file.fileSize   = descr.FileSize

    file.createdAt  = descr.CreatedAt
    file.updatedAt  = descr.UpdatedAt

    file.fileIsDistr  = false

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
    err = file.reg.IncSpecFileDescrUC(file.fileId, file.fileVer)
    if err != nil {
        return &file, dserr.Err(err)
    }
    file.fileIsOpen     = true
    file.fileIsDeleted  = false
    file.fileIsDistr    = false

    return &file, dserr.Err(err)
}

func (file *File) Write(reader io.Reader, need int64) (int64, error) {
    var err error
    var written int64

    exitFunc := func() {
        if written > 0 {
            file.reg.DecSpecFileDescrUC(file.fileId, file.fileVer)
            file.reg.DecSpecFileDescrUC(file.fileId, file.fileVer)
            file.fileSize += written
            file.fileVer    = time.Now().UnixNano()
            file.updatedAt  = dscom.UnixTime(time.Now().Unix())
            descr := file.toDescr()
            descr.UCounter = 2
            file.reg.AddNewFileDescr(descr)
        }
    }
    defer exitFunc()

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

    if written > 0 {
        file.reg.DecSpecFileDescrUC(file.fileId, file.fileVer)

        file.fileSize += written
        file.fileVer    = time.Now().UnixNano()
        file.updatedAt  = dscom.UnixTime(time.Now().Unix())
        descr := file.toDescr()
        descr.UCounter = 2
        err = file.reg.AddNewFileDescr(descr)
        if err != nil {
            return written, dserr.Err(err)
        }
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
    if file.fileIsDeleted {
        return dserr.Err(err)
    }
    if file.fileIsOpen {
        err = file.reg.DecSpecFileDescrUC(file.fileId, file.fileVer)
        if err != nil {
                return dserr.Err(err)
        }
        file.fileIsOpen = false
    }
    // Delete batchs
    for i := int64(0); i < file.batchCount(); i++ {
        err := file.batchs[i].Delete()
        if err != nil {
            return dserr.Err(err)
        }
    }
    err = file.reg.DecSpecFileDescrUC(file.fileId, file.fileVer)
    if err != nil {
        return dserr.Err(err)
    }
    file.fileIsDeleted = true
    return dserr.Err(err)
}

func (file *File) Erase() error {
    var err error
    // Close file
    if file.fileIsOpen {
        err = file.reg.DecSpecFileDescrUC(file.fileId, file.fileVer)
        if err != nil {
                return dserr.Err(err)
        }
        file.fileIsOpen = false
    }
    if file.fileIsDeleted {
        return dserr.Err(err)
    }
    // Delete batchs
    //for i := int64(0); i < file.batchCount(); i++ {
    //    err := file.batchs[i].Erase()
    //    if err != nil {
    //        return dserr.Err(err)
    //    }
    //}
    file.batchs = make([]*Batch, 0)
    // Erase file descrs
    err = file.reg.EraseSpecFileDescr(file.fileId, file.fileVer)
    if err != nil {
        return dserr.Err(err)
    }
    file.fileIsDeleted = true
    return dserr.Err(err)
}


func (file *File) Close() error {
    var err error
    if file.fileIsDeleted {
        return dserr.Err(err)
    }
    if !file.fileIsOpen {
        return dserr.Err(err)
    }
    err = file.reg.DecSpecFileDescrUC(file.fileId, file.fileVer)
    if err != nil {
        return dserr.Err(err)
    }
    for i := int64(0); i < file.batchCount(); i++ {
        err := file.batchs[i].Close()
        if err != nil {
            return dserr.Err(err)
        }
    }
    file.fileIsOpen = false
    return dserr.Err(err)
}


func (file *File) Distribute(distr dscom.IBlockDistr) error {
    var err error
    if file.fileIsDeleted {
        return dserr.Err(err)
    }
    if !file.fileIsOpen {
        return dserr.Err(err)
    }
    if !file.fileIsOpen {
        return dserr.Err(err)
    }
    for i := int64(0); i < file.batchCount(); i++ {
        err := file.batchs[i].Distribute(distr)
        if err != nil {
            return dserr.Err(err)
        }
    }
    return dserr.Err(err)
}

func (file *File)  batchCount() int64 {
    return int64(len(file.batchs))
}

func (file *File) toDescr() *dscom.FileDescr {
    descr := dscom.NewFileDescr()
    descr.FileId        = file.fileId
    descr.FileVer       = file.fileVer
    descr.BatchSize     = file.batchSize
    descr.BlockSize     = file.blockSize
    descr.FileSize      = file.fileSize
    descr.BatchCount    = file.batchCount()
    descr.CreatedAt     = file.createdAt
    descr.UpdatedAt     = dscom.UnixTime(time.Now().Unix())
    return descr
}
