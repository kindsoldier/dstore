package fsfile

import (
    "fmt"
    "io"
    "time"
    "ndstore/dscom"
    "ndstore/dserr"
)

type Batch struct {
    reg         dscom.IFSReg
    baseDir     string

    fileId      int64
    batchId     int64
    batchVer    int64

    batchSize   int64
    blockSize   int64

    createdAt   int64
    updatedAt   int64

    blocks      []*Block

    batchIsDeleted  bool
    batchIsOpen     bool
}

func NewBatch(reg dscom.IFSReg, baseDir string, fileId, batchId, batchSize, blockSize int64) (*Batch, error) {
    var batch Batch
    var err error

    batch.reg       = reg
    batch.baseDir   = baseDir

    exists, _, err := reg.GetNewestBatchDescr(fileId, batchId)
    if err != nil {
        return &batch, dserr.Err(err)
    }
    if exists {
        err = fmt.Errorf("batch %s yet exists", batch.getIdString())
        return &batch, dserr.Err(err)
    }
    newBatch, err := ForceNewBatch(reg, baseDir, fileId, batchId, batchSize, blockSize)
    return newBatch, dserr.Err(err)
}


func ForceNewBatch(reg dscom.IFSReg, baseDir string, fileId, batchId, batchSize, blockSize int64) (*Batch, error) {
    var batch Batch
    var err error

    batch.reg       = reg
    batch.baseDir   = baseDir

    batch.fileId    = fileId
    batch.batchId   = batchId
    batch.batchSize = batchSize
    batch.blockSize = blockSize

    batch.createdAt = time.Now().Unix()
    batch.updatedAt = batch.createdAt

    blockType := dscom.BTypeData
    batch.blocks = make([]*Block, batch.batchSize)

    for i := int64(0); i < batch.batchSize; i++ {
        blockId := i
        block, err := OpenBlock(reg, baseDir, fileId, batchId, blockId, blockType)
        if block != nil {
            err = block.Erase()
            if err != nil {
                return &batch, dserr.Err(err)
            }
        }
        batch.blocks[i] = block
    }

    for i := int64(0); i < batch.batchSize; i++ {
        blockId := i
        block, err := ForceNewBlock(reg, baseDir, fileId, batchId, blockId, blockType, blockSize)
        if err != nil {
            return &batch, dserr.Err(err)
        }
        batch.blocks[i] = block
    }

    batch.batchVer = time.Now().UnixNano()
    descr := batch.toDescr()
    descr.UCounter  = 2
    err = batch.reg.AddNewBatchDescr(descr)
    if err != nil {
        return &batch, dserr.Err(err)
    }
    batch.batchIsDeleted    = false
    batch.batchIsOpen       = true
    return &batch, dserr.Err(err)
}

func OpenBatch(reg dscom.IFSReg, baseDir string, fileId, batchId int64) (*Batch, error) {
    var err error
    var batch Batch

    batch.baseDir   = baseDir
    batch.reg       = reg

    exists, descr, err := reg.GetNewestBatchDescr(fileId, batchId)
    if err != nil {
        return &batch, dserr.Err(err)
    }
    if !exists {
        err = fmt.Errorf("batch %s not exist", batch.getIdString())
        return &batch, dserr.Err(err)
    }

    batch.fileId    = descr.FileId
    batch.batchId   = descr.BatchId
    batch.batchSize = descr.BatchSize
    batch.blockSize = descr.BlockSize
    batch.batchVer  = descr.BatchVer

    blockType := dscom.BTypeData
    batch.blocks = make([]*Block, batch.batchSize)
    blockErr := false
    for i := int64(0); i < batch.batchSize; i++ {
        blockId := i
        block, err := OpenBlock(reg, baseDir, fileId, batchId, blockId, blockType)
        if block != nil {
            batch.blocks[i] = block
        }
        if err != nil {
            block.Close()
            blockErr = true
        }
    }
    if blockErr {
        err = fmt.Errorf("some block in batch %s not open", batch.getIdString())
        return &batch, dserr.Err(err)
    }
    err = batch.reg.IncSpecBatchDescrUC(1, batch.fileId, batch.batchId, batch.batchVer)
    if err != nil {
            return &batch, dserr.Err(err)
    }
    batch.batchIsOpen = true
    return &batch, dserr.Err(err)
}

func ForceOpenBatch(reg dscom.IFSReg, baseDir string, fileId, batchId int64) (*Batch, error) {
    var err error
    var batch Batch

    batch.baseDir   = baseDir
    batch.reg       = reg

    exists, descr, err := reg.GetNewestBatchDescr(fileId, batchId)
    if err != nil {
        return &batch, dserr.Err(err)
    }
    if !exists {
        err = fmt.Errorf("batch %s not exist", batch.getIdString())
        return &batch, dserr.Err(err)
    }

    batch.fileId    = descr.FileId
    batch.batchId   = descr.BatchId
    batch.batchSize = descr.BatchSize
    batch.blockSize = descr.BlockSize
    batch.batchVer  = descr.BatchVer

    blockType := dscom.BTypeData
    batch.blocks = make([]*Block, batch.batchSize)
    for i := int64(0); i < batch.batchSize; i++ {
        blockId := i
        block, _ := OpenBlock(reg, baseDir, fileId, batchId, blockId, blockType)
        if block != nil {
            batch.blocks[i] = block
        }
    }
    err = batch.reg.IncSpecBatchDescrUC(1, batch.fileId, batch.batchId, batch.batchVer)
    if err != nil {
            return &batch, dserr.Err(err)
    }
    batch.batchIsOpen = true
    return &batch, dserr.Err(err)
}


func OpenSpecUnusedBatch(reg dscom.IFSReg, baseDir string, fileId, batchId, batchVer int64) (*Batch, error) {
    var err error
    var batch Batch

    batch.baseDir   = baseDir
    batch.reg       = reg

    exists, descr, err := reg.GetSpecUnusedBatchDescr(fileId, batchId, batchVer)
    if err != nil {
        return &batch, dserr.Err(err)
    }
    if !exists {
        err = fmt.Errorf("batch %s not exist", batch.getIdString())
        return &batch, dserr.Err(err)
    }

    batch.fileId    = descr.FileId
    batch.batchId   = descr.BatchId
    batch.batchSize = descr.BatchSize
    batch.blockSize = descr.BlockSize
    batch.batchVer  = descr.BatchVer

    batch.blocks = make([]*Block, batch.batchSize)

    return &batch, dserr.Err(err)
}

func (batch *Batch) Write(reader io.Reader, need int64) (int64, error) {
    var err error
    var written int64
    if !batch.batchIsOpen {
        err = fmt.Errorf("batch %s not open or open witch error", batch.getIdString())
        return written, dserr.Err(err)
    }
    if batch.batchIsDeleted {
        err = fmt.Errorf("batch %s is deleted", batch.getIdString())
        return written, dserr.Err(err)
    }
    for i := 0; i < batch.countBlocks(); i++ {
        if need < 1 {
            return written, err
        }
        blockWritten, err := batch.blocks[i].Write(reader, need)
        written += blockWritten
        if err == io.EOF {
            err = nil
            return written, dserr.Err(err)
        }
        if err != nil {
            return written, dserr.Err(err)
        }
        need -= blockWritten
    }

    //batch.updatedAt = time.Now().Unix()

    return written, dserr.Err(err)
}


func (batch *Batch) Read(writer io.Writer) (int64, error) {
    var err error
    var read int64
    if !batch.batchIsOpen {
        err = fmt.Errorf("batch %s not open or open witch error", batch.getIdString())
        return read, dserr.Err(err)
    }
    if batch.batchIsDeleted {
        err = fmt.Errorf("batch %s is deleted", batch.getIdString())
        return read, dserr.Err(err)
    }
    for i := 0; i < batch.countBlocks(); i++ {
        blockRead, err := batch.blocks[i].Read(writer)
        read += blockRead
        if err != nil {
            return read, dserr.Err(err)
        }
    }
    return read, dserr.Err(err)
}

func (batch *Batch) Distribute(distr dscom.IFileDistr) (bool, error) {
    var err error
    batchIsDistr := false
    if !batch.batchIsOpen {
        err = fmt.Errorf("batch %s not open or open witch error", batch.getIdString())
        return batchIsDistr, dserr.Err(err)
    }
    if batch.batchIsDeleted {
        err = fmt.Errorf("batch %s is deleted", batch.getIdString())
        return batchIsDistr, dserr.Err(err)
    }
    allBlockDistr := true
    for i := int64(0); i < batch.batchSize; i++ {
        blockIsDistributed, err := batch.blocks[i].Distribute(distr)
        if err != nil {
            return batchIsDistr, dserr.Err(err)
        }
        if !blockIsDistributed {
            allBlockDistr = false
        }
    }
    batchIsDistr = allBlockDistr
    return batchIsDistr, dserr.Err(err)
}

func (batch *Batch) Delete() error {
    var err error
    // Return if wrong block
    if !batch.batchIsOpen {
        err = fmt.Errorf("batch %s not open or open with error", batch.getIdString())
        return dserr.Err(err)
    }
    // Delete blocks
    for i := 0; i < batch.countBlocks(); i++ {
        if batch.blocks[i] != nil {
            err := batch.blocks[i].Delete()
            if err != nil {
                return dserr.Err(err)
            }
        }
    }
    // Close batch if open
    if batch.batchIsOpen {
        descr := batch.toDescr()
        err = batch.reg.DecSpecBatchDescrUC(1, descr.FileId, descr.BatchId, descr.BatchVer)
        if err != nil {
                return dserr.Err(err)
        }
        batch.batchIsOpen = false
    }
    if !batch.batchIsDeleted {
        // Descrease usage counter of the batch descr
        err = batch.reg.DecSpecBatchDescrUC(1, batch.fileId, batch.batchId, batch.batchVer)
        if err != nil {
                return dserr.Err(err)
        }
        batch.batchIsDeleted = true
    }
    return dserr.Err(err)
}

func (batch *Batch) Erase() error {
    var err error
    // Close block
    //if batch.batchIsOpen {
    //    err = batch.reg.DecSpecBatchDescrUC(1, batch.fileId, batch.batchId, batch.batchVer)
    //    if err != nil {
    //            return dserr.Err(err)
    //    }
    //    batch.batchIsOpen = false
    //}
    // Descrease usage counter of the batch descr
    err = batch.reg.EraseSpecBatchDescr(batch.fileId, batch.batchId, batch.batchVer)
    if err != nil {
            return dserr.Err(err)
    }
    batch.batchIsOpen = false
    batch.batchIsDeleted = true
    return dserr.Err(err)
}

func (batch *Batch) Close() error {
    var err error
    if batch.batchIsDeleted {
        return dserr.Err(err)
    }
    // Always close blocks
    for i := 0; i < batch.countBlocks(); i++ {
        if batch.blocks[i] != nil {
            batch.blocks[i].Close()
        }
    }
    if batch.batchIsOpen {
        err = batch.reg.DecSpecBatchDescrUC(1, batch.fileId, batch.batchId, batch.batchVer)
        if err != nil {
                return dserr.Err(err)
        }
        batch.batchIsOpen = false
    }
    return dserr.Err(err)
}

func (batch *Batch) getIdString() string {
    return fmt.Sprintf("%d,%d,%d", batch.fileId, batch.batchId, batch.batchVer)
}

func (batch *Batch) countBlocks() int {
    return len(batch.blocks)
}

func (batch *Batch) toDescr() *dscom.BatchDescr {
    descr := dscom.NewBatchDescr()
    descr.FileId    = batch.fileId
    descr.BatchId   = batch.batchId
    descr.BatchSize = batch.batchSize
    descr.BlockSize = batch.blockSize
    descr.BatchVer  = batch.batchVer
    descr.CreatedAt = batch.createdAt
    descr.UpdatedAt = batch.updatedAt
    return descr
}
