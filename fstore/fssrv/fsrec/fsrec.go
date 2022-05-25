/*
 * Copyright 2022 Oleg Borodin  <borodin@unix7.org>
 */

package fsrec

import (
    "fmt"
    "io/fs"
    "io"

    "ndstore/dscom"
    "ndstore/fstore/fssrv/fsfile"
)

const blockFileExt string = ".blk"
const storeDBName  string = "file.db"


type Store struct {
    dataRoot string
    dirPerm   fs.FileMode
    filePerm  fs.FileMode
    metaFile  *dscom.FileMI
}

func NewStore(dataRoot string) *Store {
    var store Store
    store.dataRoot  = dataRoot
    store.dirPerm   = 0755
    store.filePerm  = 0644
    return &store
}

func (store *Store) SetPerm(dirPerm, filePerm fs.FileMode) {
    store.dirPerm = dirPerm
    store.filePerm = filePerm
}

func (store *Store) SaveFile(fileName string, fileReader io.Reader, fileSize int64) error {
    var err error
    var fileId      int64 = 15
    var batchSize   int64 = 5
    var blockSize   int64 = 1024
    file := fsfile.NewFile(store.dataRoot, fileId, batchSize, blockSize)
    err = file.Open()
    defer file.Close()
    if err != nil {
        return err
    }

    _, err = file.Write(fileReader)
    if err != nil && err != io.EOF {
        return err
    }
    store.metaFile = file.Meta()
    return err
}

func (store *Store) FileExists(fileName string) (int64, error) {
    var err error
    var fileSize int64
    return fileSize, err
}

func (store *Store) LoadFile(fileName string, fileWriter io.Writer) error {
    var err error

    file := fsfile.RenewFile(store.dataRoot, store.metaFile)
    err = file.Open()
    defer file.Close()
    if err != nil {
        return err
    }
    _, err = file.Read(fileWriter)
    if err != nil {
        return err
    }
    return err
}

func (store *Store) DeleteFile(fileName string) error {
    var err error

    //var fileId      int64 = 15
    //var batchSize   int64 = 5
    //var blockSize   int64 = 1024
    //file := fsfile.NewFile(store.dataRoot, fileId, batchSize, blockSize)
    if store.metaFile == nil {
        fmt.Println("metafile is nil")
    }
    file := fsfile.RenewFile(store.dataRoot, store.metaFile)
    err = file.Open()
    defer file.Close()
    if err != nil {
        return err
    }
    err = file.Purge()
    if err != nil {
        return err
    }

    return err
}

func (store *Store) ListFiles(dirName string) ([]*dscom.DirEntry, error) {
    var err error
    files := make([]*dscom.DirEntry, 0)
    return files, err
}
