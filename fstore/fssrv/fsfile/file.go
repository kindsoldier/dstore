package fsfile

import (
    "io"
    "ndstore/dscom"
)

type File struct {
    baseDir     string
    fileId      int64
    batchSize   int64
    blockSize   int64
    batchs      []*Batch
}

func NewFile(baseDir string, fileId, batchSize, blockSize int64) *File {
    var file File
    file.batchs     = make([]*Batch, 0)
    file.fileId     = fileId
    file.batchSize  = batchSize
    file.blockSize  = blockSize
    file.baseDir    = baseDir
    return &file
}

func RenewFile(baseDir string, meta *dscom.FileDescr) *File {
    var file File
    file.baseDir    = baseDir
    file.fileId     = meta.FileId
    file.batchSize  = meta.BatchSize
    file.blockSize  = meta.BlockSize
    file.batchs     = make([]*Batch, 0)
    for i := int64(0); i < meta.BatchCount; i++ {
        batch := NewBatch(file.baseDir, file.fileId, i, file.batchSize, file.blockSize)
        file.batchs = append(file.batchs, batch)
    }
    return &file
}

func (file *File) Meta() (*dscom.FileDescr, error) {
    var err error
    fileMeta := dscom.NewFileDescr()
    fileMeta.FileId     = file.fileId
    fileMeta.BatchCount = file.batchCount()
    fileMeta.BatchSize  = file.batchSize
    fileMeta.BlockSize  = file.blockSize
    fileMeta.FileSize, err = file.Size()
    if err != nil {
        return fileMeta, err
    }

    for i := range file.batchs {
        batchMeta, err := file.batchs[i].Meta()
        if err != nil {
            return fileMeta, err
        }
        fileMeta.Batchs = append(fileMeta.Batchs, batchMeta)
    }
    return fileMeta, err
}

func (file *File) Open() error {
    var err error
    for i := range file.batchs {
        err = file.batchs[i].Open()
        if err != nil {
            return err
        }
    }
    return err
}

func (file *File) Write(reader io.Reader) (int64, error) {
    var err error
    var written int64

    for i := range file.batchs {
        batchWritten, err := file.batchs[i].Write(reader)
        written += batchWritten
        if err != nil {
            return written, err
        }
    }
    for {
        batchNumber := file.batchCount()

        batch := NewBatch(file.baseDir, file.fileId, batchNumber, file.batchSize, file.blockSize)
        file.batchs = append(file.batchs, batch)
        err = batch.Open()
        if err != nil {
            return written, err
        }
        batchWritten, err := batch.Write(reader)
        written += batchWritten
        if err != nil {
            return written, err
        }
    }
    return written, err
}

func (file *File) Read(writer io.Writer) (int64, error) {
    var err error
    var read int64
    for i := int64(0); i < file.batchCount(); i++ {
        blockRead, err := file.batchs[i].Read(writer)
        read += blockRead
        if err != nil {
            return read, err
        }
    }
    return read, err
}

func (file *File) Close() error {
    var err error
    for i := int64(0); i < file.batchCount(); i++ {
        err := file.batchs[i].Close()
        if err != nil {
            return err
        }
    }
    return err
}

func (file *File) Truncate() error {
    var err error
    for i := int64(0); i < file.batchCount(); i++ {
        err := file.batchs[i].Truncate()
        if err != nil {
            return err
        }
    }
    return err
}

func (file *File) Purge() error {
    var err error
    for i := int64(0); i < file.batchCount(); i++ {
        err := file.batchs[i].Purge()
        if err != nil {
            return err
        }
    }
    return err
}

func (file *File) Size() (int64, error) {
    var err error
    var size int64
    for i := int64(0); i < file.batchCount(); i++ {
        blockSize, err := file.batchs[i].Size()
        size += blockSize
        if err != nil {
            return size, err
        }
    }
    return size, err
}

func (file *File) ToBegin() error {
    var err error
    for i := int64(0); i < file.batchCount(); i++ {
        err := file.batchs[i].ToBegin()
        if err != nil {
            return err
        }
    }
    return err
}

func (file *File) ToEnd() error {
    var err error
    for i := int64(0); i < file.batchCount(); i++ {
        err := file.batchs[i].ToEnd()
        if err != nil {
            return err
        }
    }
    return err
}

func (file *File) batchCount() int64 {
    return int64(len(file.batchs))
}
