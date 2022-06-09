/*
 * Copyright 2022 Oleg Borodin  <borodin@unix7.org>
 */

package bsblock

import (
    "fmt"
    "io/fs"
    "io"
    "path/filepath"
    "os"
    "errors"
    "encoding/hex"
    "crypto/md5"
    "time"

    "ndstore/bstore/bssrv/bsbreg"
    "ndstore/dsrpc"
    "ndstore/dscom"
)

const blockFileExt string = ".blk"

type Store struct {
    dataRoot string
    reg *bsbreg.Reg
    dirPerm   fs.FileMode
    filePerm  fs.FileMode
}

func NewStore(dataRoot string, reg *bsbreg.Reg) *Store {
    var store Store
    store.dataRoot  = dataRoot
    store.reg       = reg
    store.dirPerm   = 0755
    store.filePerm  = 0644
    return &store
}

func (store *Store) SetDirPerm(dirPerm fs.FileMode) {
    store.dirPerm = dirPerm
}

func (store *Store) SetFilePerm(filePerm fs.FileMode) {
    store.filePerm = filePerm
}

func (store *Store) SaveBlock(fileId, batchId, blockId, blockSize, dataSize int64, blockReader io.Reader,
                                    binSize int64, blockType, hashAlg, hashInit, hashSum string) error {
    var err error

    blockExists, err := store.reg.BlockDescrExists(fileId, batchId, blockId, blockType)
    if err != nil {
        return err
    }
    if blockExists {
        return errors.New("block yet exists")
    }

    err = validateFileId(fileId)
    if err != nil {
        return err
    }

    err = validateBatchId(batchId)
    if err != nil {
        return err
    }

    err = validateBlockId(blockId)
    if err != nil {
        return err
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
    _, err = dsrpc.CopyBytes(blockReader, blockFile, binSize)
    if err != nil {
        return err
    }
    filePath := filepath.Join(subdirName, fileName)
    err = store.reg.AddBlockDescr(fileId, batchId, blockId, blockSize, dataSize, filePath,
                                                            blockType, hashAlg, hashInit, hashSum)

    if err != nil {
        os.Remove(filePath)
        return err
    }
    return err
}

func (store *Store) BlockExists(fileId, batchId, blockId int64, blockType string) (int64, error) {
    var err error
    var filePath string
    filePath, blockSize, err := store.reg.GetBlockFilePath(fileId, batchId, blockId, blockType)
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
                                                    blockWriter io.Writer, blockType string) error {
    var err error
    var filePath string
    filePath, blockSize, err := store.reg.GetBlockFilePath(fileId, batchId, blockId, blockType)
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

func (store *Store) DeleteBlock(fileId, batchId, blockId int64, blockType string) error {
    var err error
    var filePath string
    filePath, _, err = store.reg.GetBlockFilePath(fileId, batchId, blockId, blockType)
    if err != nil {
        return err
    }
    filePath = filepath.Join(store.dataRoot, filePath)
    err = os.Remove(filePath)
    if err != nil {
        return err
    }
    err = store.reg.DeleteBlockDescr(fileId, batchId, blockId, blockType)
    if err != nil {
        return err
    }
    return err
}

func (store *Store) ListBlocks() ([]*dscom.BlockDescr, error) {
    var err error
    blocks, err := store.reg.ListBlockDescrs()
    if err != nil {
        return blocks, err
    }
    return blocks, err
}

func MakeBlockName(fileId, batchId, blockId int64) string {
    ts := time.Now().UnixNano()
    origin := fmt.Sprintf("%020d-%020d-%020d-%020d", fileId, batchId, blockId, ts)
    hasher := md5.New()
    hasher.Write([]byte(origin))
    hashSum := hasher.Sum(nil)
    hashHex := hex.EncodeToString(hashSum)
    fileName := hashHex + blockFileExt
    return fileName
}

func MakeDirName(fileName string) string {
    hasher := md5.New()
    hasher.Write([]byte(fileName))
    hashSum := hasher.Sum(nil)
    hashHex := make([]byte, hex.EncodedLen(len(hashSum)))
    hex.Encode(hashHex, hashSum)
    l1 := string(hashHex[0:1])
    l2 := string(hashHex[2:3])
    dirName := filepath.Join(l1, l2)
    return dirName
}


func validateBlockId(id int64) error {
    var err error
    if id < 1 {
        err = errors.New("block id must be greater than 0")
    }
    return err
}

func validateFileId(id int64) error {
    var err error
    if id < 1 {
        err = errors.New("file id must be greater than 0")
    }
    return err
}

func validateBatchId(id int64) error {
    var err error
    if id < 1 {
        err = errors.New("batch id must be greater than 0")
    }
    return err
}
