package dscom

type UserDescr struct {
    UserId  int64       `json:"userId"      db:"user_id"`
    Login   string      `json:"login"       db:"login"`
    Pass    string      `json:"pass"        db:"pass"`
    State   string      `json:"state"       db:"state"`
    Role    string      `json:"role"        db:"role"`
}

func NewUserDescr() *UserDescr{
    return &UserDescr{}
}

type BStoreDescr struct {
    BStoreId int64      `json:"bStoreId"    db:"bstore_id"`
    Address string      `json:"address"     db:"address"`
    Port    string      `json:"port"        db:"port"`
    Login   string      `json:"login"       db:"login"`
    Pass    string      `json:"pass"        db:"pass"`
    State   string      `json:"state"       db:"state"`
}

func NewBStoreDescr() *BStoreDescr{
    return &BStoreDescr{}
}

type EntryDescr struct {
    EntryId     string      `json:"entryId"     db:"entryId"`
    DirPath     string      `json:"dirPath"     db:"dir_path"`
    FileName    string      `json:"fileName"    db:"file_name"`
    FileId      int64       `json:"fileId"      db:"file_id"`
}

func NewEntryDescr() *EntryDescr {
    var entry EntryDescr
    return &entry
}


type FileDescr struct {
    FileId      int64       `json:"fileId"      db:"file_id"`
    FileSize    int64       `json:"fileSize"    db:"file_size"`
    BatchSize   int64       `json:"batchSize"   db:"batch_size"`
    BlockSize   int64       `json:"blockSize"   db:"block_size"`
    BatchCount  int64       `json:"batchCount"  db:"batch_count"`
    Batchs      []*BatchDescr  `json:"batchs"   db:"-"`
}

func NewFileDescr() *FileDescr {
    var file FileDescr
    file.Batchs = make([]*BatchDescr, 0)
    return &file
}

type BatchDescr struct {
    FileId      int64       `json:"fileId"      db:"file_id"`
    BatchId     int64       `json:"batchId"     db:"batch_id"`
    BatchSize   int64       `json:"batchSize"   db:"batch_size"`
    BlockSize   int64       `json:"blockSize"   db:"block_size"`
    Blocks      []*BlockDescr  `json:"blocks"   db:"-"`
}

func NewBatchDescr() *BatchDescr {
    var batch BatchDescr
    batch.Blocks = make([]*BlockDescr, 0)
    return &batch
}

type BlockDescr struct {
    FileId      int64       `json:"fileId"      db:"file_id"`
    BatchId     int64       `json:"batchId"     db:"batch_id"`
    BlockId     int64       `json:"blockId"     db:"block_id"`
    BlockSize   int64       `json:"blockSize"   db:"block_size"`
    FilePath    string      `json:"filePath"    db:"file_path"`
    DataSize    int64       `json:"dataSize"    db:"data_size"`
    HashAlg     string      `json:"hashAlg"     db:"hash_alg"`
    HashSum     string      `json:"hashSum"     db:"hash_sum"`
    HashInit    string      `json:"hashInit"    db:"hash_init"`
}

func NewBlockDescr() *BlockDescr {
    var block BlockDescr
    return &block
}
