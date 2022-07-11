/*
 * Copyright 2022 Oleg Borodin  <borodin@unix7.org>
 */

package dscom

import (
    "encoding/json"
)

type Block struct {
    BlockSize   int64       `json:"blockSize"`
    DataSize    int64       `json:"dataSize"`
    CreatedAt   int64       `json:"createdAt"`
    UpdatedAt   int64       `json:"updatedAt"`
    FilePath    string      `json:"filePath"`
}

func NewBlock() *Block {
    var descr Block
    return &descr
}

func UnpackBlock(descrBin []byte) (*Block, error) {
    var err error
    var descr Block
    err = json.Unmarshal(descrBin, &descr)
    return &descr, err
}

func (descr *Block) Pack() ([]byte, error) {
    var err error
    descrBin, err := json.Marshal(descr)
    return descrBin, err
}

type Batch struct {
    BatchSize   int64
    BlockSize   int64
    CreatedAt   int64
    UpdatedAt   int64
}

func NewBatch() *Batch {
    var descr Batch
    return &descr
}

func UnpackBatch(descrBin []byte) (*Batch, error) {
    var err error
    var descr Batch
    err = json.Unmarshal(descrBin, &descr)
    return &descr, err
}

func (descr *Batch) Pack() ([]byte, error) {
    var err error
    descrBin, err := json.Marshal(descr)
    return descrBin, err
}



type User struct {
    Login       string      `json:"login"`
    Passw       string      `json:"passw"`
    CreatedAt   int64       `json:"updatedAt"`
    UpdatedAt   int64       `json:"createdAt"`
}

func NewUser() *User {
    var descr User
    return &descr
}

func UnpackUser(descrBin []byte) (*User, error) {
    var err error
    var descr User
    err = json.Unmarshal(descrBin, &descr)
    return &descr, err
}

func (descr *User) Pack() ([]byte, error) {
    var err error
    descrBin, err := json.Marshal(descr)
    return descrBin, err
}

type Alloc struct {
    TopId   int64           `json:"topId"`
    FreeIds []int64         `json:"freeIds"`
}

func NewAlloc() *Alloc {
    var descr Alloc
    return &descr
}

func UnpackAlloc(descrBin []byte) (*Alloc, error) {
    var err error
    var descr Alloc
    err = json.Unmarshal(descrBin, &descr)
    return &descr, err
}

func (descr *Alloc) Pack() ([]byte, error) {
    var err error
    descrBin, err := json.Marshal(descr)
    return descrBin, err
}
