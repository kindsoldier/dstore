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
        dir_path   TEXT,
        file_name  TEXT,
        file_id    INTEGER
    );
    DROP INDEX IF EXISTS entry_idx;
    CREATE UNIQUE INDEX IF NOT EXISTS entry_idx
        ON entries (dir_path, file_name);
    `

func (reg *Reg) AddEntryDescr(dirPath, fileName string, fileId int64) error {
    var err error
    request := `
        INSERT INTO entries(dir_path, file_name, file_id)
        VALUES ($1, $2, $3);`
    _, err = reg.db.Exec(request, dirPath, fileName, fileId)
    if err != nil {
        return err
    }
    return err
}

func (reg *Reg) GetEntryDescr(dirPath, fileName string) (*dscom.EntryDescr, bool, error) {
    var err error
    var exists bool
    var entry *dscom.EntryDescr
    request := `
        SELECT dir_path, file_name, file_id
        FROM entries
        WHERE dir_path = $1
            AND file_name = $2
        LIMIT 1;`
    entries := make([]*dscom.EntryDescr, 0)
    err = reg.db.Select(&entries, request, dirPath, fileName)
    if err != nil {
        return entry, exists, err

    }
    if len(entries) > 0 {
        exists = true
        entry = entries[0]
    }
    return entry, exists, err
}

func (reg *Reg) EntryDescrExists(dirPath, fileName string) (bool, error) {
    var err error
    var exists bool
    request := `
        SELECT dir_path, file_name, file_id
        FROM entries
        WHERE dir_path = $1
            AND file_name = $2
        LIMIT 1;`
    entries := make([]*dscom.EntryDescr, 0)
    err = reg.db.Select(&entries, request, dirPath, fileName)
    if err != nil {
        return exists, err
    }
    if len(entries) > 0 {
        exists = true
    }
    return exists, err
}

func (reg *Reg) DeleteEntryDescr(dirPath, fileName string) error {
    var err error
    request := `
        DELETE FROM entries
        WHERE dir_path = $1
            AND file_name = $2;`
    _, err = reg.db.Exec(request, dirPath, fileName)
    if err != nil {
        return err
    }
    return err
}

func (reg *Reg) ListEntryDescr(dirPath string) ([]*dscom.EntryDescr, error) {
    var err error
    request := `
        SELECT dir_path, file_name, file_id
        FROM entries
        WHERE dir_path = $1;`
    entries := make([]*dscom.EntryDescr, 0)
    err = reg.db.Select(&entries, request, dirPath)
    if err != nil {
        return entries, err
    }
    return entries, err
}
