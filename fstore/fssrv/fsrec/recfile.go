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
    if err != nil {
        file.Delete()
        return dserr.Err(err)
    }
    written, err := file.Write(fileReader, fileSize)
    if err == io.EOF {
        dslog.LogDebugf("file %d write error is %v", fileId, err)
        file.Delete()
        return dserr.Err(err)
    }
    if err != nil  {
        file.Delete()
        return dserr.Err(err)
    }
    if written != fileSize {
    //    file.Delete()
    //    err = fmt.Errorf("file %d size mismatch, file size %d, written %d ", fileId, fileSize, written)
    //    return dserr.Err(err)
    }

    err = store.reg.AddEntryDescr(userId, dirPath, fileName, fileId)
    if err != nil {
        file.Delete()
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

    exists, fileDescr, err := store.reg.GetNewestFileDescr(entry.FileId)
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

    err = store.reg.EraseEntryDescr(userId, dirPath, fileName)
    if err != nil {
        return dserr.Err(err)
    }
    file, err := fsfile.ForceOpenFile(store.reg, store.dataRoot, fileId)
    if err != nil {
        return dserr.Err(err)
    }
    file.Delete()
    if err != nil {
        return dserr.Err(err)
    }
    go store.pushFileWC()
    return dserr.Err(err)
}

func (store *Store) pushFileWC() {
    time.Sleep(1 * time.Second) // todo: how much size of timeout?
    if cap(store.fileWCChan) - len(store.fileWCChan) > 1 {
        store.fileWCChan <- 0xff
    }
}


func (store *Store) StoredFileDistributing() {
    for {
        select {
            case <-store.fileWCChan:
            case <-time.After(time.Second * 3):
        }

        exists, descr, err := store.reg.GetAnyNotDistrFileDescr()
        //dslog.LogDebug("file disributor call", exists, err)
        if exists && err == nil {
            dslog.LogDebug("distrubute file:", descr.FileId, descr.FileVer, descr.IsDistr)
            file, err := fsfile.OpenFile(store.reg, store.dataRoot, descr.FileId)
            if err != nil {
                dslog.LogDebug("distribute file open err:", dserr.Err(err))
                file.Close()
                continue
            }

            distr := NewFileDistr(store.dataRoot, store.reg)
            distr.LoadPool()
            err = file.Distribute(distr)
            if err != nil {
                dslog.LogDebug("distribute file err:", dserr.Err(err))
            }
            file.Close()
            //continue
        }
        select {
            case <-store.fileWCChan:
            case <-time.After(time.Second * 3):
        }
    }
}

func (store *Store) WasteFileCollecting() {
    for {
        //dslog.LogDebug("file waste collecr call")
        exists, descr, err := store.reg.GetAnyUnusedFileDescr()
        if exists && err == nil {
            dslog.LogDebug("delete waste file:", descr.FileId)
            file, err := fsfile.OpenSpecUnusedFile(store.reg, store.dataRoot, descr.FileId, descr.FileVer)
            err = file.Erase()
            if err != nil {
                dslog.LogDebug("delete file err:", dserr.Err(err))
            }
            file.Close()
            continue
        }
        select {
            case <-store.fileWCChan:
            case <-time.After(time.Second * 3):
        }
    }
}

func (store *Store) WasteBatchCollecting() {
    for {
        //dslog.LogDebug("batch waste collecr call")
        exists, descr, err := store.reg.GetAnyUnusedBatchDescr()
        if exists && err == nil {
            //dslog.LogDebug("delete waste batch:", descr.FileId, descr.BatchId)
            batch, err := fsfile.OpenSpecUnusedBatch(store.reg, store.dataRoot, descr.FileId,
                                                                descr.BatchId, descr.BatchVer)
            err = batch.Erase()
            if err != nil {
                dslog.LogDebug("delete batch err:", dserr.Err(err))
            }
            batch.Close()
            continue
        }
        select {
            case <-store.batchWCChan:
            case <-time.After(time.Second * 3):
        }
    }
}

func (store *Store) WasteBlockCollecting() {
    for {
        //dslog.LogDebug("block waste collecr call")
        exists, descr, err := store.reg.GetAnyUnusedBlockDescr()
        if exists && err == nil {
            //dslog.LogDebug("delete waste block:", descr.FileId, descr.BatchId, descr.BlockId)
            block, err := fsfile.OpenSpecUnusedBlock(store.reg, store.dataRoot, descr.FileId, descr.BatchId, descr.BlockId,
                                                        descr.BlockType, descr.BlockVer)
            err = block.Erase()
            if err != nil {
                dslog.LogDebug("delete batch err:", dserr.Err(err))
            }
            block.Close()
            continue
        }
        select {
            case <-store.blockWCChan:
            case <-time.After(time.Second * 3):
        }
    }
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
