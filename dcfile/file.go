package dcfile

import (
    "io"
    "fmt"
)

type File struct {
    baseDir     string
    fileId      int64
    batchSize   int64
    blockSize   int64
    batchs      []*Batch
}

func NewFile(basedir string, fileId, batchSize, blockSize int64) *File {
    var file File
    file.batchs = make([]*Batch, 0)
    file.fileId = fileId
    file.batchSize = batchSize
    file.blockSize = blockSize
    return &file
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
        fmt.Println(batchNumber, written)
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

func (file *File) batchCount() int64 {
    return int64(len(file.batchs))
}
