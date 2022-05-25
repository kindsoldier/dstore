/*
 * Copyright 2022 Oleg Borodin  <borodin@unix7.org>
 */

package bsrec

import (
    "fmt"
    "io/fs"
    "io"
    "path/filepath"
    "os"
    "errors"

    "ndstore/dsrpc"
    "ndstore/dscom"
    "ndstore/xtools"
    "ndstore/bstore/bssrv/bsreg"
)

const blockFileExt string = ".blk"
const storeDBName  string = "reg.db"

type Block = dscom.BlockMI


type Store struct {
    dataRoot string
    reg *bsreg.Reg
    dirPerm   fs.FileMode
    filePerm  fs.FileMode
}

func NewStore(dataRoot string) *Store {
    var store Store
    store.dataRoot  = dataRoot
    store.dirPerm   = 0755
    store.filePerm  = 0644
    return &store
}

func (store *Store) SetPerm(dirPerm, filePerm fs.FileMode) {
    store.dirPerm = dirPerm
    store.filePerm = filePerm
}

func (store *Store) OpenReg() error {
    var err error
    reg := bsreg.NewReg()
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

func (store *Store) SaveBlock(fileId, batchId, blockId int64, blockReader io.Reader,
                                                                        blockSize int64) error {
    var err error

    blockExists, err := store.reg.BlockExists(fileId, batchId, blockId)
    if err != nil {
        return err
    }
    if blockExists {
        return errors.New("block yet exists")
    }

    fileName := MakeBlockName(fileId, batchId, blockId)
    subdirName := MakeDirName(fileName)
    dirPath := filepath.Join(store.dataRoot, subdirName)
    os.MkdirAll(dirPath, store.dirPerm)

    fullFilePath := filepath.Join(dirPath, fileName)

    blockFile, err := os.OpenFile(fullFilePath, os.O_RDWR|os.O_CREATE|os.O_TRUNC, store.filePerm)
    defer blockFile.Close()
    if err != nil {
        return err
    }
    _, err = dsrpc.CopyBytes(blockReader, blockFile, blockSize)
    if err != nil {
        return err
    }

    filePath := filepath.Join(subdirName, fileName)
    err = store.reg.AddBlock(fileId, batchId, blockId, blockSize, filePath)
    if err != nil {
        os.Remove(filePath)
        return err
    }
    return err
}

func (store *Store) BlockExists(fileId, batchId, blockId int64) (int64, error) {
    var err error
    var filePath string
    filePath, blockSize, err := store.reg.GetBlock(fileId, batchId, blockId)
    if err != nil {
        return blockSize, err
    }
    filePath = filepath.Join(store.dataRoot, filePath)
    blockFile, err := os.OpenFile(filePath, os.O_RDONLY, 0)
    defer blockFile.Close()
    if err != nil {
        return blockSize, err
    }
    return blockSize, err
}


func (store *Store) LoadBlock(fileId, batchId, blockId int64,
                                                    blockWriter io.Writer) error {
    var err error
    var filePath string
    filePath, blockSize, err := store.reg.GetBlock(fileId, batchId, blockId)
    if err != nil {
        return err
    }
    filePath = filepath.Join(store.dataRoot, filePath)
    blockFile, err := os.OpenFile(filePath, os.O_RDONLY, 0)
    defer blockFile.Close()
    if err != nil {
        return err
    }
    _, err = dsrpc.CopyBytes(blockFile, blockWriter, blockSize)
    if err != nil {
        return err
    }
    return err
}

func (store *Store) DeleteBlock(fileId, batchId, blockId int64) error {
    var err error
    var filePath string
    filePath, _, err = store.reg.GetBlock(fileId, batchId, blockId)
    if err != nil {
        return err
    }
    filePath = filepath.Join(store.dataRoot, filePath)
    err = os.Remove(filePath)
    if err != nil {
        return err
    }
    err = store.reg.DeleteBlock(fileId, batchId, blockId)
    if err != nil {
        return err
    }
    return err
}

func (store *Store) ListBlocks() ([]Block, error) {
    var err error
    blocks, err := store.reg.ListBlocks()
    if err != nil {
        return blocks, err
    }
    return blocks, err
}

func MakeBlockName(fileId, batchId, blockId int64) string {
    var fileName string
    fileName = fmt.Sprintf("%020d-%020d-%020d", fileId, batchId, blockId)
    fileName = xtools.Raw2HashString([]byte(fileName))
    fileName = fileName + blockFileExt
    return fileName
}

func MakeDirName(fileName string) string {
    hash := xtools.Raw2HashBytes([]byte(fileName))
    l1 := string(hash[0:1])
    l2 := string(hash[2:3])
    dirName := filepath.Join(l1, l2)
    return dirName
}
