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
    "dstore/dscomm/dslog"
)

func (store *Store) SaveFile(login string, filePath string, fileReader io.Reader, fileSize int64) (*dsdescr.File, error) {
    var err error
    var has bool
    var descr *dsdescr.File

    filePath = cleanPath(filePath)
    has, err = store.reg.HasFile(login, filePath)
    if err != nil {
        return descr, dserr.Err(err)
    }
    if has {
        descr, err = store.reg.GetFile(login, filePath)
        if err != nil {
            return descr, dserr.Err(err)
        }
        err = fmt.Errorf("file %s already exist", filePath)
        return descr, dserr.Err(err)
    }

    var batchSize   int64 = 5
    var blockSize   int64 = 1024 * 1024 * 8

    if fileSize < blockSize * batchSize {
        blockSize = fileSize / batchSize
        rs := int64(1024 * 16)
        bs := blockSize / rs
        blockSize = (bs + 1) * rs
    }

    fileId, err := store.fileAlloc.NewId()
    if err != nil {
        return descr, dserr.Err(err)
    }
    file, err := fsfile.NewFile(store.reg, store.dataDir, login, filePath, fileId, batchSize, blockSize)
    if err != nil {
        return descr, dserr.Err(err)
    }
    wrSize, err := file.Write(fileReader, fileSize)
    if err == io.EOF {
        return descr, dserr.Err(err)
    }
    if err != nil  {
        return descr, dserr.Err(err)
    }
    if wrSize != fileSize {
        return descr, dserr.Err(err)
    }

    has, err = store.reg.HasFile(login, filePath)
    if err != nil {
        return descr, dserr.Err(err)
    }
    if !has {
        err = fmt.Errorf("file %s not saved", filePath)
        return descr, dserr.Err(err)
    }
    descr, err = store.reg.GetFile(login, filePath)
    if err != nil {
        return descr, dserr.Err(err)
    }
    return descr, dserr.Err(err)
}

func (store *Store) HasFile(login string, filePath string) (bool, *dsdescr.File, error) {
    var err error
    var has bool
    var descr *dsdescr.File
    filePath = cleanPath(filePath)
    has, err = store.reg.HasFile(login, filePath)
    if err != nil {
        return has, descr, dserr.Err(err)
    }
    if !has {
        err = fmt.Errorf("file %s not exist", filePath)
        return has, descr, dserr.Err(err)
    }
    descr, err = store.reg.GetFile(login, filePath)
    if err != nil {
        return has, descr, dserr.Err(err)
    }
    return has, descr, dserr.Err(err)
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

func (store *Store) DeleteFile(login string, filePath string) (*dsdescr.File, error) {
    var err error
    var descr *dsdescr.File
    filePath = cleanPath(filePath)
    has, err := store.reg.HasFile(login, filePath)
    if err != nil {
        return descr, dserr.Err(err)
    }
    if !has {
        return descr, dserr.Err(err)
    }
    descr, err = store.reg.GetFile(login, filePath)
    if err != nil {
        return descr, dserr.Err(err)
    }

    file, err := fsfile.ForceOpenFile(store.reg, store.dataDir, login, filePath)
    if err != nil {
        return descr, dserr.Err(err)
    }
    err = file.Clean()
    if err != nil {
        return descr, dserr.Err(err)
    }
    store.fileAlloc.FreeId(file.FileId())
    if err != nil {
        return descr, dserr.Err(err)
    }
    allocJSON, _ := store.fileAlloc.JSON()
    dslog.LogDebugf("file id alloc state: %s", string(allocJSON))
    return descr, dserr.Err(err)
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
