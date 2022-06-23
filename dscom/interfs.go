/*
 * Copyright 2022 Oleg Borodin  <borodin@unix7.org>
 */

package dscom

import (
    "io"
)
type IFileSender interface {
}

type IFile interface {
    Read(writer io.Writer) (int64, error)
    Write(reader io.Reader, need int64) (int64, error)
    Erase() error
    Close() error
}

type IBatch interface {
    Read(writer io.Writer) (int64, error)
    Write(reader io.Reader, need int64) (int64, error)
    Erase() error
    Close() error
}

type IBlock interface {
    Read(writer io.Writer) (int64, error)
    Write(reader io.Reader, need int64) (int64, error)
    Erase() error
    Close() error
}

type IFSReg interface {
    IFileReg
    IBatchReg
    IBlockReg
}

type IBlockReg interface {
    AddBlockDescr(descr *BlockDescr) error
    GetBlockDescr(fileId, batchId, blockId int64, blockType string) (bool, *BlockDescr, error)
    GetUnusedBlockDescr() (bool, *BlockDescr, error)
    //ListBlockDescrs() ([]*BlockDescr, error)
    ListBlockDescrsByFileId(fileId int64) ([]*BlockDescr, error)
    UpdateBlockDescr(descr *BlockDescr) error
    EraseBlockDescr(fileId, batchId, blockId int64, blockType string) error
}

type IBatchReg interface {
    AddBatchDescr(descr *BatchDescr) error
    EraseBatchDescr(fileId, batchId int64) error
    GetBatchDescr(fileId, batchId int64) (bool, *BatchDescr, error)
    ListBatchDescrsByFileId(fileId int64) ([]*BatchDescr, error)
    UpdateBatchDescr(descr *BatchDescr) error
}

type IFileReg interface {
    AddFileDescr(descr *FileDescr) (int64, error)
    UpdateFileDescr(descr *FileDescr) error
    EraseFileDescr(fileId int64) error
    GetFileDescr(fileId int64) (bool, *FileDescr, error)
    IncFileDescrUC(fileId int64) error
    DecFileDescrUC(fileId int64) error
    ListFileDescrs() ([]*FileDescr, error)
}
