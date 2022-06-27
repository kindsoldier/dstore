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
    IBStoreReg
    IUserReg
    IEntryReg

    IFileReg
    IBatchReg
    IBlockReg
}

type IBlockReg interface {
    AddBlockDescr(descr *BlockDescr) error
    GetBlockDescr(fileId, batchId, blockId int64, blockType string) (bool, *BlockDescr, error)
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
    GetUnusedFileDescr() (bool, *FileDescr, error)
    GetLostedFileDescr() (bool, *FileDescr, error)
}

type IBStoreReg interface {
    AddBStoreDescr(descr *BStoreDescr) (int64, error)
    EraseBStoreDescr(address, port string) error
    GetBStoreDescr(address, port string) (bool, *BStoreDescr, error)
    ListBStoreDescrs() ([]*BStoreDescr, error)
    UpdateBStoreDescr(descr *BStoreDescr) error
}

type IUserReg interface {
    AddUserDescr(descr *UserDescr) (int64, error)
    EraseUserDescr(login string) error
    GetUserDescr(login string) (bool, *UserDescr, error)
    ListUserDescrs() ([]*UserDescr, error)
    UpdateUserDescr(descr *UserDescr) error
}

type IEntryReg interface {
    AddEntryDescr(userId int64, dirPath, fileName string, fileId int64) error
    EraseEntryDescr(userId int64, dirPath, fileName string) error
    EntryDescrExists(userId int64, dirPath, fileName string) (bool, error)
    GetEntryDescr(userId int64, dirPath, fileName string) (bool, *EntryDescr, error)
    ListEntryDescr(userId int64, dirPath string) ([]*EntryDescr, error)
}
