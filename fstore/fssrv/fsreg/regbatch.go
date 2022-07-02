/*
 * Copyright 2022 Oleg Borodin  <borodin@unix7.org>
 */
package fsreg

import (
    "ndstore/dscom"
    "ndstore/dserr"
)

const batchSchema = `
    --- DROP TABLE IF EXISTS fs_batchs;
    CREATE TABLE IF NOT EXISTS fs_batchs (
        file_id         BIGINT,
        batch_id        BIGINT,

        batch_ver       BIGINT,
        u_counter       BIGINT,

        batch_size      BIGINT,
        block_size      BIGINT
    );
    --- DROP INDEX IF EXISTS fs_batch_idx;
    CREATE UNIQUE INDEX IF NOT EXISTS fs_batch_idx
        ON fs_batchs(file_id, batch_id, batch_ver);`


func (reg *Reg) AddNewBatchDescr(descr *dscom.BatchDescr) error {
    var err error
    request := `
        INSERT INTO fs_batchs(file_id, batch_id, batch_ver, u_counter, batch_size, block_size)
        VALUES ($1, $2, $3, $4, $5, $6);`
    _, err = reg.db.Exec(request, descr.FileId, descr.BatchId, descr.BatchVer, descr.UCounter,
                                                                descr.BatchSize, descr.BlockSize)
    if err != nil {
        return dserr.Err(err)
    }
    return dserr.Err(err)
}

func (reg *Reg) GetNewestBatchDescr(fileId, batchId int64) (bool, *dscom.BatchDescr, error) {
    var err error
    var exists bool
    var descr *dscom.BatchDescr
    request := `
        SELECT file_id, batch_id, batch_ver, u_counter, batch_size, block_size
        FROM fs_batchs
        WHERE file_id = $1
            AND batch_id = $2
            AND u_counter > 0
        ORDER BY batch_ver DESC
        LIMIT 1;`
    descrs := make([]*dscom.BatchDescr, 0)
    err = reg.db.Select(&descrs, request, fileId, batchId)
    if err != nil {
        return exists, descr, dserr.Err(err)
    }
    if len(descrs) > 0 {
        exists = true
        descr = descrs[0]
    }
    return exists, descr, dserr.Err(err)
}

func (reg *Reg) GetSpecBatchDescr(fileId, batchId, batchVer int64) (bool, *dscom.BatchDescr, error) {
    var err error
    var exists bool
    var descr *dscom.BatchDescr
    request := `
        SELECT file_id, batch_id, batch_ver, u_counter, batch_size, block_size
        FROM fs_batchs
        WHERE file_id = $1
            AND batch_id = $2
            AND batch_ver = $3
            AND u_counter > 0
        LIMIT 1;`
    descrs := make([]*dscom.BatchDescr, 0)
    err = reg.db.Select(&descrs, request, fileId, batchId, batchVer)
    if err != nil {
        return exists, descr, dserr.Err(err)
    }
    if len(descrs) > 0 {
        exists = true
        descr = descrs[0]
    }
    return exists, descr, dserr.Err(err)
}


func (reg *Reg) GetSpecUnusedBatchDescr(fileId, batchId, batchVer int64) (bool, *dscom.BatchDescr, error) {
    var err error
    var exists bool
    var descr *dscom.BatchDescr
    request := `
        SELECT file_id, batch_id, batch_ver, u_counter, batch_size, block_size
        FROM fs_batchs
        WHERE file_id = $1
            AND batch_id = $2
            AND batch_ver = $3
            AND u_counter < 1
        LIMIT 1;`
    descrs := make([]*dscom.BatchDescr, 0)
    err = reg.db.Select(&descrs, request, fileId, batchId, batchVer)
    if err != nil {
        return exists, descr, dserr.Err(err)
    }
    if len(descrs) > 0 {
        exists = true
        descr = descrs[0]
    }
    return exists, descr, dserr.Err(err)
}


func (reg *Reg) GetAnyUnusedBatchDescr() (bool, *dscom.BatchDescr, error) {
    var err     error
    var exists  bool
    var blockDescr *dscom.BatchDescr
    blocks := make([]*dscom.BatchDescr, 0)
    request := `
        SELECT file_id, batch_id, batch_ver, u_counter, batch_size, block_size
        FROM fs_batchs
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

func (reg *Reg) IncSpecBatchDescrUC(fileId, batchId, batchVer int64) error {
    var err error
    request := `
        UPDATE fs_batchs SET
            u_counter = u_counter + 1
        WHERE file_id = $1
            AND batch_id = $2
            AND batch_ver = $3;`
    _, err = reg.db.Exec(request, fileId, batchId, batchVer)
    if err != nil {
        return dserr.Err(err)
    }
    return dserr.Err(err)
}

func (reg *Reg) DecSpecBatchDescrUC(fileId, batchId, batchVer int64) error {
    var err error
    request := `
        UPDATE fs_batchs SET
            u_counter = u_counter - 1
        WHERE file_id = $1
            AND batch_id = $2
            AND batch_ver = $3
            AND u_counter > 0;`
    _, err = reg.db.Exec(request, fileId, batchId, batchVer)
    if err != nil {
        return dserr.Err(err)
    }
    return dserr.Err(err)
}

func (reg *Reg) ListAllBatchDescrs() ([]*dscom.BatchDescr, error) {
    var err error
    blocks := make([]*dscom.BatchDescr, 0)
    request := `
        SELECT file_id, batch_id, batch_ver, u_counter, batch_size, block_size
        FROM fs_batchs;`
    err = reg.db.Select(&blocks, request)
    if err != nil {
        return blocks, dserr.Err(err)
    }
    return blocks, dserr.Err(err)
}

func (reg *Reg) EraseSpecBatchDescr(fileId, batchId, batchVer int64) error {
    var err error
    request := `
        DELETE FROM fs_batchs
        WHERE file_id = $1
            AND batch_id = $2
            AND batch_ver = $3;`
    _, err = reg.db.Exec(request, fileId, batchId, batchVer)
    if err != nil {
        return dserr.Err(err)
    }
    return dserr.Err(err)
}

func (reg *Reg) EraseAllBatchDescrs() error {
    var err error
    request := `
        DELETE FROM fs_batchs;`
    _, err = reg.db.Exec(request)
    if err != nil {
        return dserr.Err(err)
    }
    return dserr.Err(err)
}
