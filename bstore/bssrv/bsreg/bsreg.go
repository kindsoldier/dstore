/*
 * Copyright 2022 Oleg Borodin  <borodin@unix7.org>
 */
package bsreg

import (
    "errors"

    "github.com/jmoiron/sqlx"
    _ "github.com/mattn/go-sqlite3"

    "ndstore/dscom"
)

const schema = `
    DROP TABLE IF EXISTS blocks;
    CREATE TABLE IF NOT EXISTS blocks (
        file_id     INTEGER,
        batch_id    INTEGER,
        block_id    INTEGER,
        block_size  INTEGER,
        file_path   TEXT,
        hash_alg    TEXT,
        hash_sum    TEXT,
        hash_init   TEXT
    );
    DROP INDEX IF EXISTS block_idx;
    CREATE UNIQUE INDEX IF NOT EXISTS block_idx
        ON blocks (file_id, batch_id, block_id);
    `

var ErrorNilRef error = errors.New("db ref is nil")

type Block = dscom.BlockMI

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

func (reg *Reg) AddBlock(fileId, batchId, blockId, blockSize int64, filePath string) error {
    var err error
    if reg.db == nil {
        return ErrorNilRef
    }
    request := `
        INSERT
            INTO blocks(file_id, batch_id, block_id, block_size, file_path)
        VALUES ($1, $2, $3, $4, $5)`
    _, err = reg.db.Exec(request, fileId, batchId, blockId, blockSize, filePath)
    if err != nil {
        return err
    }
    return err
}


func (reg *Reg) UpdateBlock(fileId, batchId, blockId, blockSize int64, filePath string) error {
    var err error
    if reg.db == nil {
        return ErrorNilRef
    }
    tx, err := reg.db.Begin()
    var request string
    request = `DELETE FROM blocks
                WHERE file_id = $2
                    AND batch_id = $3
                    AND block_id = $4`
    _, err = tx.Exec(request, fileId, batchId, blockId)
    if err != nil {
        return err
    }
    request = `INSERT
                INTO blocks(file_id, batch_id, block_id, block_size, file_path)
                VALUES ($1, $2, $3, $4, $5)`
    _, err = tx.Exec(request, fileId, batchId, blockId, blockSize, filePath)
    if err != nil {
        return err
    }

    err = tx.Commit()
    if err != nil {
        return err
    }
    return err
}

func (reg *Reg) GetBlock(fileId, batchId, blockId int64) (string, int64, error) {
    var err error
    var filePath string
    var blockSize int64

    if reg.db == nil {
        return filePath, blockSize, ErrorNilRef
    }
    request := `SELECT file_path, block_size
                FROM blocks
                WHERE file_id = $2
                    AND batch_id = $3
                    AND block_id = $4
                LIMIT 1`

    var block Block
    err = reg.db.Get(&block, request, fileId, batchId, blockId)
    if err != nil {
        return filePath, blockSize, err
    }
    filePath    = block.FileName
    blockSize   = block.BlockSize
    return filePath, blockSize, err
}


func (reg *Reg) BlockExists(fileId, batchId, blockId int64) (bool, error) {
    var err error
    var exists bool
    if reg.db == nil {
        return exists, ErrorNilRef
    }

    request := `SELECT file_path
                FROM blocks
                WHERE file_id = $2
                    AND batch_id = $3
                    AND block_id = $4
                LIMIT 1`

    blocks := make([]Block, 0)
    err = reg.db.Select(&blocks, request, fileId, batchId, blockId)
    if err != nil {
        return exists, err
    }
    if len(blocks) > 0 {
        exists = true
    }
    return exists, err
}

func (reg *Reg) ListBlocks() ([]Block, error) {
    var err error
    blocks := make([]Block, 0)
    if reg.db == nil {
        return blocks, ErrorNilRef
    }
    request := `SELECT file_path
                FROM blocks`
    err = reg.db.Select(&blocks, request)
    if err != nil {
        return blocks, err
    }
    return blocks, err
}

func (reg *Reg) DeleteBlock(fileId, batchId, blockId int64) error {
    var err error
    if reg.db == nil {
        return ErrorNilRef
    }
    request := `DELETE FROM blocks
                WHERE file_id = $2
                    AND batch_id = $3
                    AND block_id = $4`
    _, err = reg.db.Exec(request, fileId, batchId, blockId)
    if err != nil {
        return err
    }
    return err
}

func (reg *Reg) PurgeFile(fileId int64) error {
    var err error
    if reg.db == nil {
        return ErrorNilRef
    }
    request := `DELETE FROM blocks
                WHERE file_id = $2`
    _, err = reg.db.Exec(request, fileId)
    if err != nil {
        return err
    }
    return err
}

func (reg *Reg) PurgeCluster(userId int64) error {
    var err error
    if reg.db == nil {
        return ErrorNilRef
    }
    request := `DELETE FROM blocks`
    _, err = reg.db.Exec(request, userId)
    if err != nil {
        return err
    }
    return err
}
