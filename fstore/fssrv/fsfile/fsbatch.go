package fsfile

import (
    "errors"
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
        err = errors.New("batch yet exists")
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
        err = errors.New("batch not exist")
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
            blockErr = true
        }
    }
    if blockErr {
        err = errors.New("some block not open")
        return &batch, dserr.Err(err)
    }
    err = batch.reg.IncSpecBatchDescrUC(batch.fileId, batch.batchId, batch.batchVer)
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
        err = errors.New("batch not exist")
        return &batch, dserr.Err(err)
    }

    batch.fileId    = descr.FileId
    batch.batchId   = descr.BatchId
    batch.batchSize = descr.BatchSize
    batch.blockSize = descr.BlockSize
    batch.batchVer  = descr.BatchVer

    batch.blocks = make([]*Block, batch.batchSize)
    //blockType := dscom.BTypeData
    //blockErr := false
    //for i := int64(0); i < batch.batchSize; i++ {
    //    blockId := i
    //    block, err := OpenBlock(reg, baseDir, fileId, batchId, blockId, blockType)
    //    if block != nil {
    //        batch.blocks[i] = block
    //    }
    //    if err != nil {
    //        blockErr = true
    //    }
    //}
    //if blockErr {
    //    err = errors.New("some block not open")
    //    return &batch, dserr.Err(err)
    //}
    //err = batch.reg.IncSpecBatchDescrUC(batch.fileId, batch.batchId, batch.batchVer)
    //if err != nil {
    //        return &batch, dserr.Err(err)
    //}
    //batch.batchIsOpen = true
    return &batch, dserr.Err(err)
}

func (batch *Batch) Write(reader io.Reader, need int64) (int64, error) {
    var err error
    var written int64
    if !batch.batchIsOpen {
        err = errors.New("batch not open or open witch error")
        return written, dserr.Err(err)
    }
    if batch.batchIsDeleted {
        err = errors.New("batch id deleted")
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
    return written, dserr.Err(err)
}


func (batch *Batch) Read(writer io.Writer) (int64, error) {
    var err error
    var read int64
    if !batch.batchIsOpen {
        err = errors.New("batch not open or open witch error")
        return read, dserr.Err(err)
    }
    if batch.batchIsDeleted {
        err = errors.New("batch id deleted")
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


func (batch *Batch) Distribute(distr dscom.IBlockDistr) error {
    var err error
    if !batch.batchIsOpen {
        err = errors.New("batch not open or open witch error")
        return dserr.Err(err)
    }
    if batch.batchIsDeleted {
        err = errors.New("batch id deleted")
        return dserr.Err(err)
    }
    for i := int64(0); i < batch.batchSize; i++ {
        err := batch.blocks[i].Distribute(distr)
        if err != nil {
            return dserr.Err(err)
        }
    }
    return dserr.Err(err)
}

func (batch *Batch) Delete() error {
    var err error
    // Return if wrong block
    if !batch.batchIsOpen {
        err = errors.New("batch not open or open with error")
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
        err = batch.reg.DecSpecBatchDescrUC(descr.FileId, descr.BatchId, descr.BatchVer)
        if err != nil {
                return dserr.Err(err)
        }
        batch.batchIsOpen = false
    }
    if !batch.batchIsDeleted {
        // Descrease usage counter of the batch descr
        err = batch.reg.DecSpecBatchDescrUC(batch.fileId, batch.batchId, batch.batchVer)
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
    if batch.batchIsOpen {
        err = batch.reg.DecSpecBatchDescrUC(batch.fileId, batch.batchId, batch.batchVer)
        if err != nil {
                return dserr.Err(err)
        }
        batch.batchIsOpen = false
    }
    // Remove underline blocks
    //for i := 0; i < batch.countBlocks(); i++ {
    //    block := batch.blocks[i]
    //    if block != nil {
    //        err := block.Erase()
    //        if err != nil {
    //            return dserr.Err(err)
    //        }
    //    }
    //}
    // Descrease usage counter of the batch descr
    err = batch.reg.EraseSpecBatchDescr(batch.fileId, batch.batchId, batch.batchVer)
    if err != nil {
            return dserr.Err(err)
    }
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
        err = batch.reg.DecSpecBatchDescrUC(batch.fileId, batch.batchId, batch.batchVer)
        if err != nil {
                return dserr.Err(err)
        }
        batch.batchIsOpen = false
    }
    return dserr.Err(err)
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
    return descr
}
