/*
 * Copyright 2022 Oleg Borodin  <borodin@unix7.org>
 */

package fsrec

import (
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
    file.Open()
    defer file.Close()

    file.Write(fileReader)
    return err
}

func (store *Store) FileExists(fileName string) (int64, error) {
    var err error
    var fileSize int64
    return fileSize, err
}

func (store *Store) LoadFile(fileName string, fileWriter io.Writer) error {
    var err error

    var fileId      int64 = 15
    var batchSize   int64 = 5
    var blockSize   int64 = 1024
    file := fsfile.NewFile(store.dataRoot, fileId, batchSize, blockSize)
    file.Open()
    defer file.Close()

    file.Read(fileWriter)
    return err
}

func (store *Store) DeleteFile(fileName string) error {
    var err error
    return err
}

func (store *Store) ListFiles(dirName string) ([]*dscom.CFile, error) {
    var err error
    files := make([]*dscom.CFile, 0)
    return files, err
}
