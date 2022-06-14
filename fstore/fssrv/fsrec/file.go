/*
 * Copyright 2022 Oleg Borodin  <borodin@unix7.org>
 */

package fsrec

import (
    "io"
    "path/filepath"

    "ndstore/dscom"
    "ndstore/fstore/fssrv/fsfile"
    "ndstore/dserr"
)

const blockFileExt string = ".blk"

func (store *Store) SaveFile(userName string, filePath string, fileReader io.Reader, fileSize int64) error {
    var err error

    userId, err := store.reg.GetUserId(userName)
    if err != nil {
        return dserr.Err(err)
    }

    fileId, err := store.reg.GetNewFileId()
    if err != nil {
        return dserr.Err(err)
    }

    const batchSize   int64 = 5
    const blockSize   int64 = 1 * 1024 * 1024

    file := fsfile.NewFile(store.dataRoot, fileId, batchSize, blockSize)
    err = file.Open()
    defer file.Close()
    if err != nil {
        return dserr.Err(err)
    }
    _, err = file.Lwrite(fileReader, fileSize)
    if err != nil && err != io.EOF {
        return dserr.Err(err)
    }
    meta, err := file.Meta()
    if err != nil {
        return dserr.Err(err)
    }
    err = store.reg.AddFileDescr(meta)
    if err != nil {
        return dserr.Err(err)
    }
    dirPath, fileName := pathSplit(filePath)

    err = store.reg.AddEntryDescr(userId, dirPath, fileName, fileId)
    if err != nil {
        return dserr.Err(err)
    }

    pool := NewBSPool(store.reg)
    err = pool.LoadPool()
    if err != nil {
        return dserr.Err(err)
    }

    err = file.Save(pool)
    if err != nil {
        return dserr.Err(err)
    }

    return dserr.Err(err)
}

func (store *Store) FileExists(userName string, filePath string) (int64, error) {
    var err error
    var fileSize int64

    userId, err := store.reg.GetUserId(userName)
    if err != nil {
        return fileSize, dserr.Err(err)
    }

    dirPath, fileName := pathSplit(filePath)

    entry, err := store.reg.GetEntryDescr(userId, dirPath, fileName)
    if err != nil {
        return fileSize, dserr.Err(err)
    }
    fileMeta, err := store.reg.GetFileDescr(entry.FileId)
    if err != nil {
        return fileSize, dserr.Err(err)
    }
    fileSize = fileMeta.FileSize

    return fileSize, dserr.Err(err)
}

func (store *Store) LoadFile(userName string, filePath string, fileWriter io.Writer) error {
    var err error

    userId, err := store.reg.GetUserId(userName)
    if err != nil {
        return dserr.Err(err)
    }

    dirPath, fileName := pathSplit(filePath)

    entry, err := store.reg.GetEntryDescr(userId, dirPath, fileName)
    if err != nil {
        return dserr.Err(err)
    }
    meta, err := store.reg.GetFileDescr(entry.FileId)
    if err != nil {
        return dserr.Err(err)
    }
    file := fsfile.RenewFile(store.dataRoot, meta)
    err = file.Open()
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

func (store *Store) DeleteFile(userName string, filePath string) error {
    var err error

    userId, err := store.reg.GetUserId(userName)
    if err != nil {
        return dserr.Err(err)
    }
    dirPath, fileName := pathSplit(filePath)

    entry, err := store.reg.GetEntryDescr(userId, dirPath, fileName)
    if err != nil {
        return dserr.Err(err)
    }
    fileId := entry.FileId
    meta, err := store.reg.GetFileDescr(fileId)
    if err != nil {
        return dserr.Err(err)
    }
    file := fsfile.RenewFile(store.dataRoot, meta)
    err = file.Open()
    defer file.Close()
    if err != nil {
        return dserr.Err(err)
    }
    err = file.Purge()
    if err != nil {
        return dserr.Err(err)
    }
    err = store.reg.DeleteEntryDescr(userId, dirPath, fileName)
    if err != nil {
        return dserr.Err(err)
    }
    err = store.reg.DeleteFileDescr(fileId)
    if err != nil {
        return dserr.Err(err)
    }
    return dserr.Err(err)
}

func (store *Store) ListFiles(userName string, dirPath string) ([]*dscom.EntryDescr, error) {
    var err error
    entries := make([]*dscom.EntryDescr, 0)

    dirPath = dirConv(dirPath)

    userId, err := store.reg.GetUserId(userName)
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
