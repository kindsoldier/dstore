/*
 * Copyright 2022 Oleg Borodin  <borodin@unix7.org>
 */
package fsreg

import (
    "time"
    "ndstore/dscom"
    "ndstore/dserr"
)

const fileSchema = `
    --- DROP TABLE IF EXISTS fs_fileids;
    CREATE TABLE IF NOT EXISTS fs_fileids (
        file_id         BIGINT GENERATED ALWAYS AS IDENTITY (START 1 CYCLE),
        created_at      BIGINT
    );

    --- DROP TABLE IF EXISTS fs_files;
    CREATE TABLE IF NOT EXISTS fs_files (
        file_id         BIGINT,

        file_ver        BIGINT,
        u_counter       BIGINT,

        batch_size      BIGINT,
        block_size      BIGINT,

        batch_count     BIGINT,
        file_size       BIGINT,
        created_at      BIGINT,
        updated_at      BIGINT,

        is_distr        BOOL
    );
    --- DROP INDEX IF EXISTS fs_file_idx;
    CREATE UNIQUE INDEX IF NOT EXISTS fs_file_idx
        ON fs_files(file_id, file_ver);`



func (reg *Reg) GetNewFileId() (int64, error) {
    var err error
    var fileId int64
    request := `
        INSERT INTO fs_fileids(created_at) VALUES ($1) RETURNING file_id;`
    ts := time.Now().Unix()
    err = reg.db.Get(&fileId, request, ts)
    if err != nil {
        return fileId, dserr.Err(err)
    }
    return fileId, dserr.Err(err)
}


func (reg *Reg) AddNewFileDescr(descr *dscom.FileDescr) error {
    var err error
    request := `
        INSERT INTO fs_files(file_id, batch_count, file_ver, u_counter,
                                                        batch_size, block_size, file_size,
                                                            created_at, updated_at, is_distr)
        VALUES ($1, $2, $3, $4, $5, $6, $7, $8, $9, $10);`

    _, err = reg.db.Exec(request, descr.FileId, descr.BatchCount, descr.FileVer, descr.UCounter,
                                                    descr.BatchSize, descr.BlockSize, descr.FileSize,
                                                    descr.CreatedAt, descr.UpdatedAt, descr.IsDistr)

    if err != nil {
        return dserr.Err(err)
    }
    return dserr.Err(err)
}

func (reg *Reg) GetNewestFileDescr(fileId int64) (bool, *dscom.FileDescr, error) {
    var err error
    var exists bool
    var descr *dscom.FileDescr
    request := `
        SELECT file_id, batch_count, file_ver, u_counter, batch_size, block_size, file_size,
                                                                created_at, updated_at, is_distr
        FROM fs_files
        WHERE file_id = $1
            AND u_counter > 0
        ORDER BY file_ver DESC
        LIMIT 1;`
    descrs := make([]*dscom.FileDescr, 0)
    err = reg.db.Select(&descrs, request, fileId)
    if err != nil {
        return exists, descr, dserr.Err(err)
    }
    if len(descrs) > 0 {
        exists = true
        descr = descrs[0]
    }
    return exists, descr, dserr.Err(err)
}

func (reg *Reg) GetSpecFileDescr(fileId, fileVer int64) (bool, *dscom.FileDescr, error) {
    var err error
    var exists bool
    var descr *dscom.FileDescr
    request := `
        SELECT file_id, batch_count, file_ver, u_counter, batch_size, block_size, file_size,
                                                                created_at, updated_at, is_distr
        FROM fs_files
        WHERE file_id = $1
            AND file_ver = $2
            AND u_counter > 0
        LIMIT 1;`
    descrs := make([]*dscom.FileDescr, 0)
    err = reg.db.Select(&descrs, request, fileId, fileVer)
    if err != nil {
        return exists, descr, dserr.Err(err)
    }
    if len(descrs) > 0 {
        exists = true
        descr = descrs[0]
    }
    return exists, descr, dserr.Err(err)
}


func (reg *Reg) GetSpecUnusedFileDescr(fileId, fileVer int64) (bool, *dscom.FileDescr, error) {
    var err error
    var exists bool
    var descr *dscom.FileDescr
    request := `
        SELECT file_id, batch_count, file_ver, u_counter, batch_size, block_size, file_size,
                                                                created_at, updated_at, is_distr
        FROM fs_files
        WHERE file_id = $1
            AND file_ver = $2
            AND u_counter < 1
        LIMIT 1;`
    descrs := make([]*dscom.FileDescr, 0)
    err = reg.db.Select(&descrs, request, fileId, fileVer)
    if err != nil {
        return exists, descr, dserr.Err(err)
    }
    if len(descrs) > 0 {
        exists = true
        descr = descrs[0]
    }
    return exists, descr, dserr.Err(err)
}


func (reg *Reg) GetAnyUnusedFileDescr() (bool, *dscom.FileDescr, error) {
    var err     error
    var exists  bool
    var blockDescr *dscom.FileDescr
    blocks := make([]*dscom.FileDescr, 0)
    request := `
        SELECT file_id, batch_count, file_ver, u_counter, batch_size, block_size, file_size,
                                                                created_at, updated_at, is_distr
        FROM fs_files
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

func (reg *Reg) IncSpecFileDescrUC(count, fileId, fileVer int64) error {
    var err error
    request := `
        UPDATE fs_files SET
            u_counter = u_counter + $1
        WHERE file_id = $2
            AND file_ver = $3;`
    _, err = reg.db.Exec(request, count, fileId, fileVer)
    if err != nil {
        return dserr.Err(err)
    }
    return dserr.Err(err)
}

func (reg *Reg) DecSpecFileDescrUC(count, fileId, fileVer int64) error {
    var err error
    request := `
        UPDATE fs_files SET
            u_counter = u_counter - $1
        WHERE file_id = $2
            AND file_ver = $3
            AND u_counter > 0;`
    _, err = reg.db.Exec(request, count, fileId, fileVer)
    if err != nil {
        return dserr.Err(err)
    }
    return dserr.Err(err)
}

func (reg *Reg) ListAllFileDescrs() ([]*dscom.FileDescr, error) {
    var err error
    blocks := make([]*dscom.FileDescr, 0)
    request := `
        SELECT file_id, batch_count, file_ver, u_counter, batch_size, block_size, file_size,
                                                                created_at, updated_at, is_distr
        FROM fs_files;`
    err = reg.db.Select(&blocks, request)
    if err != nil {
        return blocks, dserr.Err(err)
    }
    return blocks, dserr.Err(err)
}

func (reg *Reg) EraseSpecFileDescr(fileId, fileVer int64) error {
    var err error
    tx, err := reg.db.Begin()
    if err != nil {
        return dserr.Err(err)
    }
    request1 := `
        DELETE FROM fs_fileids
            WHERE file_id = $1;`
    _, err = tx.Exec(request1, fileId)

    request2 := `
        DELETE FROM fs_files
            WHERE file_id = $1
            AND file_ver = $2;`
    _, err = tx.Exec(request2, fileId, fileVer)
    if err != nil {
        return dserr.Err(err)
    }
    err = tx.Commit()
    if err != nil {
        return dserr.Err(err)
    }
    return dserr.Err(err)
}

func (reg *Reg) GetAnyNotDistrFileDescr() (bool, *dscom.FileDescr, error) {
    var err     error
    var exists  bool
    var blockDescr *dscom.FileDescr
    blocks := make([]*dscom.FileDescr, 0)
    request := `
        SELECT file_id, batch_count, file_ver, u_counter, batch_size, block_size, file_size,
                                                                created_at, updated_at, is_distr
        FROM fs_files
        WHERE u_counter > 0
            AND is_distr = FALSE
            AND ($1 - updated_at) > $2
        ORDER BY file_ver DESC
        LIMIT 1;`
    now := time.Now().Unix()
    gap := int64(30)
    err = reg.db.Select(&blocks, request, now, gap)
    if err != nil {
        return exists, blockDescr, dserr.Err(err)
    }
    if len(blocks) > 0 {
        exists = true
        blockDescr = blocks[0]
    }
    return exists, blockDescr, dserr.Err(err)
}



func (reg *Reg) EraseAllFileDescrs() error {
    var err error
    request := `
        DELETE FROM fs_files;`
    _, err = reg.db.Exec(request)
    if err != nil {
        return dserr.Err(err)
    }
    return dserr.Err(err)
}
