/*
 * Copyright 2022 Oleg Borodin  <borodin@unix7.org>
 */

package fsrec

import (
    "errors"
    "io/fs"
    "io"
    "path/filepath"

    "ndstore/dscom"
    "ndstore/fstore/fssrv/fsfile"
    "ndstore/fstore/fssrv/fsreg"
)

const blockFileExt string = ".blk"
const storeDBName  string = "file.db"


type Store struct {
    dataRoot string
    dirPerm   fs.FileMode
    filePerm  fs.FileMode
    reg    *fsreg.Reg
}

func NewStore(dataRoot string, reg *fsreg.Reg) *Store {
    var store Store
    store.dataRoot  = dataRoot
    store.dirPerm   = 0755
    store.filePerm  = 0644
    store.reg       = reg
    return &store
}

func (store *Store) SetPerm(dirPerm, filePerm fs.FileMode) {
    store.dirPerm = dirPerm
    store.filePerm = filePerm
}

func (store *Store) SaveFile(filePath string, fileReader io.Reader, fileSize int64) error {
    var err error

    fileId, err := store.reg.GetNewFileId()
    if err != nil {
        return err
    }

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
    meta, err := file.Meta()
    if err != nil {
        return err
    }

    err = store.reg.AddFileDescr(meta)
    if err != nil {
        return err
    }

    fileName := filepath.Base(filePath)
    dirPath := filepath.Dir(filePath)

    err = store.reg.AddEntryDescr(dirPath, fileName, fileId)
    if err != nil {
        return err
    }
    return err
}

func (store *Store) FileExists(fileName string) (int64, error) {
    var err error
    var fileSize int64
    return fileSize, err
}

func (store *Store) LoadFile(filePath string, fileWriter io.Writer) error {
    var err error

    fileName := filepath.Base(filePath)
    dirPath := filepath.Dir(filePath)

    entry, exists, err := store.reg.GetEntryDescr(dirPath, fileName)
    if !exists || entry == nil {
        return errors.New("path not exists")
    }

    meta, err := store.reg.GetFileDescr(entry.FileId)
    if err != nil {
        return err
    }

    file := fsfile.RenewFile(store.dataRoot, meta)
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

func (store *Store) DeleteFile(filePath string) error {
    var err error

    fileName := filepath.Base(filePath)
    dirPath := filepath.Dir(filePath)

    entry, exists, err := store.reg.GetEntryDescr(dirPath, fileName)
    if !exists || entry == nil {
        return errors.New("path not exists")
    }

    fileId := entry.FileId
    meta, err := store.reg.GetFileDescr(fileId)
    if err != nil {
        return err
    }

    file := fsfile.RenewFile(store.dataRoot, meta)
    err = file.Open()
    defer file.Close()
    if err != nil {
        return err
    }

    err = file.Purge()
    if err != nil {
        return err
    }

    err = store.reg.DeleteEntryDescr(dirPath, fileName)
    if err != nil {
        return err
    }

    err = store.reg.DeleteFileDescr(fileId)
    if err != nil {
        return err
    }

    return err
}

func (store *Store) ListFiles(dirName string) ([]*dscom.EntryDescr, error) {
    var err error
    files := make([]*dscom.EntryDescr, 0)
    return files, err
}
