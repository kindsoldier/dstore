package dscom


type FileMI struct {
    FileId      int64       `json:"fileId"`
    BatchSize   int64       `json:"batchSize"`
    BlockSize   int64       `json:"blockSize"`
    BatchCount  int64       `json:"batchCount"`
    Batchs      []*BatchMI  `json:"batchs"`
}

func NewFileMI() *FileMI {
    var fileMeta FileMI
    fileMeta.Batchs = make([]*BatchMI, 0)
    return &fileMeta
}


type BatchMI struct {
    Blocks      []*BlockMI  `json:"blocks"`
}

func NewBatchMI() *BatchMI {
    var batchMeta BatchMI
    batchMeta.Blocks = make([]*BlockMI, 0)
    return &batchMeta
}


type BlockMI struct {
    ClusterId   int64       `db:"cluster_id"  json:"clusterId"`
    FileId      int64       `db:"file_id"     json:"fileId"`
    BatchId     int64       `db:"batch_id"    json:"batchId"`
    BlockId     int64       `db:"block_id"    json:"blockId"`
    BlockSize   int64       `db:"block_size"  json:"blockSize"`
    FileName    string      `db:"file_path"   json:"filePath"`
    HashAlg     string      `db:"hash_alg"    json:"hashAlg"`
    HashSum     string      `db:"hash_sum"    json:"hashSum"`
    HashInit    string      `db:"hash_init"   json:"hashInit"`
}

//type BlockMI struct {
//    FileId      int64       `json:"fileId"      db:"file_id"`
//    FileName    string      `json:"fileName"    db:"file_name"`
//    HashSum     string      `json:"hashSum"     db:"hash_sum"`
//    HashInit    string      `json:"hashInit"    db:"hash_init"`
//    Size        int64       `json:"size"        db:"size"`
//}

func NewBlockMI() *BlockMI {
    var blockMeta BlockMI
    return &blockMeta
}
