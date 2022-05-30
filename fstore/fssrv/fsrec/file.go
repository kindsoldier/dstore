/*
 * Copyright 2022 Oleg Borodin  <borodin@unix7.org>
 */

package fsrec

import (
    "errors"
    "io"
    "path/filepath"

    "ndstore/dscom"
    "ndstore/fstore/fssrv/fsfile"
)


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
    dirPath, fileName := pathSplit(filePath)

    err = store.reg.AddEntryDescr(dirPath, fileName, fileId)
    if err != nil {
        return err
    }
    return err
}

func (store *Store) FileExists(filePath string) (int64, error) {
    var err error
    var fileSize int64

    dirPath, fileName := pathSplit(filePath)

    entry, exists, err := store.reg.GetEntryDescr(dirPath, fileName)
    if !exists || entry == nil {
        return fileSize, errors.New("path not exists")
    }
    fileMeta, err := store.reg.GetFileDescr(entry.FileId)
    if err != nil {
        return fileSize, err
    }
    fileSize = fileMeta.FileSize

    return fileSize, err
}

func (store *Store) LoadFile(filePath string, fileWriter io.Writer) error {
    var err error

    dirPath, fileName := pathSplit(filePath)

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

    dirPath, fileName := pathSplit(filePath)

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

func (store *Store) ListFiles(dirPath string) ([]*dscom.EntryDescr, error) {
    var err error
    dirPath = dirConv(dirPath)
    entries, err := store.reg.ListEntryDescr(dirPath)
    if err != nil {
        return entries, err
    }
    return entries, err
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
