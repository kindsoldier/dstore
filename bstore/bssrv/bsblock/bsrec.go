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

func (store *Store) SaveBlock(descr *dscom.BlockDescr, blockReader io.Reader, binSize int64) error {

    var err error
    const uCounter int64 = 1

    dslog.LogDebugf("save block %s", getBlockDescrId(descr))


    err = validateFileId(descr.FileId)
    if err != nil {
        return dserr.Err(err)
    }
    err = validateBatchId(descr.BatchId)
    if err != nil {
        return dserr.Err(err)
    }
    err = validateBlockId(descr.BlockId)
    if err != nil {
        return dserr.Err(err)
    }

    descr.UCounter = 1
    descr.FilePath = makeFilePath()
    fullFilePath := filepath.Join(store.dataRoot, descr.FilePath)
    os.MkdirAll(filepath.Dir(fullFilePath), store.dirPerm)

    blockFile, err := os.OpenFile(fullFilePath, os.O_RDWR|os.O_CREATE|os.O_TRUNC, store.filePerm)
    defer blockFile.Close()
    if err != nil {
        os.Remove(fullFilePath)
        return dserr.Err(err)
    }

    switch {
        case len(descr.HashSum) == 0:
            hashIBytes := make([]byte, 32)
            rand.Read(hashIBytes)
            hasher, _ := highwayhash.New(hashIBytes)

            multiWriter := io.MultiWriter(blockFile, hasher)
            _, err = dsrpc.CopyBytes(blockReader, multiWriter, binSize)
            if err != nil {
                return dserr.Err(err)
            }
            hashBytes := hasher.Sum(nil)
            descr.HashSum = hex.EncodeToString(hashBytes)
            descr.HashInit = hex.EncodeToString(hashIBytes)
            descr.HashAlg = dscom.HashTypeHW
        default:

           _, err = dsrpc.CopyBytes(blockReader, blockFile, binSize)
            if err != nil {
                return dserr.Err(err)
            }
    }

    err = store.reg.AddNewBlockDescr(descr)
    if err != nil {
        os.Remove(fullFilePath)
        return dserr.Err(err)
    }
    return dserr.Err(err)
}

func (store *Store) GetBlockParams(fileId, batchId, blockId int64, blockType string, blockVer int64) (bool, int64, int64, error) {
    var err error
    var exists bool
    var dataSize int64

    dslog.LogDebugf("get block %s param", getBlockIdString(fileId, batchId, blockId,
                                                                            blockType, blockVer))

    exists, descr, err := store.reg.GetSpecBlockDescr(fileId, batchId, blockId, blockType, blockVer)
    if err != nil {
        return exists, blockVer, dataSize, dserr.Err(err)
    }
    if !exists  {
        err = fmt.Errorf("block %s not exist", getBlockIdString(fileId, batchId, blockId,
                                                                            blockType, blockVer))
        return exists, blockVer, dataSize, dserr.Err(err)
    }
    fullFilePath := filepath.Join(store.dataRoot, descr.FilePath)
    blockFile, err := os.OpenFile(fullFilePath, os.O_RDONLY, 0)
    defer blockFile.Close()
    if err != nil {
        return exists, blockVer, dataSize, dserr.Err(err)
    }
    exists = true
    blockVer = descr.BlockVer
    dataSize = descr.DataSize
    return exists, blockVer, dataSize, dserr.Err(err)
}

func (store *Store) CheckBlock(fileId, batchId, blockId int64, blockType string, blockVer int64) (bool, error) {
    var err error
    var filePath string
    correct := false
    exists, descr, err := store.reg.GetSpecBlockDescr(fileId, batchId, blockId, blockType, blockVer)
    if err != nil {
        return correct, dserr.Err(err)
    }
    if !exists {
        err = fmt.Errorf("block %s not exist", getBlockIdString(fileId, batchId, blockId,
                                                                            blockType, blockVer))
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

    if fileSize != descr.DataSize {
        return correct, fmt.Errorf("block %s data size and file size mismatch: %d %d",
                                                 getBlockDescrId(descr), descr.DataSize, fileSize)
    }

    correct = true
    return correct, dserr.Err(err)
}

func (store *Store) LoadBlock(fileId, batchId, blockId int64, blockType string, blockVer int64,
                                                                     blockWriter io.Writer) error {
    var err error

    dslog.LogDebugf("load block %s", getBlockIdString(fileId, batchId, blockId,
                                                                            blockType, blockVer))
    exists, descr, err := store.reg.GetSpecBlockDescr(fileId, batchId, blockId, blockType, blockVer)
    if err != nil {
        return dserr.Err(err)
    }
    if !exists {
        err = fmt.Errorf("block %s not exist", getBlockIdString(fileId, batchId, blockId,
                                                                            blockType, blockVer))
        return dserr.Err(err)
    }

    err = store.reg.IncSpecBlockDescrUC(descr.FileId, descr.BatchId, descr.BlockId,
                                                                 descr.BlockType, descr.BlockVer)
    if err != nil {
        return dserr.Err(err)
    }
    defer store.reg.DecSpecBlockDescrUC(descr.FileId, descr.BatchId, descr.BlockId,
                                                                 descr.BlockType, descr.BlockVer)

    filePath := filepath.Join(store.dataRoot, descr.FilePath)
    blockFile, err := os.OpenFile(filePath, os.O_RDONLY, 0)
    defer blockFile.Close()
    if err != nil {
        return dserr.Err(err)
    }
    _, err = dsrpc.CopyBytes(blockFile, blockWriter, descr.DataSize)
    if err != nil {
        return dserr.Err(err)
    }
    return dserr.Err(err)
}


func (store *Store) DeleteBlock(fileId, batchId, blockId int64, blockType string, blockVer int64) error {
    var err error

    dslog.LogDebugf("delete block %s", getBlockIdString(fileId, batchId, blockId, blockType, blockVer))

    exists, descr, err := store.reg.GetSpecBlockDescr(fileId, batchId, blockId, blockType, blockVer)
    if err != nil {
        return dserr.Err(err)
    }
    if exists {
        err = store.reg.DecSpecBlockDescrUC(descr.FileId, descr.BatchId, descr.BlockId, descr.BlockType, descr.BlockVer)
        defer store.pushWC()
        if err != nil {
            return dserr.Err(err)
        }
    }
    return dserr.Err(err)
}


func (store *Store) LinkBlock(fileId, batchId, blockId int64, blockType string, oldBlockVer, newBlockVer int64) error {
    var err error
    exists, descr, err := store.reg.GetSpecBlockDescr(fileId, batchId, blockId, blockType, oldBlockVer)
    if err != nil {
        return dserr.Err(err)
    }
    if !exists {
        err = fmt.Errorf("block %s not exist", getBlockIdString(fileId, batchId, blockId, blockType, oldBlockVer))
        return dserr.Err(err)
    }

    newFilePath := makeFilePath()
    fullNewFilePath := filepath.Join(store.dataRoot, newFilePath)
    err = os.MkdirAll(filepath.Dir(fullNewFilePath), store.dirPerm)
    if err != nil {
        return dserr.Err(err)
    }
    fullOldFilePath := filepath.Join(store.dataRoot, descr.FilePath)
    if err != nil {
        return dserr.Err(err)
    }
    err = os.Link(fullOldFilePath, fullNewFilePath)
    if err != nil {
        os.Remove(fullNewFilePath)
        return dserr.Err(err)
    }
    descr.FilePath = fullNewFilePath
    descr.BlockVer = newBlockVer

    return dserr.Err(err)
}


func (store *Store) EraseAll() error {
    var err error
    var filePath string
    descrs, err := store.reg.ListAllBlockDescrs()
    if err != nil {
        return dserr.Err(err)
    }
    for _, descr := range descrs {
        filePath = filepath.Join(store.dataRoot, descr.FilePath)
        _ = os.Remove(filePath)
        //if err != nil {
        //    return dserr.Err(err)
        //}
        err = store.reg.EraseSpecBlockDescr(descr.FileId, descr.BatchId, descr.BlockId, descr.BlockType, descr.BlockVer)
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
    descrs, err := store.reg.ListAllBlockDescrs()
    if err != nil {
        return descrs, dserr.Err(err)
    }
    return descrs, dserr.Err(err)
}

func (store *Store) WasteCollector() {
    for {
        exists, descr, err := store.reg.GetAnyUnusedBlockDescr()
        if exists && err == nil {
            dslog.LogDebugf("delete waste block %s", getBlockDescrId(descr))
            err = store.eraseBlock(descr.FileId, descr.BatchId, descr.BlockId, descr.BlockType, descr.BlockVer)
            if err != nil {
                dslog.LogErrorf("delete waste block %s err: %s", getBlockDescrId(descr), dserr.Err(err))
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

func (store *Store) eraseBlock(fileId, batchId, blockId  int64, blockType string, blockVer int64) error {
    var err error
    exists, descr, err := store.reg.GetSpecUnusedBlockDescr(fileId, batchId, blockId, blockType, blockVer)
    if err != nil {
        return dserr.Err(err)
    }
    if !exists {
        return dserr.Err(err)
    }

    var filePath string
    if len(descr.FilePath) > 0 {
        filePath = filepath.Join(store.dataRoot, descr.FilePath)
        err = os.Remove(filePath)
        if err != nil {
            return dserr.Err(err)
        }
    }
    err = store.reg.EraseSpecBlockDescr(fileId, batchId, blockId, blockType, blockVer)
    if err != nil {
        return dserr.Err(err)
    }
    return dserr.Err(err)
}


func getBlockDescrId(descr *dscom.BlockDescr) string {
    return fmt.Sprintf("%d,%d,%d,%s,%d", descr.FileId, descr.BatchId, descr.BlockId, descr.BlockType,
                                                                                    descr.BlockVer)
}

func getBlockIdString(fileId, batchId, blockId int64, blockType string, blockVer int64) string {
    return fmt.Sprintf("%d,%d,%d,%s,%d", fileId, batchId, blockId, blockType, blockVer)
}


func makeFilePath() string {
    const blockFileExt string = ".block"
    origin := make([]byte, 128)
    rand.Read(origin)
    hasher := sha1.New()
    hasher.Write(origin)
    hashSum := hasher.Sum(nil)
    hashHex := hex.EncodeToString(hashSum)
    fileName := hashHex + blockFileExt
    l1 := string(hashHex[0:1])
    l2 := string(hashHex[1:3])
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
