/*
 * Copyright 2022 Oleg Borodin  <borodin@unix7.org>
 */
package fsreg

import (
    "ndstore/dscom"
)

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
        return err
    }

    blockRequest := `
        INSERT INTO blocks(file_id, batch_id, block_id, block_size, file_path,
                                                                hash_init, hash_sum, data_size)
                VALUES ($1, $2, $3, $4, $5, $6, $7, $8);`

    for _, batch := range file.Batchs {
        for _, block := range batch.Blocks {
            _, err = tx.Exec(blockRequest, block.FileId, block.BatchId, block.BlockId,
                                             block.BlockSize, block.FilePath, block.HashInit,
                                                                  block.HashSum, block.DataSize)
            if err != nil {
                return err
            }
        }
    }

    batchRequest := `
        INSERT INTO batchs(file_id, batch_id, batch_size, block_size)
        VALUES ($1, $2, $3, $4);`
    for _, batch := range file.Batchs {
        _, err = tx.Exec(batchRequest, batch.FileId, batch.BatchId, batch.BatchSize, batch.BlockSize)
        if err != nil {
            return err
        }
    }

    fileRequest := `
        INSERT INTO files(file_id, batch_count, batch_size, block_size, file_size)
        VALUES ($1, $2, $3, $4, $5);`
    _, err = tx.Exec(fileRequest, file.FileId, file.BatchCount, file.BatchSize,
                                                            file.BlockSize, file.FileSize)
    if err != nil {
        return err
    }

    err = tx.Commit()
    if err != nil {
        return err
    }
    return err
}

func (reg *Reg) GetFileDescr(fileId int64) (*dscom.FileDescr, error) {
    var err error

    fileRequest := `
        SELECT file_id, batch_count, batch_size, block_size, file_size
        FROM files
        WHERE file_id = $1
        LIMIT 1;`

    file := dscom.NewFileDescr()
    err = reg.db.Get(file, fileRequest, fileId)
    if err != nil {
        return file, err
    }

    batchRequest := `
        SELECT file_id, batch_id, block_size, batch_size
        FROM batchs
        WHERE file_id = $1
        ORDER BY file_id, batch_id;`
    batchs := make([]*dscom.BatchDescr, 0)

    err = reg.db.Select(&batchs, batchRequest, fileId)
    if err != nil {
        return file, err
    }
    file.Batchs = batchs

    blockRequest := `
        SELECT file_id, batch_id, block_id, block_size, file_path, hash_init, hash_sum, data_size
        FROM blocks
        WHERE file_id = $1
            AND batch_id = $2
        ORDER BY file_id, batch_id, block_id;`
    for i := range file.Batchs {
        blocks := make([]*dscom.BlockDescr, 0)
        err = reg.db.Select(&blocks, blockRequest, fileId, file.Batchs[i].BatchId)
        if err != nil {
            return file, err
        }
        file.Batchs[i].Blocks = blocks
    }
    return file, err
}

func (reg *Reg) FileDescrExists(fileId int64) (bool, error) {
    var err error
    var exists bool

    request := `
        SELECT file_id, batch_count, batch_size, block_size, file_size
        FROM files
        WHERE file_id = $1
        LIMIT 1;`
    files := make([]*dscom.FileDescr, 0)
    err = reg.db.Select(&files, request, fileId)
    if err != nil {
        return exists, err
    }
    if len(files) > 0 {
        exists = true
    }
    return exists, err
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
        return err
    }

    var request string
    request = `
        DELETE FROM blocks
        WHERE file_id = $1;`
    _, err = tx.Exec(request, fileId)
    if err != nil {
        return err
    }
    request = `
        DELETE FROM batchs
        WHERE file_id = $1;`
    _, err = tx.Exec(request, fileId)
    if err != nil {
        return err
    }
    request = `
        DELETE FROM files
        WHERE file_id = $1;`
    _, err = tx.Exec(request, fileId)
    if err != nil {
        return err
    }
    err = tx.Commit()
    if err != nil {
        return err
    }
    return err
}
