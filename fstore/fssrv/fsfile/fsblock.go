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
    "io"
    "io/fs"
    "math/rand"
    "os"
    "path/filepath"
    "time"

    "github.com/minio/highwayhash"

    "ndstore/dscom"
    "ndstore/dserr"
    "ndstore/dslog"
    "ndstore/dsrpc"

    "ndstore/bstore/bsfunc"
)

type Block struct {
    reg         dscom.IBlockReg
    baseDir     string

    fileId      int64
    batchId     int64
    blockId     int64
    blockType   string
    blockVer    int64

    blockSize   int64
    dataSize    int64
    filePath    string

    hashAlg     string
    hashSum     []byte
    hashInit    []byte

    createdAt   int64
    updatedAt   int64

    savedLoc    bool
    savedRem    bool
    locUpdated  bool
    bstoreId    int64
    fstoreId    int64

    blockIsOpen     bool
    blockIsDeleted   bool
}

const fileMode fs.FileMode  = 0644
const dirMode fs.FileMode   = 0755

func NewBlock(reg dscom.IBlockReg, baseDir string, fileId, batchId, blockId int64, blockType string,
                                                                blockSize int64) (*Block, error) {
    var block Block
    var err error

    block.reg       = reg
    block.baseDir   = baseDir

    exists, _, err := reg.GetNewestBlockDescr(fileId, batchId, blockId, blockType)
    if err != nil {
        return &block, dserr.Err(err)
    }
    if exists {
        err = fmt.Errorf("block %s yet exists", block.getIdString())
        return &block, dserr.Err(err)
    }
    newBlock, err := ForceNewBlock(reg, baseDir, fileId, batchId, blockId, blockType, blockSize)
    return newBlock, dserr.Err(err)
}

func ForceNewBlock(reg dscom.IBlockReg, baseDir string, fileId, batchId, blockId int64, blockType string,
                                                                blockSize int64) (*Block, error) {
    var block Block
    var err error

    block.reg       = reg
    block.baseDir   = baseDir

    block.fileId    = fileId
    block.batchId   = batchId
    block.blockId   = blockId
    block.blockType = blockType
    block.blockVer  = 0

    block.blockSize = blockSize
    block.dataSize  = 0
    block.filePath  = makeFilePath()

    block.hashAlg   = dscom.HashTypeHW
    block.hashInit  = make([]byte, 32)
    rand.Read(block.hashInit)
    block.hashSum   = make([]byte, 0)

    block.createdAt = time.Now().Unix()
    block.updatedAt = block.createdAt

    block.savedLoc  = false
    block.savedRem  = false
    block.bstoreId  = 0
    block.fstoreId  = 0
    block.locUpdated = false

    block.blockVer  = time.Now().UnixNano()
    descr := block.toDescr()
    descr.UCounter = 2
    err = block.reg.AddNewBlockDescr(descr)
    if err != nil {
        return &block, dserr.Err(err)
    }
    block.blockIsDeleted = false
    block.blockIsOpen   = true
    return &block, dserr.Err(err)
}

func OpenBlock(reg dscom.IBlockReg, baseDir string, fileId, batchId, blockId int64, blockType string) (*Block, error) {
    var err error
    var block Block

    block.baseDir   = baseDir
    block.reg       = reg

    exists, descr, err := reg.GetNewestBlockDescr(fileId, batchId, blockId, blockType)
    if err != nil {
        return &block, dserr.Err(err)
    }
    if !exists {
        err = fmt.Errorf("block not exists", block.getIdString())
        return &block, dserr.Err(err)
    }

    err = block.reg.IncSpecBlockDescrUC(1, descr.FileId, descr.BatchId, descr.BlockId, descr.BlockType, descr.BlockVer)
    if err != nil {
            return &block, dserr.Err(err)
    }

    block.fileId    = descr.FileId
    block.batchId   = descr.BatchId
    block.blockId   = descr.BlockId
    block.blockType = descr.BlockType
    block.blockVer  = descr.BlockVer

    block.blockSize = descr.BlockSize
    block.dataSize  = descr.DataSize
    block.filePath  = descr.FilePath
    block.hashAlg   = descr.HashAlg

    block.hashInit, err = hex.DecodeString(descr.HashInit)
    if err != nil {
        return &block, dserr.Err(err)
    }
    block.hashSum, err   = hex.DecodeString(descr.HashSum)
    if err != nil {
        return &block, dserr.Err(err)
    }
    block.savedLoc  = descr.SavedLoc
    block.savedRem  = descr.SavedRem
    block.fstoreId  = descr.FStoreId
    block.bstoreId  = descr.BStoreId
    block.locUpdated = descr.LocUpdated

    block.blockIsDeleted = false
    block.blockIsOpen   = true
    return &block, dserr.Err(err)
}


func OpenSpecUnusedBlock(reg dscom.IBlockReg, baseDir string, fileId, batchId, blockId int64, blockType string,
                                                        blockVer int64) (*Block, error) {
    var err error
    var block Block
    exists, descr, err := reg.GetSpecUnusedBlockDescr(fileId, batchId, blockId, blockType, blockVer)
    if err != nil {
        return &block, dserr.Err(err)
    }
    if !exists {
        err = fmt.Errorf("block not exists", block.getIdString())
        return &block, dserr.Err(err)
    }

    block.baseDir   = baseDir
    block.reg       = reg

    block.fileId    = descr.FileId
    block.batchId   = descr.BatchId
    block.blockId   = descr.BlockId
    block.blockType = descr.BlockType
    block.blockVer  = descr.BlockVer

    block.blockSize = descr.BlockSize
    block.dataSize  = descr.DataSize
    block.filePath  = descr.FilePath
    block.hashAlg   = descr.HashAlg

    block.hashInit, err = hex.DecodeString(descr.HashInit)
    if err != nil {
        return &block, dserr.Err(err)
    }
    block.hashSum, err   = hex.DecodeString(descr.HashSum)
    if err != nil {
        return &block, dserr.Err(err)
    }
    block.savedLoc  = descr.SavedLoc
    block.savedRem  = descr.SavedRem
    block.fstoreId  = descr.FStoreId
    block.bstoreId  = descr.BStoreId
    block.locUpdated = descr.LocUpdated

    block.blockIsDeleted = false
    block.blockIsOpen   = true
    return &block, dserr.Err(err)
}


func (block *Block) Write(reader io.Reader, need int64) (int64, error) {
    var err error
    var written int64
    // Return if wrong block
    if !block.blockIsOpen {
        err = fmt.Errorf("block %s not open or open with error", block.getIdString())
        return written, dserr.Err(err)
    }
    // Return if block just erased
    if block.blockIsDeleted {
        err = fmt.Errorf("block %s is deleted", block.getIdString())
        return written, dserr.Err(err)
    }
    // Nothing if writing zero
    if need < 1 {
        return written, dserr.Err(err)
    }
    // Nothing if block full
    if block.dataSize >= block.blockSize {
        return written, dserr.Err(err)
    }
    // Prepare env
    newFilePath := makeFilePath()
    newFullFilePath := filepath.Join(block.baseDir, newFilePath)
    newDirPath := filepath.Dir(newFullFilePath)
    err = os.MkdirAll(newDirPath, dirMode)
    if err != nil {
            return written, dserr.Err(err)
    }
    // Open new version of underline file
    newFile, err := os.OpenFile(newFullFilePath, os.O_CREATE|os.O_WRONLY, fileMode)
    defer newFile.Close()
    if err != nil {
            return written, dserr.Err(err)
    }
    // Create hasher
    hasher, _ := highwayhash.New(block.hashInit)
    multiWriter := io.MultiWriter(newFile, hasher)
    // Copy old version data if need
    if block.dataSize > 0 {
        fullFilePath := filepath.Join(block.baseDir, block.filePath)
        file, err := os.OpenFile(fullFilePath, os.O_RDONLY, 0)
        defer file.Close()
        if err != nil {
                return written, dserr.Err(err)
        }
        reWritten, err := copyBytes(file, multiWriter, block.dataSize)
        if err != nil {
                return written, dserr.Err(err)
        }
        if reWritten != block.dataSize {
            err = fmt.Errorf("incorrect prev block %s file size", block.getIdString())
            return written, dserr.Err(err)
        }
    }
    // Read and write remain date
    remains := block.blockSize - block.dataSize
    if remains > need {
        remains = need
    }
    written, err = copyBytes(reader, multiWriter, remains)
    if err == io.EOF {
        err = fmt.Errorf("block %s reading error %v", block.getIdString(), err)
        return written, dserr.Err(err)
    }
    if err != nil {
        return written, dserr.Err(err)
    }
    // Add new version of block descr
    descr := block.toDescr()
    block.dataSize  += written
    block.hashSum    = hasher.Sum(nil)
    block.savedLoc   = true
    block.locUpdated = true
    block.filePath   = newFilePath
    block.blockVer   = time.Now().UnixNano()
    block.updatedAt  = time.Now().Unix()
    newDescr := block.toDescr()
    newDescr.UCounter = 2
    err = block.reg.AddNewBlockDescr(newDescr)
    if err != nil {
            return written, dserr.Err(err)
    }
    // Descreate usage old block descr
    err = block.reg.DecSpecBlockDescrUC(2, descr.FileId, descr.BatchId, descr.BlockId, descr.BlockType, descr.BlockVer)
    if err != nil {
            return written, dserr.Err(err)
    }
    return written, dserr.Err(err)
}



func (block *Block) Read(writer io.Writer) (int64, error) {
    var err error
    var written int64
    // Return with error if wrong block
    if !block.blockIsOpen {
        err = fmt.Errorf("block %s not open or open with error", block.getIdString())
        return written, dserr.Err(err)
    }
    // Return with error if block just erased
    if block.blockIsDeleted {
        err = fmt.Errorf("block %s is deleted", block.getIdString())
        return written, dserr.Err(err)
    }
    // Return if block is empty
    if block.dataSize < 1 {
        return written, dserr.Err(err)
    }
    if !block.savedLoc && !block.savedRem {
        err = fmt.Errorf("block %s not stored remote or local", block.getIdString())
        return written, dserr.Err(err)
    }

    // Get block from remote store
    var fullFilePath string

    switch {
        case block.savedLoc:
            fullFilePath = filepath.Join(block.baseDir, block.filePath)
            file, err := os.OpenFile(fullFilePath, os.O_RDONLY, 0)
            defer file.Close()
            if err != nil {
                    return written, dserr.Err(err)
            }
            // Check hashsum and size of local file
            hasher, _ := highwayhash.New(block.hashInit)
            written, err = copyBytes(file, hasher, block.dataSize)
            if err != nil {
                return written, dserr.Err(err)
            }
            if written != block.dataSize {
                err = fmt.Errorf("incorrect block %s local file size", block.getIdString())
                return written, dserr.Err(err)
            }
            hashSum := hasher.Sum(nil)
            if bytes.Compare(hashSum, block.hashSum) != 0 {
                err = fmt.Errorf("incorrect block %s hash sum", block.getIdString())
                return written, dserr.Err(err)
            }
        default: // Load from remote block store
            bstoreExists, bstoreDescr, err := block.reg.GetBStoreDescrById(block.bstoreId)
            if err != nil {
                return written, dserr.Err(err)
            }
            if !bstoreExists {
                err = fmt.Errorf("bstore %d not exists", block.bstoreId)
                return written, dserr.Err(err)
            }

            filePath := makeTmpFilePath()
            fullFilePath = filepath.Join(block.baseDir, filePath)
            dirPath := filepath.Dir(fullFilePath)
            err = os.MkdirAll(dirPath, dirMode)
            if err != nil {
                    return written, dserr.Err(err)
            }
            file, err := os.OpenFile(fullFilePath, os.O_CREATE|os.O_TRUNC|os.O_WRONLY, 0)
            exitFunc := func() {
                file.Close()
                os.Remove(fullFilePath)
            }
            defer exitFunc()
            if err != nil {
                    return written, dserr.Err(err)
            }
            uri     := fmt.Sprintf("%s:%s", bstoreDescr.Address, bstoreDescr.Port)
            login   := bstoreDescr.Login
            pass    := bstoreDescr.Pass
            auth    := dsrpc.CreateAuth([]byte(login), []byte(pass))
             //Write remote date to multiwriter and check hashsum
            hasher, _ := highwayhash.New(block.hashInit)
            mWriter := io.MultiWriter(file, hasher)
            err = bsfunc.LoadBlock(uri, auth, block.fileId, block.batchId, block.blockId, mWriter,
                                                                    block.blockType, block.blockVer)
            if err != nil {
                    return written, dserr.Err(err)
            }
            fileStat, err := file.Stat()
            if err != nil {
                    return written, dserr.Err(err)
            }
            if fileStat.Size() != block.dataSize {
                err = fmt.Errorf("incorrect block %s local file size %s: %d and block size %d",
                                    block.getIdString(), fullFilePath, fileStat.Size(), block.dataSize)
                return written, dserr.Err(err)
            }
            hashSum := hasher.Sum(nil)
            if bytes.Compare(hashSum, block.hashSum) != 0 {
                err = fmt.Errorf("incorrect block %s hash sum", block.getIdString())
                return written, dserr.Err(err)
            }
    }
    file, err := os.OpenFile(fullFilePath, os.O_RDONLY, fileMode)
    defer file.Close()
    written, err = copyBytes(file, writer, block.dataSize)
    if err != nil {
            return written, dserr.Err(err)
    }
    return written, dserr.Err(err)
}

func (block *Block) Clean() error {
    var err error
    if !block.blockIsOpen {
        err = fmt.Errorf("block %s not open or open with error")
        return dserr.Err(err)
    }
    if block.blockIsDeleted {
        err = fmt.Errorf("block %s is deleted")
        return dserr.Err(err)
    }
    // Close block if open
    if block.blockIsOpen {
        descr := block.toDescr()
        err = block.reg.DecSpecBlockDescrUC(1, descr.FileId, descr.BatchId, descr.BlockId, descr.BlockType, descr.BlockVer)
        if err != nil {
                return dserr.Err(err)
        }
        block.blockIsOpen = false
    }
    // Remove local data
    fullFilePath := filepath.Join(block.baseDir, block.filePath)
    err = removeFile(fullFilePath)
    if err != nil {
            return dserr.Err(err)
    }
    // Add new version of block descr
    descr := block.toDescr()
    block.dataSize  = 0
    block.savedLoc  = false
    block.hashSum   = make([]byte, 0)
    block.blockVer  = time.Now().UnixNano()
    newDescr := block.toDescr()
    newDescr.UCounter = 2
    err = block.reg.AddNewBlockDescr(newDescr)
    if err != nil {
            return dserr.Err(err)
    }
    // Descreate usage old block descr
    err = block.reg.DecSpecBlockDescrUC(2, descr.FileId, descr.BatchId, descr.BlockId, descr.BlockType, descr.BlockVer)
    if err != nil {
            return dserr.Err(err)
    }
    return dserr.Err(err)
}


func (block *Block) Distribute(distr dscom.IFileDistr) (bool, error) {
    var err error

    blockIsDistibuted := false
    var bstoreId int64 = 0

    // Return id not opened
    if !block.blockIsOpen {
        err = fmt.Errorf("block %s not open or open with error", block.getIdString())
        return blockIsDistibuted, dserr.Err(err)
    }
    // Return if reased
    if block.blockIsDeleted {
        err = fmt.Errorf("block %s is deleted", block.getIdString())
        return blockIsDistibuted, dserr.Err(err)
    }

    if block.savedRem && !block.locUpdated {
        blockIsDistibuted = true
        return blockIsDistibuted, dserr.Err(err)
    }
    // Upload block
    // Store current state to descr
    remDescr := block.toDescr()
    newBlockVer := time.Now().UnixNano()
    remDescr.BlockVer = newBlockVer
    switch  {
        case block.dataSize > 0:
            blockIsDistibuted, bstoreId, err = distr.SaveBlock(remDescr);
            if err != nil {
                return blockIsDistibuted, dserr.Err(err)
            }
        default:
            blockIsDistibuted = true
            bstoreId = 0
    }
    if !blockIsDistibuted {
        err = fmt.Errorf("block %s is not distributed", block.getIdString())
        return blockIsDistibuted, dserr.Err(err)
    }

    // Make new filename patch
    newFilePath := makeFilePath()
    descr := block.toDescr()

    savedLoc := true
    // Link data to new file name
    if block.dataSize > 0 {
        newFullFilePath := filepath.Join(block.baseDir, newFilePath)
        newDirPath := filepath.Dir(newFullFilePath)

        oldFullFilePath := filepath.Join(block.baseDir, block.filePath)
        err = os.MkdirAll(newDirPath, dirMode)
        if err != nil {
                return blockIsDistibuted, dserr.Err(err)
        }
        err = os.Link(oldFullFilePath, newFullFilePath)
        if err != nil {
                return blockIsDistibuted, dserr.Err(err)
        }
        removeErr := os.Remove(oldFullFilePath)
        if removeErr == nil {
            savedLoc = false
        }
    }
    // Add new version of block descr
    block.filePath  = newFilePath
    block.savedRem  = blockIsDistibuted
    block.savedLoc  = savedLoc
    block.bstoreId  = bstoreId
    block.blockVer  = newBlockVer
    newDescr := block.toDescr()
    newDescr.UCounter = 2
    err = block.reg.AddNewBlockDescr(newDescr)
    if err != nil {
        return blockIsDistibuted, dserr.Err(err)
    }
    // Decrease usage old block descr
    err = block.reg.DecSpecBlockDescrUC(2, descr.FileId, descr.BatchId, descr.BlockId, descr.BlockType, descr.BlockVer)
    if err != nil {
        return blockIsDistibuted, dserr.Err(err)

    }
    dslog.LogDebugf("block %s save to store %d ", block.getIdString(), bstoreId)
    return blockIsDistibuted, dserr.Err(err)
}

func (block *Block) Delete() error {
    var err error
    // Return if wrong block
    if !block.blockIsOpen {
        err = fmt.Errorf("block %s not open or open with error", block.getIdString())
        return dserr.Err(err)
    }
    // Close block if open
    if block.blockIsOpen {
        descr := block.toDescr()
        err = block.reg.DecSpecBlockDescrUC(1, descr.FileId, descr.BatchId, descr.BlockId, descr.BlockType, descr.BlockVer)
        if err != nil {
                return dserr.Err(err)
        }
        block.blockIsOpen = false
    }
    // Descrease usage counter of the block descr
    err = block.reg.DecSpecBlockDescrUC(1, block.fileId, block.batchId, block.blockId, block.blockType, block.blockVer)
    if err != nil {
            return dserr.Err(err)
    }
    block.blockIsDeleted = true
    return dserr.Err(err)
}

func (block *Block) Erase() error {
    var err error
    // Close block
    if block.blockIsOpen {
        err = block.reg.DecSpecBlockDescrUC(1, block.fileId, block.batchId, block.blockId, block.blockType, block.blockVer)
        if err != nil {
                return dserr.Err(err)
        }
        block.blockIsOpen = false
    }
    if block.savedLoc && block.dataSize > 0 {
        // Remove underline file
        if len(block.filePath) > 0 {
            fullFilePath := filepath.Join(block.baseDir, block.filePath)
            err = removeFile(fullFilePath)
            if err != nil {
                    return dserr.Err(err)
            }
        }
    }
    if block.savedRem && block.dataSize > 0 {
            bstoreExists, bstoreDescr, err := block.reg.GetBStoreDescrById(block.bstoreId)
            if err != nil {
                return dserr.Err(err)
            }
            if !bstoreExists {
                err = fmt.Errorf("bstore %d not exists", block.bstoreId)
                return dserr.Err(err)
            }
            uri     := fmt.Sprintf("%s:%s", bstoreDescr.Address, bstoreDescr.Port)
            login   := bstoreDescr.Login
            pass    := bstoreDescr.Pass
            auth    := dsrpc.CreateAuth([]byte(login), []byte(pass))
            count := 3
            for i := 0; i < count; i++ {
                // Delete remote block
                err = bsfunc.DeleteBlock(uri, auth, block.fileId, block.batchId, block.blockId, block.blockType, block.blockVer)
                if err == nil {
                    break
                }
            }
            if err != nil {
                    return dserr.Err(err)
            }
    }
    // Erase block descr
    err = block.reg.EraseSpecBlockDescr(block.fileId, block.batchId, block.blockId, block.blockType, block.blockVer)
    if err != nil {
            return dserr.Err(err)
    }
    block.blockIsDeleted = true
    return dserr.Err(err)
}



func (block *Block) Close() error {
    var err error
    if block.blockIsDeleted {
        return dserr.Err(err)
    }
    if block.blockIsOpen {
        descr := block.toDescr()
        err = block.reg.DecSpecBlockDescrUC(1, descr.FileId, descr.BatchId, descr.BlockId, descr.BlockType, descr.BlockVer)
        if err != nil {
                return dserr.Err(err)
        }
        block.blockIsOpen = false
        return dserr.Err(err)
    }
    return dserr.Err(err)
}

func removeFile(filePath string) error {
    var err error
    _, err = os.Stat(filePath)
    if err == nil {
        err = os.Remove(filePath)
        if err != nil {
            return dserr.Err(err)
        }
    }
    err = nil
    return err
}


func (block *Block) getIdString() string {
    return fmt.Sprintf("%d,%d,%d,%s,%d", block.fileId, block.batchId, block.blockId, block.blockType, block.blockVer)
}

func (block *Block) toDescr() *dscom.BlockDescr {
    descr := dscom.NewBlockDescr()
    descr.FileId    = block.fileId
    descr.BatchId   = block.batchId
    descr.BlockId   = block.blockId
    descr.BlockType = block.blockType
    descr.BlockVer  = block.blockVer

    descr.BlockSize = block.blockSize
    descr.DataSize  = block.dataSize
    descr.FilePath  = block.filePath

    descr.HashAlg   = block.hashAlg
    descr.HashInit  = hex.EncodeToString(block.hashInit)
    descr.HashSum   = hex.EncodeToString(block.hashSum)

    descr.SavedLoc  = block.savedLoc
    descr.SavedRem  = block.savedRem
    descr.BStoreId  = block.bstoreId
    descr.FStoreId  = block.fstoreId
    descr.LocUpdated = block.locUpdated
    return descr
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

func makeTmpFilePath() string {
    const blockFileExt string = ".tmpbl"
    origin := make([]byte, 128)
    rand.Read(origin)
    hasher := sha1.New()
    hasher.Write(origin)
    hashSum := hasher.Sum(nil)
    hashHex := hex.EncodeToString(hashSum)
    fileName := hashHex + blockFileExt
    l1 := string(hashHex[0:1])
    l2 := string(hashHex[1:3])
    dirName := filepath.Join("tmp", l1, l2)
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
