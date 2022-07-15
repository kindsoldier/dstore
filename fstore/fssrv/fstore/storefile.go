/*
 * Copyright 2022 Oleg Borodin  <borodin@unix7.org>
 */

package fstore

import (
    "fmt"
    "io"
    "path/filepath"
    "regexp"
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

func (store *Store) FileStats(login, pattern, regular string) (int64, int64, error) {
    var err error
    var usage int64
    var count int64
    descrs, err := store.ListFiles(login, pattern, regular)
    if err != nil {
        return count, usage, err
    }
    for _, descr := range descrs {
        usage += descr.DataSize
        count++
    }
    return count, usage, err
}

func (store *Store) ListFiles(login, pattern, regular string) ([]*dsdescr.File, error) {
    return store.listFiles(login, pattern, regular)
}

func (store *Store) listFiles(login, pattern, regular string) ([]*dsdescr.File, error) {
    var err error
    resDescrs := make([]*dsdescr.File, 0)
    descrs, err := store.reg.ListFiles(login)
    if err != nil {
        return resDescrs, dserr.Err(err)
    }

    usePattern := false
    useRegular := false
    if len(pattern) > 0 {
        usePattern = true
    }
    if len(regular) > 0 {
        useRegular = true
    }

    dslog.LogDebug("use pattern:", pattern, regular)

    dslog.LogDebug("use pattern:", usePattern, useRegular)

    if !usePattern && !useRegular {
        resDescrs = descrs
        return resDescrs, dserr.Err(err)
    }

    pattern = "/" + pattern
    pattern = filepath.Clean(pattern)

    re, err := regexp.CompilePOSIX(regular)
    if err != nil {
        return resDescrs, dserr.Err(err)
    }

    for _, descr := range descrs {
        ok1 := false
        ok2 := false
        switch useRegular {
            case true:
                ok1 = re.Match([]byte(descr.FilePath))
                if err != nil {
                    return resDescrs, dserr.Err(err)
                }
            default:
                ok1 = true
        }
        switch usePattern {
            case true:
                ok2, err = filepath.Match(pattern, descr.FilePath)
                if err != nil {
                    return resDescrs, dserr.Err(err)
                }
            default:
                ok2 = true
        }
        if ok1 && ok2 {
            dslog.LogDebug("res ok:", ok1, ok2)
            resDescrs = append(resDescrs, descr)
        }
    }
    return resDescrs, dserr.Err(err)
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
