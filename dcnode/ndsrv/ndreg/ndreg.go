/*
 * Copyright 2022 Oleg Borodin  <borodin@unix7.org>
 */
package ndreg

import (
    "errors"

    "github.com/jmoiron/sqlx"
    _ "github.com/mattn/go-sqlite3"

    "dcstore/dccom"
)

const schema = `
    DROP TABLE IF EXISTS blocks;
    CREATE TABLE IF NOT EXISTS blocks (
        cluster_id  INTEGER,
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
        ON blocks (cluster_id, file_id, batch_id, block_id);
    `

var ErrorNilRef error = errors.New("db ref is nil")

type Block = dccom.Block

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

func (reg *Reg) AddBlock(clusterId, fileId, batchId, blockId, blockSize int64,
                                                                    filePath string) error {
    var err error
    if reg.db == nil {
        return ErrorNilRef
    }
    request := `
        INSERT
            INTO blocks(cluster_id, file_id, batch_id, block_id, block_size, file_path)
        VALUES ($1, $2, $3, $4, $5, $6)`
    _, err = reg.db.Exec(request, clusterId, fileId, batchId, blockId, blockSize, filePath)
    if err != nil {
        return err
    }
    return err
}


func (reg *Reg) UpdateBlock(clusterId, fileId, batchId, blockId, blockSize int64,
                                                                    filePath string) error {
    var err error
    if reg.db == nil {
        return ErrorNilRef
    }
    tx, err := reg.db.Begin()
    var request string
    request = `DELETE FROM blocks
                WHERE cluster_id = $1
                    AND file_id = $2
                    AND batch_id = $3
                    AND block_id = $4`
    _, err = tx.Exec(request, clusterId, fileId, batchId, blockId)
    if err != nil {
        return err
    }
    request = `INSERT
                INTO blocks(cluster_id, file_id, batch_id, block_id, block_size, file_path)
                VALUES ($1, $2, $3, $4, $5, $6, $7)`
    _, err = tx.Exec(request, clusterId, fileId, batchId, blockId, blockSize, filePath)
    if err != nil {
        return err
    }

    err = tx.Commit()
    if err != nil {
        return err
    }
    return err
}

func (reg *Reg) GetBlock(clusterId, fileId, batchId, blockId int64) (string, int64, error) {
    var err error
    var filePath string
    var blockSize int64

    if reg.db == nil {
        return filePath, blockSize, ErrorNilRef
    }
    request := `SELECT file_path, block_size
                FROM blocks
                WHERE cluster_id = $1
                    AND file_id = $2
                    AND batch_id = $3
                    AND block_id = $4
                LIMIT 1`

    var block Block
    err = reg.db.Get(&block, request, clusterId, fileId, batchId, blockId)
    if err != nil {
        return filePath, blockSize, err
    }
    filePath    = block.FileName
    blockSize   = block.BlockSize
    return filePath, blockSize, err
}


func (reg *Reg) BlockExists(clusterId, fileId, batchId, blockId int64) (bool, error) {
    var err error
    var exists bool
    if reg.db == nil {
        return exists, ErrorNilRef
    }

    request := `SELECT file_path
                FROM blocks
                WHERE cluster_id = $1
                    AND file_id = $2
                    AND batch_id = $3
                    AND block_id = $4
                LIMIT 1`

    blocks := make([]Block, 0)
    err = reg.db.Select(&blocks, request, clusterId, fileId, batchId, blockId)
    if err != nil {
        return exists, err
    }
    if len(blocks) > 0 {
        exists = true
    }
    return exists, err
}

func (reg *Reg) ListBlocks(clusterId int64) ([]Block, error) {
    var err error
    blocks := make([]Block, 0)
    if reg.db == nil {
        return blocks, ErrorNilRef
    }
    request := `SELECT file_path
                FROM blocks
                WHERE cluster_id = $1`
    err = reg.db.Select(&blocks, request, clusterId)
    if err != nil {
        return blocks, err
    }
    return blocks, err
}

func (reg *Reg) DeleteBlock(clusterId, fileId, batchId, blockId int64) error {
    var err error
    if reg.db == nil {
        return ErrorNilRef
    }
    request := `DELETE FROM blocks
                WHERE cluster_id = $1
                    AND file_id = $2
                    AND batch_id = $3
                    AND block_id = $4`
    _, err = reg.db.Exec(request, clusterId, fileId, batchId, blockId)
    if err != nil {
        return err
    }
    return err
}

func (reg *Reg) PurgeFile(clusterId, fileId int64) error {
    var err error
    if reg.db == nil {
        return ErrorNilRef
    }
    request := `DELETE FROM blocks
                WHERE cluster_id = $1
                    AND file_id = $2`
    _, err = reg.db.Exec(request, clusterId, fileId)
    if err != nil {
        return err
    }
    return err
}

func (reg *Reg) PurgeCluster(clusterId int64) error {
    var err error
    if reg.db == nil {
        return ErrorNilRef
    }
    request := `DELETE FROM blocks
                WHERE cluster_id = $1`
    _, err = reg.db.Exec(request, clusterId)
    if err != nil {
        return err
    }
    return err
}