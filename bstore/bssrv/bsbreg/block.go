/*
 * Copyright 2022 Oleg Borodin  <borodin@unix7.org>
 */
package bsbreg

import (
    "ndstore/dscom"
)


const blockSchema = `
    DROP TABLE IF EXISTS blocks;
    CREATE TABLE IF NOT EXISTS blocks (
        file_id     INTEGER,
        batch_id    INTEGER,
        block_id    INTEGER,
        block_size  INTEGER,
        data_size   INTEGER,
        file_path   TEXT DEFAULT '',
        hash_alg    TEXT DEFAULT '',
        hash_sum    TEXT DEFAULT '',
        hash_init   TEXT DEFAULT ''
    );
    DROP INDEX IF EXISTS block_idx;
    CREATE UNIQUE INDEX IF NOT EXISTS block_idx
        ON blocks (file_id, batch_id, block_id);`


func (reg *Reg) AddBlockDescr(fileId, batchId, blockId, blockSize, dataSize int64, filePath string) error {
    var err error
    if reg.db == nil {
        return ErrorNilRef
    }
    request := `
        INSERT
            INTO blocks(file_id, batch_id, block_id, block_size, data_size, file_path)
        VALUES ($1, $2, $3, $4, $5, $6);`
    _, err = reg.db.Exec(request, fileId, batchId, blockId, blockSize, dataSize, filePath)
    if err != nil {
        return err
    }
    return err
}


func (reg *Reg) UpdateBlockDescr(fileId, batchId, blockId, blockSize, dataSize int64, filePath string) error {
    var err error
    if reg.db == nil {
        return ErrorNilRef
    }
    var request string
    request = `
        UPDATE blocks SET
            block_size = $1,
            data_size = $3,
            file_path = $2
        WHERE file_id = $4
            AND batch_id = $5
            AND block_id = $6;`
    _, err = reg.db.Exec(request, blockSize, dataSize, filePath, fileId, batchId, blockId)
    if err != nil {
        return err
    }
    return err
}

func (reg *Reg) GetBlockFilePath(fileId, batchId, blockId int64) (string, int64, error) {
    var err error
    var filePath string
    var blockSize int64

    if reg.db == nil {
        return filePath, blockSize, ErrorNilRef
    }
    request := `
        SELECT file_path, block_size
        FROM blocks
        WHERE file_id = $1
            AND batch_id = $2
            AND block_id = $3
        LIMIT 1;`

    var block dscom.BlockDescr
    err = reg.db.Get(&block, request, fileId, batchId, blockId)
    if err != nil {
        return filePath, blockSize, err
    }
    filePath    = block.FilePath
    blockSize   = block.BlockSize
    return filePath, blockSize, err
}


func (reg *Reg) BlockDescrExists(fileId, batchId, blockId int64) (bool, error) {
    var err error
    var exists bool
    if reg.db == nil {
        return exists, ErrorNilRef
    }

    request := `
        SELECT file_path
        FROM blocks
        WHERE file_id = $1
            AND batch_id = $2
            AND block_id = $3
        LIMIT 1;`

    blocks := make([]dscom.BlockDescr, 0)
    err = reg.db.Select(&blocks, request, fileId, batchId, blockId)
    if err != nil {
        return exists, err
    }
    if len(blocks) > 0 {
        exists = true
    }
    return exists, err
}

func (reg *Reg) ListBlockDescrs() ([]*dscom.BlockDescr, error) {
    var err error
    blocks := make([]*dscom.BlockDescr, 0)
    if reg.db == nil {
        return blocks, ErrorNilRef
    }
    request := `
        SELECT file_id, batch_id, block_id, block_size, data_size, file_path,
            hash_alg, hash_sum, hash_init
        FROM blocks;`
    err = reg.db.Select(&blocks, request)
    if err != nil {
        return blocks, err
    }
    return blocks, err
}

func (reg *Reg) DeleteBlockDescr(fileId, batchId, blockId int64) error {
    var err error
    if reg.db == nil {
        return ErrorNilRef
    }
    request := `
        DELETE FROM blocks
        WHERE file_id = $1
            AND batch_id = $2
            AND block_id = $3;`
    _, err = reg.db.Exec(request, fileId, batchId, blockId)
    if err != nil {
        return err
    }
    return err
}

func (reg *Reg) xPurgeFile(fileId int64) error {
    var err error
    if reg.db == nil {
        return ErrorNilRef
    }
    request := `
        DELETE FROM blocks
        WHERE file_id = $1;`
    _, err = reg.db.Exec(request, fileId)
    if err != nil {
        return err
    }
    return err
}

func (reg *Reg) xPurgeCluster(userId int64) error {
    var err error
    if reg.db == nil {
        return ErrorNilRef
    }
    request := `
        DELETE FROM blocks;`
    _, err = reg.db.Exec(request, userId)
    if err != nil {
        return err
    }
    return err
}
