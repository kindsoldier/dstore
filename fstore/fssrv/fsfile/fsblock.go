/*
 * Copyright 2022 Oleg Borodin  <borodin@unix7.org>
 */


package fsfile

import (
    "bytes"
    "crypto/sha1"
    "encoding/hex"
    "errors"
    "fmt"
    "hash"
    "io"
    "io/fs"
    "math/rand"
    "os"
    "path/filepath"

    "github.com/minio/highwayhash"

    "ndstore/fstore/fssrv/fsreg"
    "ndstore/dscom"
    "ndstore/dserr"
)



type Block struct {
    reg         *fsreg.Reg
    baseDir     string
    fileId      int64
    batchId     int64
    blockId     int64
    blockType   string
    blockSize   int64
    dataSize    int64
    filePath    string
    savedLoc    bool
    savedRem    bool
    locUpdated  bool
    bstoreId    int64
    fstoreId    int64

    hasher      hash.Hash
    hashAlg     string
    hashSum     []byte
    hashInit    []byte
}

const fileMode fs.FileMode  = 0644
const dirMode fs.FileMode   = 0755

func NewBlock(reg *fsreg.Reg, baseDir string, fileId, batchId, blockId int64, blockType string,
                                                                blockSize int64) (*Block, error) {
    var block Block
    var err error
    block.reg       = reg
    block.baseDir   = baseDir
    block.fileId    = fileId
    block.batchId   = batchId
    block.blockId   = blockId
    block.blockType = blockType
    block.blockSize = blockSize
    block.filePath  = makeFilePath(block.fileId, block.batchId, block.blockId, block.blockType)
    block.hashAlg   = dscom.HashTypeHW
    block.hashInit  = make([]byte, 32)
    block.hashSum   = make([]byte, 0)
    rand.Read(block.hashInit)
    block.hasher, _ = highwayhash.New(block.hashInit)

    err = block.addBlockDescr()
    if err != nil {
        return &block, dserr.Err(err)
    }
    return &block, dserr.Err(err)
}

func OpenBlock(reg *fsreg.Reg, baseDir string, fileId, batchId, blockId int64, blockType string) (*Block, error) {
    var err error
    var block Block
    exists, descr, err := reg.GetBlockDescr(fileId, batchId, blockId, blockType)
    if err != nil {
        return &block, dserr.Err(err)
    }
    if !exists {
        err = errors.New("block not exists")
        return &block, dserr.Err(err)
    }

    block.reg       = reg
    block.baseDir   = baseDir

    block.fileId    = descr.FileId
    block.batchId   = descr.BatchId
    block.blockId   = descr.BlockId
    block.blockType = descr.BlockType
    block.blockSize = descr.BlockSize

    block.dataSize  = descr.DataSize

    block.filePath  = descr.FilePath
    block.hashAlg   = descr.HashAlg
    block.hashInit, err  = hex.DecodeString(descr.HashInit)
    if err != nil {
        return &block, dserr.Err(err)
    }
    block.hashSum, err   = hex.DecodeString(descr.HashSum)
    if err != nil {
        return &block, dserr.Err(err)
    }
    block.hasher, err = highwayhash.New(block.hashInit)
    if err != nil {
        return &block, dserr.Err(err)
    }
    block.fstoreId  = descr.FStoreId
    block.bstoreId  = descr.BStoreId
    block.savedLoc  = descr.SavedLoc
    block.savedRem  = descr.SavedRem
    block.locUpdated = descr.LocUpdated

    return &block, dserr.Err(err)
}

func (block *Block) Close() error {
    var err error
    err = block.updateBlockDescr()
    if err != nil {
            return dserr.Err(err)
    }
    return dserr.Err(err)
}

func (block *Block) Write(reader io.Reader, need int64) (int64, error) {
    var err error
    var written int64
    var openMode int
    var file *os.File

    // Return if block full
    if need < 1 {
        return written, dserr.Err(err)
    }
    // Return if block full
    if block.dataSize >= block.blockSize {
        return written, dserr.Err(err)
    }
    closer := func() {
        if file != nil {
            file.Close()
        }
    }
    defer closer()
    // Prepare env
    fullFilePath := filepath.Join(block.baseDir, block.filePath)
    fullDirPath := filepath.Dir(fullFilePath)

    err = os.MkdirAll(fullDirPath, dirMode)
    if err != nil {
            return written, dserr.Err(err)
    }
    // Open local file
    openMode = os.O_APPEND|os.O_CREATE|os.O_RDWR
    file, err = os.OpenFile(fullFilePath, openMode, fileMode)
    if err != nil {
            return written, dserr.Err(err)
    }
    // Read and write remain date
    remains := block.blockSize - block.dataSize
    if remains > need {
        remains = need
    }
    written, err = copyBytes(reader, file, remains)
    block.dataSize += written
    if err == io.EOF {
        err = nil
    }
    if err != nil {
            return written, dserr.Err(err)
    }
    // Seek to begin of file
    _, err = file.Seek(0, 0)
    if err != nil {
            return written, dserr.Err(err)
    }
    // Create hashsum of data
    block.hasher.Reset()
    hWritten, err := copyBytes(file, block.hasher, block.dataSize)
    if err != nil {
            return written, dserr.Err(err)
    }
    if hWritten != block.dataSize {
        err = errors.New("incorrect block file size")
        return written, dserr.Err(err)
    }

    block.hashSum = block.hasher.Sum(nil)
    block.savedLoc   = true
    block.locUpdated = true
    // Sync block descr
    err = block.updateBlockDescr()
    if err != nil {
            return written, dserr.Err(err)
    }
    return written, dserr.Err(err)
}

func (block *Block) Read(writer io.Writer) (int64, error) {
    var err error
    var written int64
    var openMode int
    var file *os.File
    // Return if block is empty
    if block.dataSize < 1 {
        //err = errors.New("empty block")
        //return written, dserr.Err(err) //io.EOF
        return written, dserr.Err(err)
    }
    closer := func() {
        if file != nil {
            file.Close()
        }
    }
    defer closer()
    // Prepare env
    fullFilePath := filepath.Join(block.baseDir, block.filePath)
    fullDirPath := filepath.Dir(fullFilePath)
    err = os.MkdirAll(fullDirPath, dirMode)
    if err != nil {
            return written, dserr.Err(err)
    }
    // Open local file
    openMode = os.O_RDONLY
    file, err = os.OpenFile(fullFilePath, openMode, 0)
    if err != nil {
            return written, dserr.Err(err)
    }
    // Create hashsum of data
    block.hasher.Reset()
    written, err = copyBytes(file, block.hasher, block.dataSize)
    if err != nil {
        return written, dserr.Err(err)
    }
    if written != block.dataSize {
        err = errors.New("incorrect block file size")
        return written, dserr.Err(err)
    }

    hashSum := block.hasher.Sum(nil)
    if bytes.Compare(hashSum, block.hashSum) != 0 {
        err = errors.New("incorrect block hash sum")
        return written, dserr.Err(err)
    }
    // Seek to begin of file
    _, err = file.Seek(0, 0)
    if err != nil {
            return written, dserr.Err(err)
    }
    // Read and write date
    written, err = copyBytes(file, writer, block.dataSize)
    if err != nil {
            return written, dserr.Err(err)
    }
    // Sync block descr
    err = block.updateBlockDescr()
    if err != nil {
            return written, dserr.Err(err)
    }
    return written, dserr.Err(err)
}

//func (block *Block) Clean() error {
//    var err error
//    fullFilePath := filepath.Join(block.baseDir, block.filePath)
//    // Remove file
//    if block.blockSize < 1 {        // todo: more strong validation
//        err = os.Remove(fullFilePath)
//        if err != nil {
//                return dserr.Err(err)
//        }
//    }
//    // Clean metadata
//    block.dataSize = 0
//    block.savedLoc = false
//    block.hashSum = make([]byte, 0)
//    // Sync block descr
//    err = block.updateBlockDescr()
//    if err != nil {
//            return dserr.Err(err)
//    }
//    return dserr.Err(err)
//}

func (block *Block) Erase() error {
    var err error
    fullFilePath := filepath.Join(block.baseDir, block.filePath)
    // Remove file
    if block.blockSize < 1 {        // todo: more strong validation
        err = os.Remove(fullFilePath)
        if err != nil {
                return dserr.Err(err)
        }
    }
    // Erase block descr
    err = block.reg.EraseBlockDescr(block.fileId, block.batchId, block.blockId, block.blockType)
    if err != nil {
            return dserr.Err(err)
    }
    return dserr.Err(err)
}

func (block *Block) addBlockDescr() error {
    descr := block.toDescr()
    return block.reg.AddBlockDescr(descr)
}

func (block *Block) updateBlockDescr() error {
    descr := block.toDescr()
    return block.reg.UpdateBlockDescr(descr)
}


func (block *Block) toDescr() *dscom.BlockDescr {
    descr := dscom.NewBlockDescr()
    descr.FileId    = block.fileId
    descr.BatchId   = block.batchId
    descr.BlockId   = block.blockId
    descr.BlockSize = block.blockSize
    descr.DataSize  = block.dataSize
    descr.FilePath  = block.filePath
    descr.BlockType = block.blockType
    descr.HashAlg   = block.hashAlg
    descr.HashInit  = hex.EncodeToString(block.hashInit)
    descr.HashSum   = hex.EncodeToString(block.hashSum)
    descr.FStoreId  = block.fstoreId
    descr.BStoreId  = block.bstoreId
    descr.SavedLoc  = block.savedLoc
    descr.SavedRem  = block.savedRem
    descr.LocUpdated = block.locUpdated
    return descr
}

func makeFilePath(fileId, batchId, blockId int64, blockType string) string {
    const blockFileExt string = ".blk"
    origin := fmt.Sprintf("%020d-%020d-%020d-%s", fileId, batchId, blockId, blockType)
    hasher := sha1.New()
    hasher.Write([]byte(origin))
    hashSum := hasher.Sum(nil)
    hashHex := hex.EncodeToString(hashSum)
    fileName := hashHex + blockFileExt
    l1 := string(hashHex[0:1])
    l2 := string(hashHex[1:3])
    dirName := filepath.Join(l1, l2)
    return filepath.Join(dirName, fileName)
}

func copyBytes(reader io.Reader, writer io.Writer, size int64) (int64, error) {
    var err error
    var bufSize int64 = 1024 * 8
    var total   int64 = 0
    var remains int64 = size
    buffer := make([]byte, bufSize)

    for {
        if remains == 0 {
            return total, dserr.Err(err)
        }
        if remains < bufSize {
            bufSize = remains
        }
        received, err := reader.Read(buffer[0:bufSize])
        if err != nil {
            return total, dserr.Err(err)
        }
        written, err := writer.Write(buffer[0:received])
        if err != nil {
            return total, dserr.Err(err)
        }
        if written != received {
            err = errors.New("write error")
            return total, dserr.Err(err)
        }
        total += int64(written)
        remains -= int64(written)
    }
    return total, dserr.Err(err)
}
