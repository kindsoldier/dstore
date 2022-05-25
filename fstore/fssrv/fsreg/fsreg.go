/*
 * Copyright 2022 Oleg Borodin  <borodin@unix7.org>
 */
package fsreg

import (
    "errors"
    "fmt"

    "github.com/jmoiron/sqlx"
    _ "github.com/mattn/go-sqlite3"

    "ndstore/dscom"
)


const schema = `
    DROP TABLE IF EXISTS entries;
    CREATE TABLE IF NOT EXISTS entries (
        entry_dir   TEXT,
        entry_name  TEXT,
        file_id     INTEGER
    );

    DROP TABLE IF EXISTS files;
    CREATE TABLE IF NOT EXISTS files (
        file_id     INTEGER,
        batch_count INTEGER,
        batch_size  INTEGER,
        block_size  INTEGER,
        file_size   INTEGER
    );

    DROP TABLE IF EXISTS batchs;
    CREATE TABLE IF NOT EXISTS batchs (
        file_id     INTEGER,
        batch_id    INTEGER,
        batch_size  INTEGER,
        block_size  INTEGER
    );

    DROP TABLE IF EXISTS blocks;
    CREATE TABLE IF NOT EXISTS blocks (
        file_id     INTEGER,
        batch_id    INTEGER,
        block_id    INTEGER,
        block_size  INTEGER,
        file_path   TEXT,
        hash_init   TEXT,
        hash_sum    TEXT,
        data_size  INTEGER
    );`

var ErrorNilRef error = errors.New("db ref is nil")

type Reg struct {
    db *sqlx.DB
}

func NewReg() *Reg {
    var reg Reg
    return &reg
}

func (reg *Reg) OpenDB(dbPath string) error {
    var err error
    db, err := sqlx.Open("sqlite3", dbPath)
    if err != nil {
        return err
    }
    err = db.Ping()
    if err != nil {
        return err
    }
    reg.db = db
    return err
}

func (reg *Reg) CloseDB() error {
    var err error
    if reg.db == nil {
        return ErrorNilRef
    }
    reg.db.Close()
    return err
}

func (reg *Reg) MigrateDB() error {
    var err error
    if reg.db == nil {
        return ErrorNilRef
    }
    _, err = reg.db.Exec(schema)
    if err != nil {
        return err
    }
    return err
}

func (reg *Reg) AddFileDescr(file *dscom.FileDescr) error {
    var err error

    tx, err := reg.db.Begin()

    blockRequest := `
        INSERT INTO blocks(file_id, batch_id, block_id, block_size, file_path,
                                                                hash_init, hash_sum, data_size)
                VALUES ($1, $2, $3, $4, $5, $6, $7, $8)`

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
        VALUES ($1, $2, $3, $4)`
    for _, batch := range file.Batchs {
        _, err = tx.Exec(batchRequest, batch.FileId, batch.BatchId, batch.BatchSize, batch.BlockSize)
        if err != nil {
            return err
        }
    }

    fileRequest := `
        INSERT INTO files(file_id, batch_count, batch_size, block_size, file_size)
        VALUES ($1, $2, $3, $4, $5)`
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
        LIMIT 1`

    file := dscom.NewFileDescr()
    err = reg.db.Get(file, fileRequest, fileId)
    if err != nil {
        return file, err
    }

    batchRequest := `
        SELECT file_id, batch_id, block_size
        FROM batchs
        WHERE file_id = $1
        ORDER BY file_id, batch_id
        `
    batchs := make([]*dscom.BatchDescr, 0)


    err = reg.db.Select(&batchs, batchRequest, fileId)
    if err != nil {
        return file, err
    }
    file.Batchs = batchs

    fmt.Println(batchs)


    blockRequest := `
        SELECT file_id, batch_id, block_id, block_size, file_path, hash_init, hash_sum, data_size
        FROM blocks
        WHERE file_id = $1
            AND batch_id = $2
        ORDER BY file_id, batch_id, block_id`
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

func (reg *Reg) DeleteFileDescr() error {
    var err error
    return err
}
