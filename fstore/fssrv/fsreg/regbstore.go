/*
 * Copyright 2022 Oleg Borodin  <borodin@unix7.org>
 */
package fsreg

import (
    "ndstore/dscom"
    "ndstore/dserr"
)

const bstoreSchema = `
    DROP TABLE IF EXISTS fs_bstores;
    CREATE TABLE IF NOT EXISTS fs_bstores (
        bstore_id   INTEGER GENERATED ALWAYS AS IDENTITY (START 1 CYCLE),
        address     TEXT,
        port        TEXT,
        login       TEXT,
        pass        TEXT,
        state       TEXT
    );
    DROP INDEX IF EXISTS fs_bstore_idx;
    CREATE UNIQUE INDEX IF NOT EXISTS fs_bstore_idx
        ON fs_bstores(address, port);
    `

func (reg *Reg) AddBStoreDescr(descr *dscom.BStoreDescr) (int64, error) {
    var err error
    request := `
        INSERT INTO fs_bstores(address, port, login, pass, state)
        VALUES ($1, $2, $3, $4, $5)
        RETURNING bstore_id;`
    var storeId int64
    err = reg.db.Get(&storeId, request, descr.Address, descr.Port, descr.Login, descr.Pass, descr.State)
    if err != nil {
        return storeId, dserr.Err(err)
    }
    return storeId, dserr.Err(err)
}

func (reg *Reg) UpdateBStoreDescr(descr *dscom.BStoreDescr) error {
    var err error
    request := `
        UPDATE fs_bstores
        SET address = $1, port = $2, login = $3, pass = $4, state = $5
        WHERE address = $1 AND port = $2;`
    _, err = reg.db.Exec(request, descr.Address, descr.Port, descr.Login, descr.Pass, descr.State)
    if err != nil {
        return dserr.Err(err)
    }
    return dserr.Err(err)
}

func (reg *Reg) BStoreDescrExists(address, port string) (bool, error) {
    var err error
    var exists bool
    request := `
        SELECT count(bstore_id) AS count
        FROM fs_bstores
        WHERE address = $1 AND port = $2
        LIMIT 1;`
    var count int64
    err = reg.db.Get(&count, request, address, port)
    if err != nil {
        return exists, dserr.Err(err)
    }
    if count > 0 {
        exists = true
    }
    return exists, dserr.Err(err)
}

func (reg *Reg) GetBStoreDescr(address, port string) (*dscom.BStoreDescr, error) {
    var err error
    request := `
        SELECT bstore_id, address, port, login, pass, state
        FROM fs_bstores
        WHERE address = $1 AND port = $2
        LIMIT 1;`
    bstore := dscom.NewBStoreDescr()
    err = reg.db.Get(bstore, request, address, port)
    if err != nil {
        return bstore, dserr.Err(err)
    }
    return bstore, dserr.Err(err)
}

func (reg *Reg) EraseBStoreDescr(address, port string) error {
    var err error
    request := `
        DELETE FROM fs_bstores
        WHERE address = $1 AND port = $2;`
    _, err = reg.db.Exec(request, address, port)
    if err != nil {
        return dserr.Err(err)
    }
    return dserr.Err(err)
}

func (reg *Reg) ListBStoreDescrs() ([]*dscom.BStoreDescr, error) {
    var err error
    request := `
        SELECT bstore_id, address, port, login, pass, state
        FROM fs_bstores;`
    bstores := make([]*dscom.BStoreDescr, 0)
    err = reg.db.Select(&bstores, request)
    if err != nil {
        return bstores, dserr.Err(err)
    }
    return bstores, dserr.Err(err)
}