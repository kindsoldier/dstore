package fsfile

import (
    "fmt"
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
    createdAt       int64
    updatedAt       int64
    isDistr         bool
    isCompl         bool
    isDeleted       bool

    batchs          []*Batch

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
    file.createdAt  = time.Now().Unix()
    file.updatedAt  = file.createdAt
    file.isDistr    = false
    file.isCompl    = false
    file.isDeleted  = false

    descr := file.toDescr()
    descr.UCounter = 2
    err = reg.AddNewFileDescr(descr)
    if err != nil {
        return fileId, &file, dserr.Err(err)
    }
    file.fileIsOpen     = true
    file.fileIsDeleted  = false

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
        err = fmt.Errorf("file %s not exists", file.getIdString())
        return &file, dserr.Err(err)
    }
    file.fileId     = descr.FileId
    file.fileVer    = descr.FileVer

    file.batchSize  = descr.BatchSize
    file.blockSize  = descr.BlockSize
    file.fileSize   = descr.FileSize

    file.createdAt  = descr.CreatedAt
    file.updatedAt  = descr.UpdatedAt

    file.isDistr    = descr.IsDistr
    file.isCompl    = descr.IsCompl
    file.isDeleted  = descr.IsDeleted

    batchCount := descr.BatchCount

    file.batchs = make([]*Batch, batchCount)

    file.fileIsOpen     = true
    file.fileIsDeleted  = false

    return &file, dserr.Err(err)
}

func OpenFile(reg dscom.IFSReg, baseDir string, fileId int64) (*File, error) {
    force := false
    return openFile(force, reg, baseDir, fileId)
}


func ForceOpenFile(reg dscom.IFSReg, baseDir string, fileId int64) (*File, error) {
    force := true
    return openFile(force, reg, baseDir, fileId)
}

func openFile(force bool, reg dscom.IFSReg, baseDir string, fileId int64) (*File, error) {
    var err error
    var file File


    file.reg        = reg
    file.baseDir    = baseDir

    exists, descr, err := reg.GetNewestFileDescr(fileId)
    if err != nil {
        return &file, dserr.Err(err)
    }
    if !exists {
        err = fmt.Errorf("file %s not exists", file.getIdString())
        return &file, dserr.Err(err)
    }
    file.fileId     = descr.FileId
    file.fileVer    = descr.FileVer

    file.batchSize  = descr.BatchSize
    file.blockSize  = descr.BlockSize
    file.fileSize   = descr.FileSize

    file.createdAt  = descr.CreatedAt
    file.updatedAt  = descr.UpdatedAt

    file.isDistr    = descr.IsDistr

    batchCount := descr.BatchCount

    file.batchs = make([]*Batch, batchCount)
    for i := int64(0); i < batchCount; i++ {
        batchId := i
        batch, err := ForceOpenBatch(reg, baseDir, fileId, batchId)
        if !force {
            if err != nil {
                batch.Close()
                return &file, dserr.Err(err)
            }
        }
        if batch != nil {
            file.batchs[i] = batch
        }
    }
    err = file.reg.IncSpecFileDescrUC(1, file.fileId, file.fileVer)
    if err != nil {
        return &file, dserr.Err(err)
    }
    file.fileIsOpen     = true
    file.fileIsDeleted  = false

    return &file, dserr.Err(err)
}

func (file *File) Write(reader io.Reader, need int64) (int64, error) {
    var err error
    var written int64

    exitFunc := func(file *File, err error) {
        if written > 0 {
            descr := file.toDescr()
            file.fileSize += written
            file.fileVer    = time.Now().UnixNano()
            file.updatedAt  = time.Now().Unix()
            if err == nil {
                file.isCompl = true
                file.isDistr = false
            }
            newDescr := file.toDescr()
            newDescr.UCounter = 2
            file.reg.AddNewFileDescr(newDescr)
            file.reg.DecSpecFileDescrUC(2, descr.FileId, descr.FileVer)
        }
    }
    defer exitFunc(file, err)

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
    //if written > 0 {
    //    // Save prev state to descr
    //    descr := file.toDescr()
    //    // Update state
    //    file.fileSize += written // TODO: change to file.size()
    //    file.fileVer    = time.Now().UnixNano()
    //    file.updatedAt  = time.Now().Unix()
    //    // Save new state to descr
    //    newDescr := file.toDescr()
    //    descr.UCounter = 2
    //    // Add new version of descr
    //    err = file.reg.AddNewFileDescr(newDescr)
    //    if err != nil {
    //        return written, dserr.Err(err)
    //    }
    //    // Descrease usage counter of old state
    //    file.reg.DecSpecFileDescrUC(2, descr.FileId, descr.FileVer)
    //    if err != nil {
    //        return written, dserr.Err(err)
    //    }
    //}
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
        err = file.reg.DecSpecFileDescrUC(1, file.fileId, file.fileVer)
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
    err = file.reg.DecSpecFileDescrUC(1, file.fileId, file.fileVer)
    if err != nil {
        return dserr.Err(err)
    }
    file.fileIsDeleted = true
    return dserr.Err(err)
}

func (file *File) Erase() error {
    var err error
    // The same as close file
    if file.fileIsOpen {
        err = file.reg.DecSpecFileDescrUC(1, file.fileId, file.fileVer)
        if err != nil {
                return dserr.Err(err)
        }
        file.fileIsOpen = false
    }
    // Do nothing if yet erased
    if file.fileIsDeleted {
        return dserr.Err(err)
    }
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
    // Always close batchs
    for i := int64(0); i < file.batchCount(); i++ {
        batch := file.batchs[i]
        if batch != nil {
            err := batch.Close()
            if err != nil {
                return dserr.Err(err)
            }
        }
    }
    if !file.fileIsOpen {
        return dserr.Err(err)
    }
    err = file.reg.DecSpecFileDescrUC(1, file.fileId, file.fileVer)
    if err != nil {
        return dserr.Err(err)
    }
    file.fileIsOpen = false
    return dserr.Err(err)
}


func (file *File) Distribute(distr dscom.IFileDistr) error {
    var err error
    // For dummy call: do nothing if deleted
    if file.fileIsDeleted {
        return dserr.Err(err)
    }
    // Return with err if not open
    if !file.fileIsOpen {
        err = fmt.Errorf("file %s is not open or open with errors", file.getIdString())
        return dserr.Err(err)
    }
    // Distibute all underline objects
    allBatchIsDisr := true
    for i := int64(0); i < file.batchCount(); i++ {
        batchIsDisr, err := file.batchs[i].Distribute(distr)
        if err != nil {
            return dserr.Err(err)
        }
        if !batchIsDisr {
            allBatchIsDisr = false
        }
    }
    // Add new ver of file descr if distr state was changed
    if allBatchIsDisr != file.isDistr {
        // Save prev state to descr
        descr := file.toDescr()
        // Update state
        file.isDistr    = allBatchIsDisr
        file.fileVer    = time.Now().UnixNano()
        file.updatedAt  = time.Now().Unix()
        // Save new state to descr
        newDescr := file.toDescr()
        newDescr.UCounter = 2
        // Add new version of descr
        err = file.reg.AddNewFileDescr(newDescr)
        if err != nil {
            return dserr.Err(err)
        }
        // Descrease usage counter of old state
        file.reg.DecSpecFileDescrUC(2, descr.FileId, descr.FileVer)
        if err != nil {
            return dserr.Err(err)
        }
    }
    return dserr.Err(err)
}

func (file *File) getIdString() string {
    return fmt.Sprintf("%d,%d", file.fileId, file.fileVer)
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
    descr.UpdatedAt     = time.Now().Unix()
    descr.IsDistr       = file.isDistr
    descr.IsCompl       = file.isCompl
    descr.IsDeleted     = file.isDeleted
    return descr
}
