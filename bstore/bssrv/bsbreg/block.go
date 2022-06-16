/*
 * Copyright 2022 Oleg Borodin  <borodin@unix7.org>
 */
package bsbreg

import (
    "ndstore/dscom"
    "ndstore/dserr"
)

const blockSchema = `
    --- DROP TABLE IF EXISTS blocks;
    CREATE TABLE IF NOT EXISTS blocks (
        file_id         INTEGER,
        batch_id        INTEGER,
        block_id        INTEGER,
        u_counter       INTEGER,
        block_size      INTEGER,
        data_size       INTEGER,
        block_type      TEXT DEFAULT '',
        file_path       TEXT DEFAULT '',
        hash_alg        TEXT DEFAULT '',
        hash_sum        TEXT DEFAULT '',
        hash_init       TEXT DEFAULT ''
    );
    --- DROP INDEX IF EXISTS block_idx;
    CREATE UNIQUE INDEX IF NOT EXISTS block_idx
        ON blocks (file_id, batch_id, block_id, block_type);`


func (reg *Reg) AddBlockDescr(fileId, batchId, blockId, uCounter, blockSize, dataSize int64,
                                            filePath, blockType, hashAlg, hashInit, hashSum string) error {
    var err error
    request := `
        INSERT
            INTO blocks(file_id, batch_id, block_id, u_counter, block_size, data_size,
                                                    file_path, block_type, hash_alg, hash_init, hash_sum)
        VALUES ($1, $2, $3, $4, $5, $6, $7, $8, $9, $10, $11);`
    _, err = reg.db.Exec(request, fileId, batchId, blockId, uCounter, blockSize, dataSize,
                                                    filePath, blockType, hashAlg, hashInit, hashSum)
    if err != nil {
        return dserr.Err(err)
    }
    return dserr.Err(err)
}

func (reg *Reg) GetBlockParams(fileId, batchId, blockId int64, blockType string) (bool, bool, string, int64, error) {
    var err error
    var exists bool
    var used bool
    var filePath string
    var dataSize int64
    request := `
        SELECT file_path, data_size, u_counter
        FROM blocks
        WHERE file_id = $1
            AND batch_id = $2
            AND block_id = $3
            AND block_type = $4
        LIMIT 1;`
    blocks := make([]*dscom.BlockDescr, 0)
    err = reg.db.Select(&blocks, request, fileId, batchId, blockId, blockType)
    if err != nil {
        return exists, used, filePath, dataSize, dserr.Err(err)
    }
    if len(blocks) > 0 {
        exists   = true
        filePath = blocks[0].FilePath
        dataSize = blocks[0].DataSize
        if blocks[0].UCounter > 0 {
            used = true
        }
    }
    return exists, used, filePath, dataSize, dserr.Err(err)
}

func (reg *Reg) GetUnusedBlockDescr() (bool, *dscom.BlockDescr, error) {
    var err     error
    var exists  bool
    var blockDescr *dscom.BlockDescr
    blocks := make([]*dscom.BlockDescr, 0)
    request := `
        SELECT file_id, batch_id, block_id, u_counter, block_size, data_size,
                                    file_path, hash_alg, hash_sum, hash_init, block_type
        FROM blocks
        WHERE u_counter < 1
        LIMIT 1;`
    err = reg.db.Select(&blocks, request)
    if err != nil {
        return exists, blockDescr, dserr.Err(err)
    }
    if len(blocks) > 0 {
        exists = true
        blockDescr = blocks[0]
    }
    return exists, blockDescr, dserr.Err(err)
}

func (reg *Reg) IncBlockDescrUC(fileId, batchId, blockId int64, blockType string) error {
    var err error
    request := `
        UPDATE blocks SET
            u_counter = u_counter + 1
        WHERE file_id = $1
            AND batch_id = $2
            AND block_id = $3
            AND block_type = $4;`
    _, err = reg.db.Exec(request, fileId, batchId, blockId, blockType)
    if err != nil {
        return dserr.Err(err)
    }
    return dserr.Err(err)
}

func (reg *Reg) DecBlockDescrUC(fileId, batchId, blockId int64, blockType string) error {
    var err error
    request := `
        UPDATE blocks SET
            u_counter = u_counter - 1
        WHERE file_id = $1
            AND batch_id = $2
            AND block_id = $3
            AND block_type = $4
            AND u_counter > 0;`
    _, err = reg.db.Exec(request, fileId, batchId, blockId, blockType)
    if err != nil {
        return dserr.Err(err)
    }
    return dserr.Err(err)
}

func (reg *Reg) ListBlockDescrs() ([]*dscom.BlockDescr, error) {
    var err error
    blocks := make([]*dscom.BlockDescr, 0)
    request := `
        SELECT file_id, batch_id, block_id, u_counter, block_size, data_size,
                                    file_path, hash_alg, hash_sum, hash_init, block_type
        FROM blocks;`
    err = reg.db.Select(&blocks, request)
    if err != nil {
        return blocks, dserr.Err(err)
    }
    return blocks, dserr.Err(err)
}

func (reg *Reg) DropBlockDescr(fileId, batchId, blockId int64, blockType string) error {
    var err error
    request := `
        DELETE FROM blocks
        WHERE file_id = $1
            AND batch_id = $2
            AND block_id = $3
            AND block_type = $4;`
    _, err = reg.db.Exec(request, fileId, batchId, blockId, blockType)
    if err != nil {
        return dserr.Err(err)
    }
    return dserr.Err(err)
}

func (reg *Reg) PurgeAllDescrs() error {
    var err error
    request := `
        DELETE FROM blocks;`
    _, err = reg.db.Exec(request)
    if err != nil {
        return dserr.Err(err)
    }
    return dserr.Err(err)
}
