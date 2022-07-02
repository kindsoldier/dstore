/*
 * Copyright 2022 Oleg Borodin  <borodin@unix7.org>
 */

package dscom

type UserDescr struct {
    UserId      int64       `json:"userId"      db:"user_id"`
    Login       string      `json:"login"       db:"login"`
    Pass        string      `json:"pass"        db:"pass"`
    State       string      `json:"state"       db:"state"`
    Role        string      `json:"role"        db:"role"`
    UpdatedAt   UnixTime    `json:"updatedAt"   db:"updated_at"`
    CreatedAt   UnixTime    `json:"createdAt"   db:"created_at"`
}

func NewUserDescr() *UserDescr{
    return &UserDescr{}
}

const BStateNormal      string = "normal"
const BStateDisabled    string = "disabled"
const BStateWrong       string = "wrong"

type BStoreDescr struct {
    BStoreId    int64       `json:"bStoreId"    db:"bstore_id"`
    Address     string      `json:"address"     db:"address"`
    Port        string      `json:"port"        db:"port"`
    Login       string      `json:"login"       db:"login"`
    Pass        string      `json:"pass"        db:"pass"`
    State       string      `json:"state"       db:"state"`
    UpdatedAt   UnixTime    `json:"updatedAt"   db:"updated_at"`
    CreatedAt   UnixTime    `json:"createdAt"   db:"created_at"`
}

func NewBStoreDescr() *BStoreDescr{
    return &BStoreDescr{}
}

type EntryDescr struct {
    EntryId     string      `json:"entryId"     db:"entry_id"`
    UserId      int64       `json:"userId"      db:"user_id"`
    DirPath     string      `json:"dirPath"     db:"dir_path"`
    FileName    string      `json:"fileName"    db:"file_name"`
    FileId      int64       `json:"fileId"      db:"file_id"`
    // From userDescr
    UserName    string      `json:"userName"    db:"user_name"`
    // From fileDescr
    FileSize    int64       `json:"fileSize"    db:"file_size"`
    UpdatedAt   UnixTime    `json:"updatedAt"   db:"updated_at"`
    CreatedAt   UnixTime    `json:"createdAt"   db:"created_at"`
}

func NewEntryDescr() *EntryDescr {
    var entry EntryDescr
    return &entry
}


type FileDescr struct {
    FileId      int64       `json:"fileId"      db:"file_id"`
    BatchSize   int64       `json:"batchSize"   db:"batch_size"`
    BlockSize   int64       `json:"blockSize"   db:"block_size"`
    BatchCount  int64       `json:"batchCount"  db:"batch_count"`

    FileVer     int64       `json:"fileVer"     db:"file_ver"`
    UCounter    int64       `json:"uCounter"    db:"u_counter"`

    FileSize    int64       `json:"fileSize"    db:"file_size"`
    UpdatedAt   UnixTime    `json:"updatedAt"   db:"updated_at"`
    CreatedAt   UnixTime    `json:"createdAt"   db:"created_at"`
}

func NewFileDescr() *FileDescr {
    var file FileDescr
    return &file
}

type BatchDescr struct {
    FileId      int64       `json:"fileId"      db:"file_id"`
    BatchId     int64       `json:"batchId"     db:"batch_id"`
    BatchSize   int64       `json:"batchSize"   db:"batch_size"`
    BlockSize   int64       `json:"blockSize"   db:"block_size"`

    BatchVer    int64       `json:"batchVer"    db:"batch_ver"`
    UCounter    int64       `json:"uCounter"    db:"u_counter"`
}

func NewBatchDescr() *BatchDescr {
    var batch BatchDescr
    return &batch
}

const BTypeData     string = "data"
const BTypeRecov    string = "reco"

const HashTypeHW    string = "hway"

type BlockDescr struct {
    FileId      int64       `json:"fileId"      db:"file_id"`
    BatchId     int64       `json:"batchId"     db:"batch_id"`
    BlockId     int64       `json:"blockId"     db:"block_id"`
    BlockType   string      `json:"blockType"   db:"block_type"`

    BlockVer    int64       `json:"blockVer"    db:"block_ver"`
    UCounter    int64       `json:"uCounter"    db:"u_counter"`

    BlockSize   int64       `json:"blockSize"   db:"block_size"`
    DataSize    int64       `json:"dataSize"    db:"data_size"`
    FilePath    string      `json:"filePath"    db:"file_path"`
    HashAlg     string      `json:"hashAlg"     db:"hash_alg"`
    HashSum     string      `json:"hashSum"     db:"hash_sum"`
    HashInit    string      `json:"hashInit"    db:"hash_init"`

    SavedLoc    bool        `json:"savedLoc"    db:"saved_loc"`
    SavedRem    bool        `json:"savedRem"    db:"saved_rem"`
    LocUpdated  bool        `json:"locUpdated"  db:"loc_updated"`

    FStoreId    int64       `json:"fstoreId"    db:"fstore_id"`
    BStoreId    int64       `json:"bstoreId"    db:"bstore_id"`

}

func NewBlockDescr() *BlockDescr {
    var block BlockDescr
    return &block
}
