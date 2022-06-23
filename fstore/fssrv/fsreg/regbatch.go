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
        file_id         INTEGER,
        batch_id        INTEGER,
        batch_size      INTEGER,
        block_size      INTEGER
    );
    --- DROP INDEX IF EXISTS fs_batch_idx;
    CREATE UNIQUE INDEX IF NOT EXISTS fs_batch_idx
        ON fs_batchs(file_id, batch_id);`


func (reg *Reg) AddBatchDescr(descr *dscom.BatchDescr) error {
    var err error
    request := `
        INSERT INTO fs_batchs(file_id, batch_id, batch_size, block_size)
        VALUES ($1, $2, $3, $4);`
    _, err = reg.db.Exec(request, descr.FileId, descr.BatchId, descr.BatchSize, descr.BlockSize)
    if err != nil {
        return dserr.Err(err)
    }
    return dserr.Err(err)
}

func (reg *Reg) UpdateBatchDescr(descr *dscom.BatchDescr) error {
    var err error
    request := `
        UPDATE fs_batchs SET batch_size = $1, block_size = $2
        WHERE file_id = $3, batch_id = $3;`
    _, err = reg.db.Exec(request, descr.BatchSize, descr.BlockSize, descr.FileId, descr.BatchId)
    if err != nil {
        return dserr.Err(err)
    }
    return dserr.Err(err)
}

func (reg *Reg) GetBatchDescr(fileId, batchId int64) (bool, *dscom.BatchDescr, error) {
    var err error
    exists := false
    var batchDescr *dscom.BatchDescr

    batchDescrs := make([]*dscom.BatchDescr, 0)
    request := `
        SELECT file_id, batch_id, batch_size, block_size
        FROM fs_batchs
        WHERE file_id = $1 AND batch_id = $2
        LIMIT 1;`
    err = reg.db.Select(&batchDescrs, request, fileId, batchId)
    if err != nil {
        return exists, batchDescr, dserr.Err(err)
    }
    if len(batchDescrs) > 0 {
        exists = true
        batchDescr = batchDescrs[0]
    }
    return exists, batchDescr, dserr.Err(err)
}


func (reg *Reg) ListBatchDescrsByFileId(fileId int64) ([]*dscom.BatchDescr, error) {
    var err error
    batchs := make([]*dscom.BatchDescr, 0)
    request := `
        SELECT file_id, batch_id, batch_size, block_size
        FROM fs_batchs
        WHERE file_id = $1
        ORDER BY file_id, batch_id;`
    err = reg.db.Select(&batchs, request, fileId)
    if err != nil {
        return batchs, dserr.Err(err)
    }
    return batchs, dserr.Err(err)
}


func (reg *Reg) EraseBatchDescr(fileId, batchId int64) error {
    var err error
    request := `
        DELETE FROM fs_batchs
        WHERE file_id = $1 AND batch_id = $2;`
    _, err = reg.db.Exec(request, fileId, batchId)
    if err != nil {
        return dserr.Err(err)
    }
    return dserr.Err(err)
}
