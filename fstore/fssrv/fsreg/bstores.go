/*
 * Copyright 2022 Oleg Borodin  <borodin@unix7.org>
 */
package fsreg

import (
    "ndstore/dscom"
)

const bstoresSchema = `
    DROP TABLE IF EXISTS bstores;
    CREATE TABLE IF NOT EXISTS bstores (
        id          INTEGER UNIQUE,
        address     TEXT,
        login       TEXT,
        pass        TEXT,
        state       TEXT
    );
    DROP INDEX IF EXISTS bstore_idx;
    CREATE UNIQUE INDEX IF NOT EXISTS bstore_idx
        ON bstores (id, address);
    `

func (reg *Reg) AddBStoreDescr(id int64, address, login, pass, state string) error {
    var err error
    request := `
        INSERT INTO bstores(id, address, login, pass, state)
        VALUES ($1, $2, $3, $4, $5);`
    _, err = reg.db.Exec(request, id, address, login, pass, state)
    if err != nil {
        return err
    }
    return err
}

func (reg *Reg) UpdateBStoreDescr(id int64, address, login, pass, state string) error {
    var err error
    request := `
        UPDATE bstores
        SET address = $2, login = $3, pass = $4, state = $5
        WHERE id = $1;`
    _, err = reg.db.Exec(request, id, address, login, pass, state)
    if err != nil {
        return err
    }
    return err
}

func (reg *Reg) RenewBStoreDescr(descr *dscom.BStoreDescr) error {
    var err error
    request := `
        UPDATE bstores
        SET address = $2, login = $3, pass = $4, state = $5
        WHERE id = $1;`
    _, err = reg.db.Exec(request, descr.Id, descr.Address, descr.Login, descr.Pass, descr.State)
    if err != nil {
        return err
    }
    return err
}


func (reg *Reg) GetBStoreDescr(id int64) (*dscom.BStoreDescr, bool, error) {
    var err error
    var exists bool
    var bstore *dscom.BStoreDescr
    request := `
        SELECT id, address, login, pass, state
        FROM bstores
        WHERE id = $1
        LIMIT 1;`
    bstores := make([]*dscom.BStoreDescr, 0)
    err = reg.db.Select(&bstores, request, id)
    if err != nil {
        return bstore, exists, err

    }
    if len(bstores) > 0 {
        exists = true
        bstore = bstores[0]
    }
    return bstore, exists, err
}

func (reg *Reg) BStoreDescrExists(id int64) (bool, error) {
    var err error
    var exists bool
    request := `
        SELECT id, address, login, pass, state
        FROM bstores
        WHERE id = $1
        LIMIT 1;`
    bstores := make([]*dscom.BStoreDescr, 0)
    err = reg.db.Select(&bstores, request, id)
    if err != nil {
        return exists, err
    }
    if len(bstores) > 0 {
        exists = true
    }
    return exists, err
}

func (reg *Reg) DeleteBStoreDescr(id int64) error {
    var err error
    request := `
        DELETE FROM bstores
        WHERE id = $1;`
    _, err = reg.db.Exec(request, id)
    if err != nil {
        return err
    }
    return err
}

func (reg *Reg) ListBStoresDescr() ([]*dscom.BStoreDescr, error) {
    var err error
    request := `
        SELECT id, address, login, pass, state
        FROM bstores
        WHERE;`
    bstores := make([]*dscom.BStoreDescr, 0)
    err = reg.db.Select(&bstores, request)
    if err != nil {
        return bstores, err
    }
    return bstores, err
}
