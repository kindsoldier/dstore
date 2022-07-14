/*
 * Copyright 2022 Oleg Borodin  <borodin@unix7.org>
 */

package fstore

import (
    "fmt"
    "io"
    "path/filepath"
    "dstore/fstore/fssrv/fsfile"
    "dstore/dscomm/dserr"
    "dstore/dscomm/dsdescr"
)

func (store *Store) SaveFile(login string, filePath string, fileReader io.Reader, fileSize int64) error {
    var err error
    var has bool

    filePath = cleanPath(filePath)
    has, err = store.reg.HasFile(login, filePath)
    if err != nil {
        return dserr.Err(err)
    }
    if has {
        err = fmt.Errorf("file %s already exist", filePath)
        return dserr.Err(err)
    }

    var batchSize   int64 = 5
    var blockSize   int64 = 1000 * 1000

    if fileSize < blockSize * batchSize {
        blockSize = fileSize / batchSize
        rs := int64(1024 * 10)
        bs := blockSize / rs
        blockSize = (bs + 1) * rs
    }

    fileId, err := store.fileAlloc.NewId()
    if err != nil {
        return dserr.Err(err)
    }
    file, err := fsfile.NewFile(store.reg, store.dataDir, login, filePath, fileId, batchSize, blockSize)
    if err != nil {
        return dserr.Err(err)
    }
    wrSize, err := file.Write(fileReader, fileSize)
    if err == io.EOF {
        return dserr.Err(err)
    }
    if err != nil  {
        return dserr.Err(err)
    }
    if wrSize != fileSize {
        return dserr.Err(err)
    }
    return dserr.Err(err)
}

func (store *Store) HasFile(login string, filePath string) (bool, int64, error) {
    var err error
    var has bool
    var fileSize int64
    filePath = cleanPath(filePath)
    has, err = store.reg.HasFile(login, filePath)
    if err != nil {
        return has, fileSize, dserr.Err(err)
    }
    if !has {
        err = fmt.Errorf("file %s not exist", filePath)
        return has, fileSize, dserr.Err(err)
    }
    descr, err := store.reg.GetFile(login, filePath)
    if err != nil {
        return has, fileSize, dserr.Err(err)
    }
    fileSize = descr.DataSize
    return has, fileSize, dserr.Err(err)
}

func (store *Store) LoadFile(login string, filePath string, fileWriter io.Writer) error {
    var err error
    filePath = cleanPath(filePath)
    has, err := store.reg.HasFile(login, filePath)
    if err != nil {
        return dserr.Err(err)
    }
    if !has {
        err = fmt.Errorf("file %s not exist", filePath)
        return dserr.Err(err)
    }
    file, err := fsfile.OpenFile(store.reg, store.dataDir, login, filePath)
    if err != nil {
        return dserr.Err(err)
    }
    _, err = file.Read(fileWriter)
    if err != nil {
        return dserr.Err(err)
    }
    return dserr.Err(err)
}


func (store *Store) ListFiles(login string, dirPath string) ([]*dsdescr.File, error) {
    var err error
    files := make([]*dsdescr.File, 0)
    files, err = store.reg.ListFiles(login)
    if err != nil {
        return files, dserr.Err(err)
    }
    return files, dserr.Err(err)
}

func (store *Store) DeleteFile(login string, filePath string) error {
    var err error
    filePath = cleanPath(filePath)
    has, err := store.reg.HasFile(login, filePath)
    if err != nil {
        return dserr.Err(err)
    }
    if !has {
        err = fmt.Errorf("file %s not exist", filePath)
        return dserr.Err(err)
    }
    file, err := fsfile.OpenFile(store.reg, store.dataDir, login, filePath)
    if err != nil {
        return dserr.Err(err)
    }
    err = file.Clean()
    if err != nil {
        return dserr.Err(err)
    }
    store.fileAlloc.FreeId(file.FileId())
    if err != nil {
        return dserr.Err(err)
    }
    return dserr.Err(err)
}

func (store *Store) checkLogin(login string) error {
    var err error
    var has bool
    has, err = store.reg.HasUser(login)
    if err != nil {
        return dserr.Err(err)
    }
    if !has {
        err = fmt.Errorf("user %s not exist", login)
        return dserr.Err(err)
    }
    return dserr.Err(err)
}

func cleanPath(filePath string) string {
    filePath = "/" + filePath
    filePath = filepath.Clean(filePath)
    return filePath
}
