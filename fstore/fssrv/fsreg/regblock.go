/*
 * Copyright 2022 Oleg Borodin  <borodin@unix7.org>
 */
package fsreg

import (
    "ndstore/dscom"
    "ndstore/dserr"
)

const blockSchema = `
    DROP TABLE IF EXISTS fs_blocks;
    CREATE TABLE IF NOT EXISTS fs_blocks (
        file_id         INTEGER,
        batch_id        INTEGER,
        block_id        INTEGER,
        block_type      TEXT DEFAULT '',
        block_size      INTEGER,
        data_size       INTEGER,
        file_path       TEXT DEFAULT '',
        hash_alg        TEXT DEFAULT '',
        hash_sum        TEXT DEFAULT '',
        hash_init       TEXT DEFAULT '',
        bstore_id       INTEGER,
        saved_rem       BOOL,
        fstore_id       INTEGER,
        saved_loc       BOOL,
        loc_updated     BOOL
    );
    DROP INDEX IF EXISTS fs_block_idx;
    CREATE UNIQUE INDEX IF NOT EXISTS fs_block_idx
        ON fs_blocks(file_id, batch_id, block_id, block_type);`


func (reg *Reg) AddBlockDescr(descr *dscom.BlockDescr) error {
    var err error
    request := `
        INSERT INTO fs_blocks(file_id, batch_id, block_id, block_size, data_size,
                                    file_path, block_type, hash_alg, hash_init, hash_sum,
                                    fstore_id, bstore_id, saved_loc, saved_rem, loc_updated)
        VALUES ($1, $2, $3, $4, $5, $6, $7, $8, $9, $10, $11, $12, $13, $14, $15);`
    _, err = reg.db.Exec(request, descr.FileId, descr.BatchId, descr.BlockId, descr.BlockSize, descr.DataSize,
                                    descr.FilePath, descr.BlockType, descr.HashAlg, descr.HashInit, descr.HashSum,
                                    descr.FStoreId, descr.BStoreId, descr.SavedLoc, descr.SavedRem, descr.LocUpdated)
    if err != nil {
        return dserr.Err(err)
    }
    return dserr.Err(err)
}


func (reg *Reg) UpdateBlockDescr(descr *dscom.BlockDescr) error {
    var err error
    request := `
        UPDATE fs_blocks SET block_size = $1, data_size = $2,
                                    file_path = $3, hash_alg = $4, hash_init = $5, hash_sum = $6,
                                    fstore_id = $7, bstore_id = $8, saved_loc = $9, saved_rem = $10,
                                    loc_updated = $11
        WHERE file_id = $12
            AND batch_id = $13
            AND block_id = $14
            AND block_type = $15;`
    _, err = reg.db.Exec(request, descr.BlockSize, descr.DataSize,
                                    descr.FilePath, descr.HashAlg, descr.HashInit, descr.HashSum,
                                    descr.FStoreId, descr.BStoreId, descr.SavedLoc, descr.SavedRem,
                                    descr.LocUpdated,
                                    descr.FileId, descr.BatchId, descr.BlockId, descr.BlockType)
    if err != nil {
        return dserr.Err(err)
    }
    return dserr.Err(err)
}


func (reg *Reg) GetBlockDescr(fileId, batchId, blockId int64, blockType string) (bool, *dscom.BlockDescr, error) {
    var err error
    exists := false
    var blockDescr *dscom.BlockDescr

    blockDescrs := make([]*dscom.BlockDescr, 0)
    request := `
        SELECT file_id, batch_id, block_id, block_size, data_size,
                                    file_path, block_type, hash_alg, hash_init, hash_sum,
                                    fstore_id, bstore_id, saved_loc, saved_rem, loc_updated
        FROM fs_blocks
        WHERE file_id = $1
            AND batch_id = $2
            AND block_id = $3
            AND block_type = $4
        LIMIT 1;`
    err = reg.db.Select(&blockDescrs, request, fileId, batchId, blockId, blockType )
    if err != nil {
        return exists, blockDescr, dserr.Err(err)
    }
    if len(blockDescrs) > 0 {
        exists = true
        blockDescr = blockDescrs[0]
    }
    return exists, blockDescr, dserr.Err(err)
}


func (reg *Reg) ListBlockDescrs() ([]*dscom.BlockDescr, error) {
    var err error
    blocks := make([]*dscom.BlockDescr, 0)
    request := `
        SELECT file_id, batch_id, block_id, block_size, data_size,
                                    file_path, block_type, hash_alg, hash_init, hash_sum,
                                    fstore_id, bstore_id, saved_loc, saved_rem, loc_updated
        FROM fs_blocks;`
    err = reg.db.Select(&blocks, request)
    if err != nil {
        return blocks, dserr.Err(err)
    }
    return blocks, dserr.Err(err)
}

func (reg *Reg) EraseBlockDescr(fileId, batchId, blockId int64, blockType string) error {
    var err error
    request := `
        DELETE FROM fs_blocks
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
