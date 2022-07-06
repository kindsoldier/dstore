/*
 * Copyright 2022 Oleg Borodin  <borodin@unix7.org>
 */

package dscom

type IFileDistr interface {
    SaveBlock(descr *BlockDescr) (bool, int64, error)
}

type IFile interface {
}

type IBatch interface {
}

type IBlock interface {
}

type IFSReg interface {
    IBStoreReg
    IUserReg
    IEntryReg

    IFileReg
    IBatchReg
    IBlockReg
}

type IFileReg interface {
    GetAnyNotDistrFileDescr() (bool, *FileDescr, error)
    GetNewFileId() (int64, error)
    GetNewestFileDescr(fileId int64) (bool, *FileDescr, error)
    AddNewFileDescr(descr *FileDescr) error
    DecSpecFileDescrUC(count, fileId, fileVer int64) error
    IncSpecFileDescrUC(count, fileId, fileVer int64) error
    GetSpecFileDescr(fileId, fileVer int64) (bool, *FileDescr, error)
    GetSpecUnusedFileDescr(fileId, fileVer int64) (bool, *FileDescr, error)
    EraseSpecFileDescr(fileId, fileVer int64) error
    GetAnyUnusedFileDescr() (bool, *FileDescr, error)
    ListAllFileDescrs() ([]*FileDescr, error)
    EraseAllFileDescrs() error

    GetSetNotDistrFileDescr(count int) (bool, []*FileDescr, error)
}

type IBatchReg interface {
    AddNewBatchDescr(descr *BatchDescr) error
    IncSpecBatchDescrUC(count, fileId, batchId, batchVer int64) error
    DecSpecBatchDescrUC(count, fileId, batchId, batchVer int64) error
    GetNewestBatchDescr(fileId, batchId int64) (bool, *BatchDescr, error)
    GetSpecBatchDescr(fileId, batchId, batchVer int64) (bool, *BatchDescr, error)
    GetSpecUnusedBatchDescr(fileId, batchId, batchVer int64) (bool, *BatchDescr, error)
    GetAnyUnusedBatchDescr() (bool, *BatchDescr, error)
    ListAllBatchDescrs() ([]*BatchDescr, error)
    EraseSpecBatchDescr(fileId, batchId, batchVer int64) error
    EraseAllBatchDescrs() error
}

type IBlockReg interface {
    AddNewBlockDescr(descr *BlockDescr) error
    IncSpecBlockDescrUC(count, fileId, batchId, blockId int64, blockType string, blockVer int64) error
    DecSpecBlockDescrUC(count, fileId, batchId, blockId int64, blockType string, blockVer int64) error
    EraseSpecBlockDescr(fileId, batchId, blockId int64, blockType string, blockVer int64) error
    GetAnyUnusedBlockDescr() (bool, *BlockDescr, error)
    GetNewestBlockDescr(fileId, batchId, blockId int64, blockType string) (bool, *BlockDescr, error)
    GetSpecBlockDescr(fileId, batchId, blockId int64, blockType string, blockVer int64) (bool, *BlockDescr, error)
    GetSpecUnusedBlockDescr(fileId, batchId, blockId int64, blockType string, blockVer int64) (bool, *BlockDescr, error)
    ListAllBlockDescrs() ([]*BlockDescr, error)

    GetBStoreDescrById(bstoreId int64) (bool, *BStoreDescr, error)
    GetSetUnusedBlockDescrs(count int) (bool, []*BlockDescr, error)
}


type IBStoreReg interface {
    AddBStoreDescr(descr *BStoreDescr) (int64, error)
    GetBStoreDescr(address, port string) (bool, *BStoreDescr, error)
    ListBStoreDescrs() ([]*BStoreDescr, error)
    UpdateBStoreDescr(descr *BStoreDescr) error
    EraseBStoreDescr(address, port string) error
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
