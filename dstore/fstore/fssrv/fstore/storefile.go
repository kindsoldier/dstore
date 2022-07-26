/*
 * Copyright 2022 Oleg Borodin  <borodin@unix7.org>
 */

package fstore

import (
    "fmt"
    "io"
    "path/filepath"
    "regexp"
    "math/rand"
    "encoding/hex"

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

    // Get file id
    fileId, err := store.fileAlloc.NewId()
    if err != nil {
        return descr, dserr.Err(err)
    }

    // Create tmp name
    randBin := make([]byte, 16)
    rand.Read(randBin)
    randStr := hex.EncodeToString(randBin)
    tmpFilePath := filepath.Join("/.tmp/", randStr, filePath)

    // Create file object
    file, err := fsfile.NewFile(store.dataDir, store.reg, login, tmpFilePath, fileId, batchSize, blockSize)
    if err != nil {
        return descr, dserr.Err(err)
    }
    // Save file descr with tmp name
    descr = file.Descr()
    err = store.reg.PutFile(descr)
    if err != nil {
        return descr, dserr.Err(err)
    }
    _, eof, err := file.Write(fileReader, fileSize)
    if err == io.EOF {
        err = nil
        eof = true
    }
    if err != nil   {
        dslog.LogDebugf("write error %s,%s: %v", login, filePath, err)
    }
    if eof {
        dslog.LogDebugf("eof for %s,%s", login, filePath)
    }
    // Save descr with new name
    file.SetFilePath(filePath)
    descr = file.Descr()

    err = store.reg.PutFile(descr)
    if err != nil {
        return descr, dserr.Err(err)
    }
    // Delete old descr with tmp name
    err = store.reg.DeleteFile(login, tmpFilePath)
    if err != nil {
        return descr, dserr.Err(err)
    }
    // Check descr
    has, err = store.reg.HasFile(login, filePath)
    if err != nil {
        return descr, dserr.Err(err)
    }
    if !has {
        err = fmt.Errorf("file %s not saved", filePath)
        return descr, dserr.Err(err)
    }
    // Return saved descr
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


func (store *Store) DeleteFile(login string, filePath string) (*dsdescr.File, error) {
    var err error
    var fileDescr *dsdescr.File
    filePath = cleanPath(filePath)
    has, err := store.reg.HasFile(login, filePath)
    if err != nil {
        return fileDescr, dserr.Err(err)
    }
    if !has {
        return fileDescr, dserr.Err(err)
    }
    fileDescr, err = store.reg.GetFile(login, filePath)
    if err != nil {
        return fileDescr, dserr.Err(err)
    }
    err = store.deleteFile(fileDescr)
    if err != nil {
        return fileDescr, dserr.Err(err)
    }
    return fileDescr, dserr.Err(err)
}

func (store *Store) FileStats(login, pattern, regular, gPattern string, reader io.Reader) (int64, int64, error) {
    var err error
    var usage int64
    var count int64
    cb := func(descr *dsdescr.File) error {
        var err error
        return err
    }

    descrs, err := store.loopFiles(login, pattern, regular, gPattern, cb, reader)
    if err != nil {
        return count, usage, err
    }
    for _, descr := range descrs {
        usage += descr.DataSize
        count++
    }
    return count, usage, err
}

type loopFunc = func(descr *dsdescr.File) error

func (store *Store) EraseFiles(login, pattern, regular, gPattern string, erase bool, reader io.Reader) ([]*dsdescr.File, error) {
    cb := func(descr *dsdescr.File) error {
        var err error
        if erase {
            store.deleteFile(descr)
            dslog.LogDebugf("delete file %s", descr.FilePath)
        }
        return err
    }
    return store.loopFiles(login, pattern, regular, gPattern, cb, reader)
}

func (store *Store) ListFiles(login, pattern, regular, gPattern string, reader io.Reader) ([]*dsdescr.File, error) {
    cb := func(descr *dsdescr.File) error {
        var err error
        return err
    }
    return store.loopFiles(login, pattern, regular, gPattern, cb, reader)
}

func (store *Store) loopFiles(login, pattern, regular, gPattern string, callback loopFunc, reader io.Reader) ([]*dsdescr.File, error) {
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

    errChan := make(chan error, 16)
    checker := func() {
        buf := make([]byte, 1)
        for {
            _, err := reader.Read(buf)
            if err != nil {
                errChan <- err
                return
            }
        }
    }
    go checker()

    for _, descr := range descrs {

        select {
            case err := <-errChan:
                err = fmt.Errorf("connection error: %s", err)
                return resDescrs, dserr.Err(err)
            default:
        }

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
            callback(descr)
            resDescrs = append(resDescrs, descr)
        }
    }
    return resDescrs, dserr.Err(err)
}



func (store *Store) deleteFile(fileDescr *dsdescr.File) error {
    var err error

    dslog.LogDebugf("erase #2 file %s", fileDescr.FilePath)


    blockDescrs, err := store.reg.ListBlocks(fileDescr.FileId)
    if err != nil {
        return dserr.Err(err)
    }
    cleanBlocks := true
    for _, descr := range blockDescrs {
        block, err := fsfile.OpenBlock(store.dataDir, descr)
        if block == nil && err != nil {
            cleanBlocks = false
            continue
        }
        err = block.Clean()
        if err != nil {
            cleanBlocks = false
            continue
        }
        err = store.reg.DeleteBlock(descr.FileId, descr.BatchId, descr.BlockType, descr.BlockId)
        if err != nil {
            cleanBlocks = false
            continue
        }
    }
    cleanBatchs := true
    batchDescrs, err := store.reg.ListBatchs(fileDescr.FileId)
    if err != nil {
        return dserr.Err(err)
    }
    for _, descr := range batchDescrs {
        err = store.reg.DeleteBatch(descr.FileId, descr.BatchId)
        if err != nil {
            cleanBatchs = false
            continue
        }

    }
    switch {
        case cleanBatchs && cleanBlocks:
            dslog.LogDebugf("delete file %s", fileDescr.FilePath)

            err = store.reg.DeleteFile(fileDescr.Login, fileDescr.FilePath)
            if err != nil {
                err = fmt.Errorf("cannot delete file descr for %s, err: %v", fileDescr.FilePath, err)
                return dserr.Err(err)
            }
            store.fileAlloc.FreeId(fileDescr.FileId)
            if err != nil {
                return dserr.Err(err)
            }
        default:
            dslog.LogDebugf("trash file %s", fileDescr.FilePath)

            trashDescr := dsdescr.NewFile()
            *trashDescr = *fileDescr

            randBin := make([]byte, 16)
            rand.Read(randBin)
            randStr := hex.EncodeToString(randBin)
            trashPath := filepath.Join("/.trash/", randStr, fileDescr.FilePath)

            trashDescr.FilePath = trashPath
            err = store.reg.PutFile(trashDescr)
            if err != nil {
                return dserr.Err(err)
            }
            err = store.reg.DeleteFile(fileDescr.Login, fileDescr.FilePath)
            if err != nil {
                return dserr.Err(err)
            }
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
