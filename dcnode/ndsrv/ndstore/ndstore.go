/*
 * Copyright 2022 Oleg Borodin  <borodin@unix7.org>
 */

package ndstore

import (
    "fmt"
    "io/fs"
    "io"
    "path/filepath"
    "os"

    "dcstore/dcrpc"
    "dcstore/tools"
    "dcstore/dcnode/ndsrv/ndreg"
)

const blockFileExt string = ".blk"
const storeDBName  string = "reg.db"

const dirPerm   fs.FileMode = 0777
const filePerm  fs.FileMode = 0644

type Store struct {
    dataRoot string
    reg *ndreg.Reg
}

func NewStore(dataRoot string) *Store {
    var store Store
    store.dataRoot = dataRoot
    return &store
}

func (store *Store) OpenReg() error {
    var err error
    reg := ndreg.NewReg()
    dbPath := filepath.Join(store.dataRoot, storeDBName)
    err = reg.OpenDB(dbPath)
    if err != nil {
        return err
    }
    err = reg.MigrateDB()
    if err != nil {
        return err
    }
    store.reg = reg
    return err
}

func (store *Store) CloseReg() error {
    var err error
    if store.reg != nil {
        err = store.reg.CloseDB()
    }
    return err
}

func (store *Store) SaveBlock(clusterId, fileId, batchId, blockId int64,
                                blockReader io.Reader, blockSize int64) error {
    var err error

    fileName := MakeBlockName(clusterId, fileId, batchId, blockId)
    subdirName := MakeDirName(fileName)
    dirPath := filepath.Join(store.dataRoot, subdirName)
    os.MkdirAll(dirPath, dirPerm)

    filePath := filepath.Join(dirPath, fileName)
    blockFile, err := os.OpenFile(filePath, os.O_RDWR|os.O_CREATE, filePerm)
    defer blockFile.Close()
    if err != nil {
        return err
    }
    _, err = dcrpc.CopyBytes(blockReader, blockFile, blockSize)
    if err != nil {
        return err
    }
    err = store.reg.AddBlock(clusterId, fileId, batchId, blockId, blockSize, filePath)
    if err != nil {
        os.Remove(filePath)
        return err
    }
    return err
}

func (store *Store) LoadBlock(clusterId, fileId, batchId, blockId int64,
                                                    blockWriter io.Writer) error {
    var err error
    var filePath string
    filePath, blockSize, err := store.reg.GetBlock(clusterId, fileId, batchId, blockId)
    if err != nil {
        return err
    }
    filePath = filepath.Join(store.dataRoot, filePath)
    blockFile, err := os.OpenFile(filePath, os.O_RDONLY, 0)
    defer blockFile.Close()
    if err != nil {
        return err
    }
    _, err = dcrpc.CopyBytes(blockFile, blockWriter, blockSize)
    if err != nil {
        return err
    }
    return err
}


func (store *Store) DeleteBlock(clusterId, fileId, batchId, blockId int64) error {
    var err error
    var filePath string
    filePath, _, err = store.reg.GetBlock(clusterId, fileId, batchId, blockId)
    if err != nil {
        return err
    }
    filePath = filepath.Join(store.dataRoot, filePath)
    err = os.Remove(filePath)
    if err != nil {
        return err
    }
    err = store.reg.DeleteBlock(clusterId, fileId, batchId, blockId)
    if err != nil {
        return err
    }
    return err
}


func MakeBlockName(clusterId, fileId, batchId, blockId int64) string {
    var fileName string
    fileName = fmt.Sprintf("%020d-%020d-%020d-%020d", clusterId, fileId, batchId, blockId)
    fileName = tools.Raw2HashString([]byte(fileName))
    fileName = fileName + blockFileExt
    return fileName
}

func MakeDirName(fileName string) string {
    hash := tools.Raw2HashBytes([]byte(fileName))
    l1 := string(hash[0:1])
    l2 := string(hash[2:3])
    dirName := filepath.Join(l1, l2)
    return dirName
}
