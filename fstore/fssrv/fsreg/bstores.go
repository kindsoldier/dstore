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
        bstore_id   INTEGER GENERATED ALWAYS AS IDENTITY (START 1 CYCLE),
        address     TEXT,
        port        TEXT,
        login       TEXT,
        pass        TEXT,
        state       TEXT
    );
    DROP INDEX IF EXISTS bstore_idx;
    CREATE UNIQUE INDEX IF NOT EXISTS bstore_idx
        ON bstores (address, port);
    `

func (reg *Reg) AddBStoreDescr(address, port, login, pass, state string) (int64, error) {
    var err error
    request := `
        INSERT INTO bstores(address, port, login, pass, state)
        VALUES ($1, $2, $3, $4, $5)
        RETURNING bstore_id;`
    var storeId int64
    err = reg.db.Get(&storeId, request, address, port, login, pass, state)
    if err != nil {
        return storeId, err
    }
    return storeId, err
}

func (reg *Reg) UpdateBStoreDescr(address, port, login, pass, state string) error {
    var err error
    request := `
        UPDATE bstores
        SET address = $1, port = $2, login = $3, pass = $4, state = $5
        WHERE address = $1 AND port = $2;`
    _, err = reg.db.Exec(request, address, port, login, pass, state)
    if err != nil {
        return err
    }
    return err
}

func (reg *Reg) RenewBStoreDescr(descr *dscom.BStoreDescr) error {
    var err error
    request := `
        UPDATE bstores
        SET address = $1, port = $2, login = $3, pass = $4, state = $5
        WHERE address = $1 AND port = $2;`
    _, err = reg.db.Exec(request, descr.Address, descr.Port, descr.Login, descr.Pass, descr.State)
    if err != nil {
        return err
    }
    return err
}

func (reg *Reg) BStoreDescrExists(address, port string) (bool, error) {
    var err error
    var exists bool
    request := `
        SELECT count(bstore_id) AS count
        FROM bstores
        WHERE address = $1 AND port = $2
        LIMIT 1;`
    var count int64
    err = reg.db.Get(&count, request, address, port)
    if err != nil {
        return exists, err
    }
    if count > 0 {
        exists = true
    }
    return exists, err
}

func (reg *Reg) GetBStoreDescr(address, port string) (*dscom.BStoreDescr, error) {
    var err error
    request := `
        SELECT bstore_id, address, port, login, pass, state
        FROM bstores
        WHERE address = $1 AND port = $2
        LIMIT 1;`
    bstore := dscom.NewBStoreDescr()
    err = reg.db.Get(bstore, request, address, port)
    if err != nil {
        return bstore, err
    }
    return bstore, err
}

func (reg *Reg) DeleteBStoreDescr(address, port string) error {
    var err error
    request := `
        DELETE FROM bstores
        WHERE address = $1 AND port = $2;`
    _, err = reg.db.Exec(request, address, port)
    if err != nil {
        return err
    }
    return err
}

func (reg *Reg) ListBStoreDescrs() ([]*dscom.BStoreDescr, error) {
    var err error
    request := `
        SELECT bstore_id, address, port, login, pass, state
        FROM bstores;`
    bstores := make([]*dscom.BStoreDescr, 0)
    err = reg.db.Select(&bstores, request)
    if err != nil {
        return bstores, err
    }
    return bstores, err
}
