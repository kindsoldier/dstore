/*
 * Copyright 2022 Oleg Borodin  <borodin@unix7.org>
 */
package fsreg

import (
    "ndstore/dscom"
)


const entriesSchema = `
    DROP TABLE IF EXISTS entries;
    CREATE TABLE IF NOT EXISTS entries (
        entry_id   INTEGER GENERATED ALWAYS AS IDENTITY (START 1 CYCLE),
        user_id    INTEGER,
        file_id    INTEGER,
        dir_path   TEXT,
        file_name  TEXT
    );
    DROP INDEX IF EXISTS entry_idx;
    CREATE UNIQUE INDEX IF NOT EXISTS entry_idx
        ON entries (user_id, dir_path, file_name);
    `

func (reg *Reg) AddEntryDescr(userId int64, dirPath, fileName string, fileId int64) error {
    var err error
    request := `
        INSERT INTO entries(user_id, dir_path, file_name, file_id)
        VALUES ($1, $2, $3, $4);`
    _, err = reg.db.Exec(request, userId, dirPath, fileName, fileId)
    if err != nil {
        return err
    }
    return err
}

func (reg *Reg) EntryDescrExists(userId int64, dirPath, fileName string) (bool, error) {
    var err error
    var exists bool
    request := `
        SELECT count(entry_id) AS count
        FROM entries
        WHERE user_id = $1 AND dir_path = $2 AND file_name = $3
        LIMIT 1;`
    var count int64
    err = reg.db.Get(&count, request, userId, dirPath, fileName)
    if err != nil {
        return exists, err
    }
    if count > 0 {
        exists = true
    }
    return exists, err
}

func (reg *Reg) GetEntryDescr(userId int64, dirPath, fileName string) (*dscom.EntryDescr, error) {
    var err error
    request := `
        SELECT dir_path, file_name, file_id, user_id
        FROM entries
        WHERE  user_id = $1 AND dir_path = $2 AND file_name = $3
        LIMIT 1;`
    entry := dscom.NewEntryDescr()
    err = reg.db.Get(entry, request, userId, dirPath, fileName)
    if err != nil {
        return entry, err
    }
    return entry, err
}

func (reg *Reg) DeleteEntryDescr(userId int64, dirPath, fileName string) error {
    var err error
    request := `
        DELETE FROM entries
        WHERE user_id = $1 AND dir_path = $2 AND file_name = $3;`
    _, err = reg.db.Exec(request, userId, dirPath, fileName)
    if err != nil {
        return err
    }
    return err
}

func (reg *Reg) ListEntryDescr(userId int64, dirPath string) ([]*dscom.EntryDescr, error) {
    var err error
    request := `
        SELECT e.entry_id, e.user_id, e.dir_path, e.file_name, e.file_id, f.file_size
        FROM entries AS e, files AS f
        WHERE e.file_id = f.file_id
            AND e.user_id = $1;`
    entries := make([]*dscom.EntryDescr, 0)
    err = reg.db.Select(&entries, request, userId)
    if err != nil {
        return entries, err
    }
    return entries, err
}
