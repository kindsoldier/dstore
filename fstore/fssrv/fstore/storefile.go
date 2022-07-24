/*
 * Copyright 2022 Oleg Borodin  <borodin@unix7.org>
 */

package fstore

import (
    "fmt"
    "io"
    "path/filepath"
    "regexp"

    "github.com/ganbarodigital/go_glob"

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

    if fileSize < blockSize * batchSize {
        batchSize = fileSize / blockSize + 1
    }


    fileId, err := store.fileAlloc.NewId()
    if err != nil {
        return descr, dserr.Err(err)
    }
    file, err := fsfile.NewFile(store.dataDir, store.reg, login, filePath, fileId, batchSize, blockSize)
    if err != nil {
        return descr, dserr.Err(err)
    }

    descr = file.Descr()
    err = store.reg.PutFile(descr)
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

    descr = file.Descr()
    err = store.reg.PutFile(descr)
    if err != nil {
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
    descr, err := store.reg.GetFile(login, filePath)
    if err != nil {
        return dserr.Err(err)
    }
    file, err := fsfile.OpenFile(store.dataDir, store.reg, descr)
    if err != nil {
        return dserr.Err(err)
    }
    _, err = file.Read(fileWriter)
    if err != nil {
        return dserr.Err(err)
    }
    return dserr.Err(err)
}

func (store *Store) FileStats(login, pattern, regular, gPattern string) (int64, int64, error) {
    var err error
    var usage int64
    var count int64
    descrs, err := store.ListFiles(login, pattern, regular, gPattern)
    if err != nil {
        return count, usage, err
    }
    for _, descr := range descrs {
        usage += descr.DataSize
        count++
    }
    return count, usage, err
}

func (store *Store) EraseFiles(login, pattern, regular, gPattern string, erase bool, reader io.Reader) ([]*dsdescr.File, error) {
    var err error

    descrs, err := store.listFiles(login, pattern, regular, gPattern)
    if !erase {
        return descrs, dserr.Err(err)
    }

    notify := make(chan error, 16)

    checker := func() {
        buf := make([]byte, 1)
        for {
            _, err := reader.Read(buf)
            if err != nil {
                dslog.LogDebugf("user %s connection err: %s", login, err)
                notify <- err
                return
            }
        }
    }
    go checker()

    resDescrs := make([]*dsdescr.File, 0)
    for _, descr := range descrs {

        select {
            case err := <-notify:
                err = fmt.Errorf("connection error: %s", err)
                return resDescrs, dserr.Err(err)
            default:
        }

        file, err := fsfile.ForceOpenFile(store.dataDir, store.reg, descr)
        if err != nil {
            err = fmt.Errorf("cannot open file %s, err: %v", descr.FilePath, err)
            return resDescrs, dserr.Err(err)
        }
        err = file.Clean()
        if err != nil {
            return resDescrs, dserr.Err(err)
        }
        err = store.reg.DeleteFile(descr.Login, descr.FilePath)
        if err != nil {
            err = fmt.Errorf("cannot delete file descr for %s, err: %v", descr.FilePath, err)
            return resDescrs, dserr.Err(err)
        }

        store.fileAlloc.FreeId(file.FileId())
        if err != nil {
            return resDescrs, dserr.Err(err)
        }

        dslog.LogDebugf("user %s file deleted: %s", login, descr.FilePath)

        //allocJSON, _ := store.fileAlloc.JSON()
        //dslog.LogDebugf("file id alloc state: %s", string(allocJSON))
        resDescrs = append(resDescrs, descr)
    }
    return resDescrs, dserr.Err(err)
}

func (store *Store) ListFiles(login, pattern, regular, gPattern string) ([]*dsdescr.File, error) {
    return store.listFiles(login, pattern, regular, gPattern)
}

func (store *Store) listFiles(login, pattern, regular, gPattern string) ([]*dsdescr.File, error) {
    var err error
    resDescrs := make([]*dsdescr.File, 0)
    descrs, err := store.reg.ListFiles(login)
    if err != nil {
        return resDescrs, dserr.Err(err)
    }

    usePattern  := false
    useRegular  := false
    useGPattern := false
    if len(pattern) > 0 {
        usePattern = true
    }
    if len(regular) > 0 {
        useRegular = true
    }
    if len(gPattern) > 0 {
        useGPattern = true
    }

    if !usePattern && !useRegular && !useGPattern {
        resDescrs = descrs
        return resDescrs, dserr.Err(err)
    }

    pattern = "/" + pattern
    pattern = filepath.Clean(pattern)

    re, err := regexp.CompilePOSIX(regular)
    if err != nil {
        return resDescrs, dserr.Err(err)
    }
    g := glob.NewGlob(gPattern)

    for _, descr := range descrs {
        ok1 := false
        ok2 := false
        ok3 := false
        switch useRegular {
            case true:
                ok1 = re.Match([]byte(descr.FilePath))
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
        switch useGPattern {
            case true:
                ok3, err = g.Match(descr.FilePath)
                if err != nil {
                    return resDescrs, dserr.Err(err)
                }
            default:
                ok3 = true
        }

        if ok1 && ok2 && ok3 {
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
    file, err := fsfile.ForceOpenFile(store.dataDir, store.reg, descr)
    if err != nil {
        return descr, dserr.Err(err)
    }
    err = file.Clean()
    if err != nil {
        return descr, dserr.Err(err)
    }
    err = store.reg.DeleteFile(login, filePath)
    if err != nil {
        return descr, dserr.Err(err)
    }

    store.fileAlloc.FreeId(file.FileId())
    if err != nil {
        return descr, dserr.Err(err)
    }
    //allocJSON, _ := store.fileAlloc.JSON()
    //dslog.LogDebugf("file id alloc state: %s", string(allocJSON))

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
