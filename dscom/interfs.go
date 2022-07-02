/*
 * Copyright 2022 Oleg Borodin  <borodin@unix7.org>
 */

package dscom

import (
    "io"
)
type IBlockDistr interface {
    SaveBlock(descr *BlockDescr) (int64, error)
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
    AddNewBlockDescr(descr *BlockDescr) error
    DecSpecBlockDescrUC(fileId, batchId, blockId int64, blockType string, blockVer int64) error
    EraseSpecBlockDescr(fileId, batchId, blockId int64, blockType string, blockVer int64) error
    GetAnyUnusedBlockDescr() (bool, *BlockDescr, error)
    GetNewestBlockDescr(fileId, batchId, blockId int64, blockType string) (bool, *BlockDescr, error)
    GetSpecBlockDescr(fileId, batchId, blockId int64, blockType string, blockVer int64) (bool, *BlockDescr, error)
    GetSpecUnusedBlockDescr(fileId, batchId, blockId int64, blockType string, blockVer int64) (bool, *BlockDescr, error)
    IncSpecBlockDescrUC(fileId, batchId, blockId int64, blockType string, blockVer int64) error
    ListAllBlockDescrs() ([]*BlockDescr, error)
}

type IBatchReg interface {
    AddNewBatchDescr(descr *BatchDescr) error
    DecSpecBatchDescrUC(fileId, batchId, batchVer int64) error
    EraseAllBatchDescrs() error
    EraseSpecBatchDescr(fileId, batchId, batchVer int64) error
    GetAnyUnusedBatchDescr() (bool, *BatchDescr, error)
    GetNewestBatchDescr(fileId, batchId int64) (bool, *BatchDescr, error)
    GetSpecBatchDescr(fileId, batchId, batchVer int64) (bool, *BatchDescr, error)
    GetSpecUnusedBatchDescr(fileId, batchId, batchVer int64) (bool, *BatchDescr, error)
    IncSpecBatchDescrUC(fileId, batchId, batchVer int64) error
    ListAllBatchDescrs() ([]*BatchDescr, error)
}

type IFileReg interface {
    GetNewFileId() (int64, error)
    AddNewFileDescr(descr *FileDescr) error
    GetNewestFileDescr(fileId int64) (bool, *FileDescr, error)
    DecSpecFileDescrUC(fileId, fileVer int64) error
    IncSpecFileDescrUC(fileId, fileVer int64) error
    GetSpecFileDescr(fileId, fileVer int64) (bool, *FileDescr, error)
    GetSpecUnusedFileDescr(fileId, fileVer int64) (bool, *FileDescr, error)
    EraseSpecFileDescr(fileId, fileVer int64) error
    GetAnyUnusedFileDescr() (bool, *FileDescr, error)
    ListAllFileDescrs() ([]*FileDescr, error)
    EraseAllFileDescrs() error

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
    EntryDescrExists(userId int64, dirPath, fileName string) (bool, error)
    GetEntryDescr(userId int64, dirPath, fileName string) (bool, *EntryDescr, error)
    ListEntryDescr(userId int64, dirPath string) ([]*EntryDescr, error)
    EraseEntryDescr(userId int64, dirPath, fileName string) error
    EraseEntryDescrsByUserId(userId int64) error
}
