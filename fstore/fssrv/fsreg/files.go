/*
 * Copyright 2022 Oleg Borodin  <borodin@unix7.org>
 */
package fsreg

import (
    "time"
    "ndstore/dscom"
    "ndstore/dserr"
)

const filesSchema = `
    DROP TABLE IF EXISTS file_ids;
    CREATE TABLE IF NOT EXISTS file_ids (
        file_id     INTEGER GENERATED ALWAYS AS IDENTITY (CYCLE),
        created_at  INTEGER
    );

    DROP TABLE IF EXISTS files;
    CREATE TABLE IF NOT EXISTS files (
        file_id     INTEGER,
        batch_count INTEGER,
        batch_size  INTEGER,
        block_size  INTEGER,
        file_size   INTEGER,
        created_at  INTEGER
    );
    DROP INDEX IF EXISTS file_idx;
    CREATE UNIQUE INDEX IF NOT EXISTS file_idx
        ON files (file_id);

    DROP TABLE IF EXISTS batchs;
    CREATE TABLE IF NOT EXISTS batchs (
        file_id     INTEGER,
        batch_id    INTEGER,
        batch_size  INTEGER,
        block_size  INTEGER
    );
    DROP INDEX IF EXISTS batch_idx;
    CREATE UNIQUE INDEX IF NOT EXISTS batch_idx
        ON batchs (file_id, batch_id);

    DROP TABLE IF EXISTS blocks;
    CREATE TABLE IF NOT EXISTS blocks (
        file_id     INTEGER,
        batch_id    INTEGER,
        block_id    INTEGER,
        block_size  INTEGER,
        block_type  TEXT DEFAULT '',
        file_path   TEXT DEFAULT '',
        data_size   INTEGER,
        hash_init   TEXT DEFAULT '',
        hash_sum    TEXT DEFAULT '',
        hash_alg    TEXT DEFAULT '',
        bstore_id   INTEGER,
        has_remote  BOOL,
        has_local   BOOL
    );
    DROP INDEX IF EXISTS block_idx;
    CREATE UNIQUE INDEX IF NOT EXISTS block_idx
        ON blocks (file_id, batch_id, block_id, block_type);`

func timestamp() int64 {
    return time.Now().UTC().Unix()
}

func (reg *Reg) GetNewFileId() (int64, error) {
    var err error
    var fileId int64
    ts := timestamp()

    var request string
    //holeRequest = `
    //    SELECT t1.file_id - 1 AS file_id
    //    FROM file_ids AS t1
    //    LEFT JOIN file_ids AS t2
    //    ON t1.file_id - 1 = t2.file_id
    //    WHERE t2.file_id IS NULL
    //    LIMIT 1;`

    request = `
        INSERT INTO file_ids(created_at) VALUES($1)
        RETURNING file_id;`
    err = reg.db.Get(&fileId, request, ts)
    if err != nil {
        return fileId, dserr.Err(err)
    }
    return fileId, dserr.Err(err)
}


func (reg *Reg) AddFileDescr(file *dscom.FileDescr) error {
    var err error

    tx, err := reg.db.Begin()
    exitFunc := func() {
        if err != nil && tx != nil {
            tx.Rollback()
        }
    }
    defer exitFunc()
    if err != nil {
        return dserr.Err(err)
    }
    blockRequest := `
        INSERT INTO blocks(file_id, batch_id, block_id, block_size, file_path,
                                        hash_alg, hash_init, hash_sum, data_size,
                                        block_type, bstore_id, has_remote, has_local)
                VALUES ($1, $2, $3, $4, $5, $6, $7, $8, $9, $10, $11, $12, $13);`
    for _, batch := range file.Batchs {
        for _, bl := range batch.Blocks {
            _, err = tx.Exec(blockRequest, bl.FileId, bl.BatchId, bl.BlockId, bl.BlockSize, bl.FilePath,
                                            bl.HashAlg, bl.HashInit, bl.HashSum, bl.DataSize,
                                            bl.BlockType, bl.BStoreId, bl.HasRemote, bl.HasLocal)
            if err != nil {
                return dserr.Err(err)
            }
        }
    }
    batchRequest := `
        INSERT INTO batchs(file_id, batch_id, batch_size, block_size)
        VALUES ($1, $2, $3, $4);`
    for _, batch := range file.Batchs {
        _, err = tx.Exec(batchRequest, batch.FileId, batch.BatchId, batch.BatchSize,
                                                                            batch.BlockSize)
        if err != nil {
            return dserr.Err(err)
        }
    }
    fileRequest := `
        INSERT INTO files(file_id, batch_count, batch_size, block_size, file_size, created_at)
        VALUES ($1, $2, $3, $4, $5, $6);`
    _, err = tx.Exec(fileRequest, file.FileId, file.BatchCount, file.BatchSize,
                                                file.BlockSize, file.FileSize, file.CreatedAt)
    if err != nil {
        return dserr.Err(err)
    }
    err = tx.Commit()
    if err != nil {
        return dserr.Err(err)
    }
    return dserr.Err(err)
}

func (reg *Reg) GetFileDescr(fileId int64) (*dscom.FileDescr, error) {
    var err error

    fileRequest := `
        SELECT file_id, batch_count, batch_size, block_size, file_size, created_at
        FROM files
        WHERE file_id = $1
        LIMIT 1;`
    file := dscom.NewFileDescr()
    err = reg.db.Get(file, fileRequest, fileId)
    if err != nil {
        return file, dserr.Err(err)
    }
    batchRequest := `
        SELECT file_id, batch_id, block_size, batch_size
        FROM batchs
        WHERE file_id = $1
        ORDER BY file_id, batch_id;`
    batchs := make([]*dscom.BatchDescr, 0)
    err = reg.db.Select(&batchs, batchRequest, fileId)
    if err != nil {
        return file, dserr.Err(err)
    }
    file.Batchs = batchs
    blockRequest := `
        SELECT file_id, batch_id, block_id, block_size, file_path,
                hash_alg, hash_init, hash_sum, data_size,
                block_type, bstore_id, has_remote, has_local
        FROM blocks
        WHERE file_id = $1
            AND batch_id = $2
        ORDER BY file_id, batch_id, block_id;`
    for i := range file.Batchs {
        blocks := make([]*dscom.BlockDescr, 0)
        err = reg.db.Select(&blocks, blockRequest, fileId, file.Batchs[i].BatchId)
        if err != nil {
            return file, dserr.Err(err)
        }
        file.Batchs[i].Blocks = blocks
    }
    return file, dserr.Err(err)
}

func (reg *Reg) FileDescrExists(fileId int64) (bool, error) {
    var err error
    var exists bool
    request := `
        SELECT count(file_id)
        FROM files
        WHERE file_id = $1
        LIMIT 1;`
    var count int64
    err = reg.db.Select(&count, request, fileId)
    if err != nil {
        return exists, dserr.Err(err)
    }
    if count > 0 {
        exists = true
    }
    return exists, dserr.Err(err)
}

func (reg *Reg) DeleteFileDescr(fileId int64) error {
    var err error
    tx, err := reg.db.Begin()
    exitFunc := func() {
        if err != nil && tx != nil {
            tx.Rollback()
        }
    }
    defer exitFunc()
    if err != nil {
        return dserr.Err(err)
    }
    var request string
    request = `
        DELETE FROM blocks
        WHERE file_id = $1;`
    _, err = tx.Exec(request, fileId)
    if err != nil {
        return dserr.Err(err)
    }
    request = `
        DELETE FROM batchs
        WHERE file_id = $1;`
    _, err = tx.Exec(request, fileId)
    if err != nil {
        return dserr.Err(err)
    }
    request = `
        DELETE FROM files
        WHERE file_id = $1;`
    _, err = tx.Exec(request, fileId)
    if err != nil {
        return dserr.Err(err)
    }
    request = `
        DELETE FROM file_ids
        WHERE file_id = $1;`
    _, err = tx.Exec(request, fileId)
    if err != nil {
        return dserr.Err(err)
    }
    err = tx.Commit()
    if err != nil {
        return dserr.Err(err)
    }
    return dserr.Err(err)
}
