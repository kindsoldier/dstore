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
    --- DROP TABLE IF EXISTS fs_files;
    CREATE TABLE IF NOT EXISTS fs_files (
        file_id         INTEGER GENERATED ALWAYS AS IDENTITY (START 1 CYCLE ),
        batch_size      INTEGER,
        block_size      INTEGER,
        batch_count     INTEGER,
        u_counter       INTEGER,
        file_size       INTEGER,
        created_at      INTEGER,
        updated_at      INTEGER
    );
    --- DROP INDEX IF EXISTS fs_file_idx;
    CREATE UNIQUE INDEX IF NOT EXISTS fs_file_idx
        ON fs_files(file_id);`


func (reg *Reg) AddFileDescr(descr *dscom.FileDescr) (int64, error) {
    var err error
    request := `
        INSERT INTO fs_files(batch_size, block_size, u_counter, batch_count, file_size, created_at, updated_at)
        VALUES ($1, $2, $3, $4, $5, $6, $7)
        RETURNING file_id;`
    var fileId int64
    createdAt := time.Now().Unix()
    updatedAt := createdAt
    err = reg.db.Get(&fileId, request, descr.BatchSize, descr.BlockSize, descr.UCounter, descr.BatchCount,
                                                                                descr.FileSize,
                                                                                createdAt, updatedAt)
    if err != nil {
        return fileId, dserr.Err(err)
    }
    return fileId, dserr.Err(err)
}

func (reg *Reg) UpdateFileDescr(descr *dscom.FileDescr) error {
    var err error
    updatedAt := time.Now().Unix()
    request := `
        UPDATE fs_files SET batch_size = $1, block_size = $2, batch_count = $3, file_size = $4, updated_at = $5
        WHERE file_id = $6;`
    _, err = reg.db.Exec(request, descr.BatchSize, descr.BlockSize, descr.BatchCount,
                                                                                descr.FileSize, updatedAt,
                                                                                descr.FileId)
    if err != nil {
        return dserr.Err(err)
    }
    return dserr.Err(err)
}

func (reg *Reg) GetFileDescr(fileId int64) (bool, *dscom.FileDescr, error) {
    var err error
    exists := false
    var fileDescr *dscom.FileDescr

    fileDescrs := make([]*dscom.FileDescr, 0)
    request := `
        SELECT file_id, batch_size, block_size, u_counter, batch_count, file_size, created_at, updated_at
        FROM fs_files
        WHERE file_id = $1
        LIMIT 1;`
    err = reg.db.Select(&fileDescrs, request, fileId)
    if err != nil {
        return exists, fileDescr, dserr.Err(err)
    }
    if len(fileDescrs) > 0 {
        exists = true
        fileDescr = fileDescrs[0]
    }
    return exists, fileDescr, dserr.Err(err)
}


func (reg *Reg) ListFileDescrs() ([]*dscom.FileDescr, error) {
    var err error
    files := make([]*dscom.FileDescr, 0)
    request := `
        SELECT file_id, batch_size, block_size, u_counter, batch_count, file_size, created_at, updated_at
        FROM fs_files
        ORDER BY file_id;`
    err = reg.db.Select(&files, request)
    if err != nil {
        return files, dserr.Err(err)
    }
    return files, dserr.Err(err)
}

//func (reg *Reg) GetUnusedBlockDescr() (bool, *dscom.BlockDescr, error) {
//    var err     error
//    var exists  bool
//    var blockDescr *dscom.BlockDescr
//    blocks := make([]*dscom.BlockDescr, 0)
//    request := `
//        SELECT b.file_id, b.batch_id, b.block_id, b.block_size, b.data_size,
//                                b.file_path, b.block_type, b.hash_alg, b.hash_init, b.hash_sum,
//                                fstore_id, b.bstore_id, b.saved_loc, b.saved_rem
//        FROM fs_blocks AS b, fs_files as f
//        WHERE f.u_counter < 1 AND b.file_id = b.file_id
//        LIMIT 1;`
//    err = reg.db.Select(&blocks, request)
//    if err != nil {
//        return exists, blockDescr, dserr.Err(err)
//    }
//    if len(blocks) > 0 {
//        exists = true
//        blockDescr = blocks[0]
//    }
//    return exists, blockDescr, dserr.Err(err)
//}

func (reg *Reg) GetUnusedFileDescr() (bool, *dscom.FileDescr, error) {
    var err     error
    var exists  bool
    var fileDescr *dscom.FileDescr
    files := make([]*dscom.FileDescr, 0)
    request := `
        SELECT file_id, batch_size, block_size, u_counter, batch_count, file_size
        FROM fs_files
        WHERE u_counter < 1
        ORDER BY file_id
        LIMIT 1;`
    err = reg.db.Select(&files, request)
    if err != nil {
        return exists, fileDescr, dserr.Err(err)
    }
    if len(files) > 0 {
        exists = true
        fileDescr = files[0]
    }
    return exists, fileDescr, dserr.Err(err)
}


func (reg *Reg) GetLostedFileDescr() (bool, *dscom.FileDescr, error) {
    var err     error
    var exists  bool
    var fileDescr *dscom.FileDescr
    files := make([]*dscom.FileDescr, 0)
    request := `
        SELECT f.* FROM fs_files AS f
        LEFT JOIN fs_entries AS e ON e.file_id = f.file_id
        WHERE e.entry_id IS NULL
        ORDER BY f.file_id
        LIMIT 1;`
    err = reg.db.Select(&files, request)
    if err != nil {
        return exists, fileDescr, dserr.Err(err)
    }
    if len(files) > 0 {
        exists = true
        fileDescr = files[0]
    }
    return exists, fileDescr, dserr.Err(err)
}


func (reg *Reg) IncFileDescrUC(fileId int64) error {
    var err error
    request := `
        UPDATE fs_files SET
            u_counter = u_counter + 1
        WHERE file_id = $1;`
    _, err = reg.db.Exec(request, fileId)
    if err != nil {
        return dserr.Err(err)
    }
    return dserr.Err(err)
}

func (reg *Reg) DecFileDescrUC(fileId int64) error {
    var err error
    request := `
        UPDATE fs_files SET
            u_counter = u_counter - 1
        WHERE file_id = $1 AND u_counter > 0;`
    _, err = reg.db.Exec(request, fileId)
    if err != nil {
        return dserr.Err(err)
    }
    return dserr.Err(err)
}

func (reg *Reg) EraseFileDescr(fileId int64) error {
    var err error
    var request string
    tx, err := reg.db.Begin()
    if err != nil {
        return dserr.Err(err)
    }
    request = `
        DELETE FROM fs_files
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
