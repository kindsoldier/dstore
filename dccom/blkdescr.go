package dccom


type Block struct {
    ClusterId   int64     `db:"cluster_id"  json:"clusterId"`
    FileId      int64     `db:"file_id"     json:"fileId"`
    BatchId     int64     `db:"batch_id"    json:"batchId"`
    BlockId     int64     `db:"block_id"    json:"blockId"`
    BlockSize   int64     `db:"block_size"  json:"blockSize"`
    FileName    string    `db:"file_path"   json:"filePath"`
    HashAlg     string    `db:"hash_alg"    json:"hashAlg"`
    HashSum     string    `db:"hash_sum"    json:"hashSum"`
    HashInit    string    `db:"hash_init"   json:"hashInit"`
}
