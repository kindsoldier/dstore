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
        block_type      TEXT DEFAULT '',
        block_ver       INTEGER,
        u_counter       INTEGER,
        block_size      INTEGER,
        data_size       INTEGER,
        file_path       TEXT DEFAULT '',
        hash_alg        TEXT DEFAULT '',
        hash_sum        TEXT DEFAULT '',
        hash_init       TEXT DEFAULT ''
    );
    --- DROP INDEX IF EXISTS block_idx;
    CREATE UNIQUE INDEX IF NOT EXISTS block_idx
        ON blocks (file_id, batch_id, block_id, block_type, block_ver);`


func (reg *Reg) AddNewBlockDescr(descr *dscom.BlockDescr) error {
    var err error
    request := `
        INSERT INTO blocks(file_id, batch_id, block_id, u_counter, block_size, data_size,
                                file_path, block_type, hash_alg, hash_init, hash_sum, block_ver)
        VALUES ($1, $2, $3, $4, $5, $6, $7, $8, $9, $10, $11, $12);`
    _, err = reg.db.Exec(request, descr.FileId, descr.BatchId, descr.BlockId, descr.UCounter, descr.BlockSize, descr.DataSize,
                                descr.FilePath, descr.BlockType, descr.HashAlg, descr.HashInit, descr.HashSum, descr.BlockVer)
    if err != nil {
        return dserr.Err(err)
    }
    return dserr.Err(err)
}

func (reg *Reg) GetNewestBlockDescr(fileId, batchId, blockId int64, blockType string) (bool, *dscom.BlockDescr, error) {
    var err error
    var exists bool
    var descr *dscom.BlockDescr
    request := `
        SELECT file_id, batch_id, block_id, u_counter, block_size, data_size,
                file_path, block_type, hash_alg, hash_init, hash_sum, block_ver
        FROM blocks
        WHERE file_id = $1
            AND batch_id = $2
            AND block_id = $3
            AND block_type = $4
            AND u_counter > 0
        ORDER BY block_ver DESC
        LIMIT 1;`
    descrs := make([]*dscom.BlockDescr, 0)
    err = reg.db.Select(&descrs, request, fileId, batchId, blockId, blockType)
    if err != nil {
        return exists, descr, dserr.Err(err)
    }
    if len(descrs) > 0 {
        exists = true
        descr = descrs[0]
    }
    return exists, descr, dserr.Err(err)
}

func (reg *Reg) GetSpecBlockDescr(fileId, batchId, blockId int64, blockType string, blockVer int64) (bool, *dscom.BlockDescr, error) {
    var err error
    var exists bool
    var descr *dscom.BlockDescr
    request := `
        SELECT file_id, batch_id, block_id, u_counter, block_size, data_size,
                file_path, block_type, hash_alg, hash_init, hash_sum, block_ver
        FROM blocks
        WHERE file_id = $1
            AND batch_id = $2
            AND block_id = $3
            AND block_type = $4
            AND block_ver = $5
            AND u_counter > 0
        LIMIT 1;`
    descrs := make([]*dscom.BlockDescr, 0)
    err = reg.db.Select(&descrs, request, fileId, batchId, blockId, blockType, blockVer)
    if err != nil {
        return exists, descr, dserr.Err(err)
    }
    if len(descrs) > 0 {
        exists = true
        descr = descrs[0]
    }
    return exists, descr, dserr.Err(err)
}


func (reg *Reg) GetSpecUnusedBlockDescr(fileId, batchId, blockId int64, blockType string, blockVer int64) (bool, *dscom.BlockDescr, error) {
    var err error
    var exists bool
    var descr *dscom.BlockDescr
    request := `
        SELECT file_id, batch_id, block_id, u_counter, block_size, data_size,
                file_path, block_type, hash_alg, hash_init, hash_sum, block_ver
        FROM blocks
        WHERE file_id = $1
            AND batch_id = $2
            AND block_id = $3
            AND block_type = $4
            AND block_ver = $5
            AND u_counter < 1
        ORDER BY block_ver DESC
        LIMIT 1;`
    descrs := make([]*dscom.BlockDescr, 0)
    err = reg.db.Select(&descrs, request, fileId, batchId, blockId, blockType, blockVer)
    if err != nil {
        return exists, descr, dserr.Err(err)
    }
    if len(descrs) > 0 {
        exists = true
        descr = descrs[0]
    }
    return exists, descr, dserr.Err(err)
}


func (reg *Reg) GetAnyUnusedBlockDescr() (bool, *dscom.BlockDescr, error) {
    var err     error
    var exists  bool
    var blockDescr *dscom.BlockDescr
    blocks := make([]*dscom.BlockDescr, 0)
    request := `
        SELECT file_id, batch_id, block_id, u_counter, block_size, data_size,
                file_path, block_type, hash_alg, hash_init, hash_sum, block_ver
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

func (reg *Reg) IncSpecBlockDescrUC(fileId, batchId, blockId int64, blockType string, blockVer int64) error {
    var err error
    request := `
        UPDATE blocks SET
            u_counter = u_counter + 1
        WHERE file_id = $1
            AND batch_id = $2
            AND block_id = $3
            AND block_type = $4
            AND block_ver = $5
            ;`
    _, err = reg.db.Exec(request, fileId, batchId, blockId, blockType, blockVer)
    if err != nil {
        return dserr.Err(err)
    }
    return dserr.Err(err)
}

func (reg *Reg) DecSpecBlockDescrUC(fileId, batchId, blockId int64, blockType string, blockVer int64) error {
    var err error
    request := `
        UPDATE blocks SET
            u_counter = u_counter - 1
        WHERE file_id = $1
            AND batch_id = $2
            AND block_id = $3
            AND block_type = $4
            AND block_ver = $5
            AND u_counter > 0;`
    _, err = reg.db.Exec(request, fileId, batchId, blockId, blockType, blockVer)
    if err != nil {
        return dserr.Err(err)
    }
    return dserr.Err(err)
}

func (reg *Reg) ListAllBlockDescrs() ([]*dscom.BlockDescr, error) {
    var err error
    blocks := make([]*dscom.BlockDescr, 0)
    request := `
        SELECT file_id, batch_id, block_id, block_ver, u_counter, block_size, data_size,
                                    file_path, hash_alg, hash_sum, hash_init, block_type
        FROM blocks;`
    err = reg.db.Select(&blocks, request)
    if err != nil {
        return blocks, dserr.Err(err)
    }
    return blocks, dserr.Err(err)
}

func (reg *Reg) EraseSpecBlockDescr(fileId, batchId, blockId int64, blockType string, blockVer int64) error {
    var err error
    request := `
        DELETE FROM blocks
        WHERE file_id = $1
            AND batch_id = $2
            AND block_id = $3
            AND block_type = $4
            AND block_ver = $5
            ;`
    _, err = reg.db.Exec(request, fileId, batchId, blockId, blockType, blockVer)
    if err != nil {
        return dserr.Err(err)
    }
    return dserr.Err(err)
}

func (reg *Reg) EraseAllDescrs() error {
    var err error
    request := `
        DELETE FROM blocks;`
    _, err = reg.db.Exec(request)
    if err != nil {
        return dserr.Err(err)
    }
    return dserr.Err(err)
}
