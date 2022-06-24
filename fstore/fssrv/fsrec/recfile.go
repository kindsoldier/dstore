/*
 * Copyright 2022 Oleg Borodin  <borodin@unix7.org>
 */

package fsrec

import (
    "fmt"
    "io"
    "path/filepath"
    "time"

    "ndstore/fstore/fssrv/fsfile"
    "ndstore/dscom"
    "ndstore/dserr"
    "ndstore/dslog"
)

const blockFileExt string = ".blk"

func (store *Store) SaveFile(userName string, filePath string, fileReader io.Reader, fileSize int64) error {
    var err error

    exists, userDescr, err := store.reg.GetUserDescr(userName)
    if !exists {
        err = fmt.Errorf("user %s not exist", userName)
        return dserr.Err(err)
    }
    if err != nil {
        return dserr.Err(err)
    }
    userId := userDescr.UserId
    dirPath, fileName := pathSplit(filePath)
    filePath = filepath.Join(dirPath, fileName)

    exists, _, err = store.reg.GetEntryDescr(userId, dirPath, fileName)
    if exists {
        err = fmt.Errorf("file entry %s exist", filePath)
        return dserr.Err(err)
    }

    const batchSize   int64 = 5
    const blockSize   int64 = 1024 * 1024

    fileId, file, err := fsfile.NewFile(store.reg, store.dataRoot, batchSize, blockSize)
    defer file.Close()

    // todo: dec file usage if exit with error
    if err != nil {
        return dserr.Err(err)
    }
    _, err = file.Write(fileReader, fileSize)
    if err != nil && err != io.EOF {
        return dserr.Err(err)
    }


    err = store.reg.AddEntryDescr(userId, dirPath, fileName, fileId)
    if err != nil {
        return dserr.Err(err)
    }

    err = file.Close()
    if err != nil {
        return dserr.Err(err)
    }
    return dserr.Err(err)
}

func (store *Store) FileExists(userName string, filePath string) (bool, int64, error) {
    var err error
    var fileSize int64
    var exists bool

    userId, err := store.getUserId(userName)
    if err != nil {
        return exists, fileSize, dserr.Err(err)
    }

    dirPath, fileName := pathSplit(filePath)

    exists, entry, err := store.reg.GetEntryDescr(userId, dirPath, fileName)
    if err != nil {
        return exists, fileSize, dserr.Err(err)
    }
    if !exists {
        filePath := filepath.Join(dirPath, fileName)
        err = fmt.Errorf("file entry for %s not exist", filePath)
        return exists, fileSize, dserr.Err(err)
    }

    exists, fileDescr, err := store.reg.GetFileDescr(entry.FileId)
    if err != nil {
        return exists, fileSize, dserr.Err(err)
    }
    if !exists {
        filePath := filepath.Join(dirPath, fileName)
        err = fmt.Errorf("file desciptor for file %s not found", filePath)
        return exists, fileSize, dserr.Err(err)
    }

    fileSize = fileDescr.FileSize

    return exists, fileSize, dserr.Err(err)
}

func (store *Store) LoadFile(userName string, filePath string, fileWriter io.Writer) error {
    var err error

    userId, err := store.getUserId(userName)
    if err != nil {
        return dserr.Err(err)
    }
    dirPath, fileName := pathSplit(filePath)
    exists, entry, err := store.reg.GetEntryDescr(userId, dirPath, fileName)
    if err != nil {
        return dserr.Err(err)
    }
    if !exists {
        filePath := filepath.Join(dirPath, fileName)
        err = fmt.Errorf("file entry for %s not found", filePath)
        return dserr.Err(err)
    }
    file, err := fsfile.OpenFile(store.reg, store.dataRoot, entry.FileId)
    defer file.Close()
    if err != nil {
        return dserr.Err(err)
    }
    _, err = file.Read(fileWriter)
    if err != nil {
        return dserr.Err(err)
    }
    return dserr.Err(err)
}

func (store *Store) getEntryFileId(userId int64, dirPath, fileName string) (int64, error) {
    var err error
    var fileId int64
    exists, entryDescr, err := store.reg.GetEntryDescr(userId, dirPath, fileName)
    if err != nil {
        return fileId, dserr.Err(err)
    }
    filePath := filepath.Join(dirPath, fileName)
    if !exists {
        err = fmt.Errorf("file %s not exists", filePath)
        return fileId, dserr.Err(err)
    }
    fileId = entryDescr.FileId
    return fileId, dserr.Err(err)
}


func (store *Store) DeleteFile(userName string, filePath string) error {
    var err error

    userId, err := store.getUserId(userName)
    if err != nil {
        return dserr.Err(err)
    }
    dirPath, fileName := pathSplit(filePath)

    exists, entry, err := store.reg.GetEntryDescr(userId, dirPath, fileName)
    if err != nil {
        return dserr.Err(err)
    }
    if !exists {
        filePath := filepath.Join(dirPath, fileName)
        err = fmt.Errorf("file %s not exist", filePath)
        return dserr.Err(err)
    }
    fileId := entry.FileId

    err = store.reg.DeleteEntryDescr(userId, dirPath, fileName)
    if err != nil {
        return dserr.Err(err)
    }
    err = store.reg.DecFileDescrUC(fileId)
    if err != nil {
        return dserr.Err(err)
    }
    store.pushWC()
    return dserr.Err(err)
}

func (store *Store) WasteCollector() {
    for {
        exists, fd, err := store.reg.GetUnusedFileDescr()
        if exists && err == nil {
            dslog.LogDebug("delete waste file :", fd.FileId)
            err = store.eraseFile(fd.FileId)
            if err != nil {
                dslog.LogDebug("delete file err:", dserr.Err(err))
            }
            continue
        }
        select {
            case <-store.wasteChan:
            case <-time.After(time.Second * 1):
        }
    }
}

func (store *Store) LostCollector() {
    return
    for {
        exists, fd, err := store.reg.GetLostedFileDescr()
        if exists && err == nil {
            dslog.LogDebug("delete lost file :", fd.FileId)
            err = store.eraseFile(fd.FileId)
            if err != nil {
                dslog.LogDebug("delete file err:", dserr.Err(err))
            }
            continue
        }
        select {
            case <-time.After(time.Second * 1):
                continue
        }
    }
}

func (store *Store) pushWC() {
    if cap(store.wasteChan) - len(store.wasteChan) > 1 {
        store.wasteChan <- 0xff
    }
}

func (store *Store) eraseFile(fileId int64) error {
    var err error
    file, err := fsfile.OpenFile(store.reg, store.dataRoot, fileId)
    if err == nil {
        err := file.Erase()
        if err != nil {
            file.Close()
            file, _ = fsfile.OpenFile(store.reg, store.dataRoot, fileId)
            file.BrutalErase()
            file.Close()
        }
    }
    return dserr.Err(err)
}

func (store *Store) ListFiles(userName string, dirPath string) ([]*dscom.EntryDescr, error) {
    var err error
    entries := make([]*dscom.EntryDescr, 0)

    dirPath = dirConv(dirPath)

    userId, err := store.getUserId(userName)
    if err != nil {
        return entries, dserr.Err(err)
    }
    entries, err = store.reg.ListEntryDescr(userId, dirPath)
    if err != nil {
        return entries, dserr.Err(err)
    }
    return entries, dserr.Err(err)
}

func pathSplit(filePath string) (string, string) {
    filePath = "/" + filePath
    filePath = filepath.Clean(filePath)
    dirPath, fileName := filepath.Split(filePath)
    dirPath = dirPath + "/"
    dirPath = filepath.Clean(dirPath)
    if dirPath == "" {
        dirPath = "/"
    }
    return dirPath, fileName
}

func dirConv(dirPath string) string {
    dirPath = "/" + dirPath + "/"
    dirPath = filepath.Clean(dirPath)
    return dirPath
}
