/*
 * Copyright 2022 Oleg Borodin  <borodin@unix7.org>
 */

package fstore

import (
    "fmt"
    "io"
    "path/filepath"
    "time"

    "dstore/fstore/fssrv/fsfile"
    "dstore/dsdescr"
    "dstore/dserr"
    //"dstore/dslog"
)

func (store *Store) SaveFile(login string, filePath string, fileReader io.Reader, fileSize int64) error {
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

    filePath = filepath.Clean(filePath)

    has, err = store.reg.HasEntry(login, filePath)
    if err != nil {
        return dserr.Err(err)
    }
    if has {
        err = fmt.Errorf("entry %s exist", filePath)
        return dserr.Err(err)
    }

    var batchSize   int64 = 5
    var blockSize   int64 = 8 * 1024 * 1024

    if fileSize < blockSize * batchSize {
        blockSize = fileSize / batchSize
        rs := int64(1024 * 10)
        bs := blockSize / rs
        blockSize = (bs + 1) * rs
    }

    var fileId int64 = 1

    file, err := fsfile.NewFile(store.dataDir, store.reg, fileId, batchSize, blockSize)
    if err != nil {
        return dserr.Err(err)
    }

    entry := dsdescr.NewEntry()
    entry.FilePath  = filePath
    entry.FileId    = fileId
    entry.CreatedAt = time.Now().Unix()
    entry.UpdatedAt = entry.CreatedAt
    err = store.reg.PutEntry(login, entry)
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
