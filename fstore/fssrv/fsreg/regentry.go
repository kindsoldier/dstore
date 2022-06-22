/*
 * Copyright 2022 Oleg Borodin  <borodin@unix7.org>
 */
package fsreg

import (
    "ndstore/dscom"
    "ndstore/dserr"
)


const entrieSchema = `
    DROP TABLE IF EXISTS fs_entries;
    CREATE TABLE IF NOT EXISTS fs_entries (
        entry_id   INTEGER GENERATED ALWAYS AS IDENTITY (START 1 CYCLE),
        user_id    INTEGER,
        file_id    INTEGER,
        dir_path   TEXT,
        file_name  TEXT
    );
    DROP INDEX IF EXISTS fs_entry_idx;
    CREATE UNIQUE INDEX IF NOT EXISTS fs_entry_idx
        ON fs_entries(user_id, dir_path, file_name);
    `

func (reg *Reg) AddEntryDescr(userId int64, dirPath, fileName string, fileId int64) error {
    var err error
    request := `
        INSERT INTO fs_entries(user_id, dir_path, file_name, file_id)
        VALUES ($1, $2, $3, $4);`
    _, err = reg.db.Exec(request, userId, dirPath, fileName, fileId)
    if err != nil {
        return dserr.Err(err)
    }
    return dserr.Err(err)
}

func (reg *Reg) EntryDescrExists(userId int64, dirPath, fileName string) (bool, error) {
    var err error
    var exists bool
    request := `
        SELECT count(entry_id) AS count
        FROM fs_entries
        WHERE user_id = $1 AND dir_path = $2 AND file_name = $3
        LIMIT 1;`
    var count int64
    err = reg.db.Get(&count, request, userId, dirPath, fileName)
    if err != nil {
        return exists, dserr.Err(err)
    }
    if count > 0 {
        exists = true
    }
    return exists, dserr.Err(err)
}

func (reg *Reg) GetEntryDescr(userId int64, dirPath, fileName string) (*dscom.EntryDescr, error) {
    var err error
    request := `
        SELECT dir_path, file_name, file_id, user_id
        FROM fs_entries
        WHERE  user_id = $1 AND dir_path = $2 AND file_name = $3
        LIMIT 1;`
    entry := dscom.NewEntryDescr()
    err = reg.db.Get(entry, request, userId, dirPath, fileName)
    if err != nil {
        return entry, dserr.Err(err)
    }
    return entry, dserr.Err(err)
}

func (reg *Reg) DeleteEntryDescr(userId int64, dirPath, fileName string) error {
    var err error
    request := `
        DELETE FROM fs_entries
        WHERE user_id = $1 AND dir_path = $2 AND file_name = $3;`
    _, err = reg.db.Exec(request, userId, dirPath, fileName)
    if err != nil {
        return dserr.Err(err)
    }
    return dserr.Err(err)
}

func (reg *Reg) ListEntryDescr(userId int64, dirPath string) ([]*dscom.EntryDescr, error) {
    var err error
    request := `
        SELECT e.entry_id, e.user_id, e.dir_path, e.file_name, e.file_id, f.file_size
        FROM fs_entries AS e, fs_files AS f
        WHERE e.file_id = f.file_id
            AND e.user_id = $1;`
    entries := make([]*dscom.EntryDescr, 0)
    err = reg.db.Select(&entries, request, userId)
    if err != nil {
        return entries, dserr.Err(err)
    }
    return entries, dserr.Err(err)
}
