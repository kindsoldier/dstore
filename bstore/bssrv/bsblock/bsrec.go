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
    "crypto/sha1"
    "math/rand"
    "time"

    "github.com/minio/highwayhash"

    "ndstore/bstore/bssrv/bsbreg"
    "ndstore/dsrpc"
    "ndstore/dscom"
    "ndstore/dserr"
    "ndstore/dslog"
)

const blockFileExt string = ".blk"

type Store struct {
    dataRoot string
    reg *bsbreg.Reg
    dirPerm   fs.FileMode
    filePerm  fs.FileMode
    wasteChan chan byte
}

func NewStore(dataRoot string, reg *bsbreg.Reg) *Store {
    var store Store
    store.dataRoot  = dataRoot
    store.reg       = reg
    store.dirPerm   = 0755
    store.filePerm  = 0644
    store.wasteChan = make(chan byte, 1024)
    return &store
}

func (store *Store) SetDirPerm(dirPerm fs.FileMode) {
    store.dirPerm = dirPerm
}

func (store *Store) SetFilePerm(filePerm fs.FileMode) {
    store.filePerm = filePerm
}

func (store *Store) SaveBlock(fileId, batchId, blockId int64, blockType string, blockSize, dataSize int64,
                            hashAlg, hashInit, hashSum string, blockReader io.Reader, binSize int64) error {

    var err error
    const uCounter int64 = 1

    exists, used, _, _, err := store.getBlockParams(fileId, batchId, blockId, blockType)

    if err != nil {
        return dserr.Err(err)
    }
    if exists && used {
        err = errors.New("block yet exists")
        return dserr.Err(err)
    }

    if exists && !used {
        err = store.dropBlock(fileId, batchId, blockId, blockType)
        if err != nil {
            return dserr.Err(err)
        }
    }

    err = validateFileId(fileId)
    if err != nil {
        return dserr.Err(err)
    }

    err = validateBatchId(batchId)
    if err != nil {
        return dserr.Err(err)
    }

    err = validateBlockId(blockId)
    if err != nil {
        return dserr.Err(err)
    }

    filePath := makeFilePath(fileId, batchId, blockId, blockType)
    fullFilePath := filepath.Join(store.dataRoot, filePath)
    os.MkdirAll(filepath.Dir(fullFilePath), store.dirPerm)

    blockFile, err := os.OpenFile(fullFilePath, os.O_RDWR|os.O_CREATE|os.O_TRUNC, store.filePerm)
    defer blockFile.Close()
    if err != nil {
        return dserr.Err(err)
    }

    switch {
        case len(hashSum) == 0:
            hashIBytes := make([]byte, 32)
            rand.Read(hashIBytes)
            hasher, _ := highwayhash.New(hashIBytes)

            multiWriter := io.MultiWriter(blockFile, hasher)
            _, err = dsrpc.CopyBytes(blockReader, multiWriter, binSize)
            if err != nil {
                return dserr.Err(err)
            }
            hashBytes := hasher.Sum(nil)
            hashSum = hex.EncodeToString(hashBytes)
            hashInit = hex.EncodeToString(hashIBytes)
            hashAlg = dscom.HashTypeHW
        default:
           _, err = dsrpc.CopyBytes(blockReader, blockFile, binSize)
            if err != nil {
                return dserr.Err(err)
            }
    }
    descr := dscom.NewBlockDescr()
    descr.FileId    = fileId
    descr.BatchId   = batchId
    descr.BlockId   = blockId
    descr.UCounter  = uCounter
    descr.BlockSize = blockSize
    descr.DataSize  = dataSize
    descr.FilePath  = filePath
    descr.BlockType = blockType
    descr.HashAlg   = hashAlg
    descr.HashInit  = hashInit
    descr.HashSum   = hashSum
    err = store.reg.AddBlockDescr(descr)

    if err != nil {
        os.Remove(filePath)
        return dserr.Err(err)
    }
    return dserr.Err(err)
}

func (store *Store) BlockParams(fileId, batchId, blockId int64, blockType string) (bool, int64, error) {
    var err error
    var filePath string
    var exists bool
    exists, used, filePath, blockSize, err := store.getBlockParams(fileId, batchId, blockId, blockType)
    if err != nil {
        return exists, blockSize, dserr.Err(err)
    }
    if !exists || (exists && !used) {
        exists = false
        return exists, blockSize, dserr.Err(err)
    }

    exists = false
    filePath = filepath.Join(store.dataRoot, filePath)
    blockFile, err := os.OpenFile(filePath, os.O_RDONLY, 0)
    defer blockFile.Close()
    if err != nil {
        return exists, blockSize, dserr.Err(err)
    }
    exists = true
    return exists, blockSize, dserr.Err(err)
}

func (store *Store) CheckBlock(fileId, batchId, blockId int64, blockType string) (bool, error) {
    var err error
    var filePath string
    var correct bool
    exists, used, filePath, dataSize, err := store.getBlockParams(fileId, batchId, blockId, blockType)
    if err != nil {
        return correct, dserr.Err(err)
    }

    if !exists || (exists && !used) {
        return correct, dserr.Err(err)
    }

    filePath = filepath.Join(store.dataRoot, filePath)
    blockFile, err := os.OpenFile(filePath, os.O_RDONLY, 0)
    defer blockFile.Close()
    if err != nil {
        return correct, dserr.Err(err)
    }

    fileInfo, err := blockFile.Stat()
    if err != nil {
        return correct, dserr.Err(err)
    }
    fileSize := fileInfo.Size()

    if fileSize != dataSize {
        return correct, fmt.Errorf("data size and file size mismatch: %d %d", dataSize, fileSize)
    }

    correct = true
    return correct, dserr.Err(err)
}

func (store *Store) LoadBlock(fileId, batchId, blockId int64, blockType string, blockWriter io.Writer) error {
    var err error
    var filePath string
    exists, used, filePath, blockSize, err := store.getBlockParams(fileId, batchId, blockId, blockType)
    if err != nil {
        return dserr.Err(err)
    }
    if !exists || !used {
        err = errors.New("block not exists")
        return dserr.Err(err)
    }

    err = store.reg.IncBlockDescrUC(fileId, batchId, blockId, blockType)
    if err != nil {
        return dserr.Err(err)
    }
    defer store.reg.DecBlockDescrUC(fileId, batchId, blockId, blockType)

    filePath = filepath.Join(store.dataRoot, filePath)
    blockFile, err := os.OpenFile(filePath, os.O_RDONLY, 0)
    defer blockFile.Close()
    if err != nil {
        return dserr.Err(err)
    }
    _, err = dsrpc.CopyBytes(blockFile, blockWriter, blockSize)
    if err != nil {
        return dserr.Err(err)
    }
    return dserr.Err(err)
}


func (store *Store) DeleteBlock(fileId, batchId, blockId int64, blockType string) error {
    var err error
    err = store.reg.DecBlockDescrUC(fileId, batchId, blockId, blockType)
    defer store.pushWC()
    if err != nil {
        return dserr.Err(err)
    }
    return dserr.Err(err)
}


func (store *Store) PurgeAll() error {
    var err error
    var filePath string
    blocks, err := store.reg.ListBlockDescrs()
    if err != nil {
        return dserr.Err(err)
    }
    for _, block := range blocks {
        filePath = filepath.Join(store.dataRoot, block.FilePath)
        _ = os.Remove(filePath)
        //if err != nil {
        //    return dserr.Err(err)
        //}
        err = store.reg.EraseBlockDescr(block.FileId, block.BatchId, block.BlockId, block.BlockType)
        if err != nil {
            return dserr.Err(err)
        }
    }
    //err = store.reg.EraseAll()
    //if err != nil {
    //    return dserr.Err(err)
    //}
    return dserr.Err(err)
}

func (store *Store) ListBlocks() ([]*dscom.BlockDescr, error) {
    var err error
    blocks, err := store.reg.ListBlockDescrs()
    if err != nil {
        return blocks, dserr.Err(err)
    }
    return blocks, dserr.Err(err)
}

func (store *Store) WasteCollector() {
    for {
        exists, bl, err := store.reg.GetUnusedBlockDescr()
        if exists && err == nil {
            err = store.dropBlock(bl.FileId, bl.BatchId, bl.BlockId, bl.BlockType)
            if err != nil {
                dslog.LogDebug("delete waste block err:", dserr.Err(err))
            }
            continue
        }
        select {
            case <- store.wasteChan:
            case <-time.After(time.Second * 10):
        }
    }
}

func (store *Store) pushWC() {
    if cap(store.wasteChan) - len(store.wasteChan) > 1 {
        store.wasteChan <- 0xff
    }
}

func (store *Store) dropBlock(fileId, batchId, blockId  int64, blockType string) error {
    var err error
    var filePath string
    exists, used, filePath, _, _ := store.getBlockParams(fileId, batchId, blockId, blockType)
    if err != nil {
        return dserr.Err(err)
    }
    if exists && used {
        return dserr.Err(err)
    }
    if !exists  {
        return dserr.Err(err)
    }

    if len(filePath) > 0 {
        filePath = filepath.Join(store.dataRoot, filePath)
        err = os.Remove(filePath)
        if err != nil {
            return dserr.Err(err)
        }
    }
    err = store.reg.EraseBlockDescr(fileId, batchId, blockId, blockType)
    if err != nil {
        return dserr.Err(err)
    }
    return dserr.Err(err)
}

func (store *Store) getBlockParams(fileId, batchId, blockId int64, blockType string) (bool, bool, string, int64, error) {
    var err error
    var exists bool
    var used bool
    var filePath string
    var dataSize int64
    exists, descr, err := store.reg.GetBlockDescr(fileId, batchId, blockId, blockType)
    if err != nil {
         return exists, used, filePath, dataSize, dserr.Err(err)
    }
    if !exists {
         return exists, used, filePath, dataSize, dserr.Err(err)
    }
    filePath = descr.FilePath
    dataSize = descr.DataSize
    if descr.UCounter > 0 {
        used = true
    }
     return exists, used, filePath, dataSize, dserr.Err(err)
}

func makeFilePath(fileId, batchId, blockId int64, blockType string) string {
    origin := fmt.Sprintf("%020d-%020d-%020d-%020d-%s", fileId, batchId, blockId, blockType)
    hasher := sha1.New()
    hasher.Write([]byte(origin))
    hashSum := hasher.Sum(nil)
    hashHex := hex.EncodeToString(hashSum)
    fileName := hashHex + blockFileExt
    l1 := string(hashHex[0:1])
    l2 := string(hashHex[2:3])
    dirName := filepath.Join(l1, l2)
    return filepath.Join(dirName, fileName)
}

func validateBlockId(id int64) error {
    var err error
    if id < 0 {
        err = errors.New("block id must be equal or greater than 0")
    }
    return dserr.Err(err)
}

func validateFileId(id int64) error {
    var err error
    if id < 0 {
        err = errors.New("file id must be equal or greater than 0")
    }
    return dserr.Err(err)
}

func validateBatchId(id int64) error {
    var err error
    if id < 0 {
        err = errors.New("batch id must be equal or greater than 0")
    }
    return dserr.Err(err)
}
