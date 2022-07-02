/*
 * Copyright 2022 Oleg Borodin  <borodin@unix7.org>
 */
package fsreg

import (
    "ndstore/dscom"
    "ndstore/dserr"
)

const blockSchema = `
    --- DROP TABLE IF EXISTS fs_blocks;
    CREATE TABLE IF NOT EXISTS fs_blocks (
        file_id         BIGINT,
        batch_id        BIGINT,
        block_id        BIGINT,
        block_type      TEXT DEFAULT '',

        block_ver       BIGINT,
        u_counter       BIGINT,

        block_size      BIGINT,
        data_size       BIGINT,
        file_path       TEXT DEFAULT '',

        hash_alg        TEXT DEFAULT '',
        hash_init       TEXT DEFAULT '',
        hash_sum        TEXT DEFAULT '',

        saved_loc       BOOL,
        saved_rem       BOOL,
        fstore_id       BIGINT,
        bstore_id       BIGINT,
        loc_updated     BOOL
    );
    --- DROP INDEX IF EXISTS fs_block_idx;
    CREATE UNIQUE INDEX IF NOT EXISTS fs_block_idx
        ON fs_blocks (file_id, batch_id, block_id, block_type, block_ver);`



func (reg *Reg) AddNewBlockDescr(descr *dscom.BlockDescr) error {
    var err error
    request := `
        INSERT INTO fs_blocks(file_id, batch_id, block_id, block_type, block_ver, u_counter,
                        block_size, data_size, file_path, hash_alg, hash_init, hash_sum,
                        saved_loc, saved_rem, fstore_id, bstore_id, loc_updated)
        VALUES ($1, $2, $3, $4, $5, $6, $7, $8, $9, $10, $11, $12, $13, $14, $15, $16, $17);`
    _, err = reg.db.Exec(request, descr.FileId, descr.BatchId, descr.BlockId, descr.BlockType,
                                                                    descr.BlockVer, descr.UCounter,
                                descr.BlockSize, descr.DataSize, descr.FilePath, descr.HashAlg,
                                                                    descr.HashInit, descr.HashSum,
                                descr.SavedLoc, descr.SavedRem, descr.FStoreId, descr.BStoreId,
                                                                                descr.LocUpdated)
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
        SELECT file_id, batch_id, block_id, block_type, block_ver, u_counter,
                        block_size, data_size, file_path, hash_alg, hash_init, hash_sum,
                        saved_loc, saved_rem, fstore_id, bstore_id, loc_updated
        FROM fs_blocks
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
        SELECT file_id, batch_id, block_id, block_type, block_ver, u_counter,
                        block_size, data_size, file_path, hash_alg, hash_init, hash_sum,
                        saved_loc, saved_rem, fstore_id, bstore_id, loc_updated
        FROM fs_blocks
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
        SELECT file_id, batch_id, block_id, block_type, block_ver, u_counter,
                        block_size, data_size, file_path, hash_alg, hash_init, hash_sum,
                        saved_loc, saved_rem, fstore_id, bstore_id, loc_updated
        FROM fs_blocks
        WHERE file_id = $1
            AND batch_id = $2
            AND block_id = $3
            AND block_type = $4
            AND block_ver = $5
            AND u_counter < 1
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
        SELECT file_id, batch_id, block_id, block_type, block_ver, u_counter,
                        block_size, data_size, file_path, hash_alg, hash_init, hash_sum,
                        saved_loc, saved_rem, fstore_id, bstore_id, loc_updated
        FROM fs_blocks
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
        UPDATE fs_blocks SET
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
        UPDATE fs_blocks SET
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
        SELECT file_id, batch_id, block_id, block_type, block_ver, u_counter,
                        block_size, data_size, file_path, hash_alg, hash_init, hash_sum,
                        saved_loc, saved_rem, fstore_id, bstore_id, loc_updated
        FROM fs_blocks;`
    err = reg.db.Select(&blocks, request)
    if err != nil {
        return blocks, dserr.Err(err)
    }
    return blocks, dserr.Err(err)
}

func (reg *Reg) EraseSpecBlockDescr(fileId, batchId, blockId int64, blockType string, blockVer int64) error {
    var err error
    request := `
        DELETE FROM fs_blocks
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

func (reg *Reg) EraseAllBlockDescrs() error {
    var err error
    request := `
        DELETE FROM fs_blocks;`
    _, err = reg.db.Exec(request)
    if err != nil {
        return dserr.Err(err)
    }
    return dserr.Err(err)
}
