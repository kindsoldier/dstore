/*
 * Copyright 2022 Oleg Borodin  <borodin@unix7.org>
 */

package dsdescr

import (
    "encoding/json"
)

const UStateEnabled     string  = "enabled"
const UStateDisabled    string  = "disabled"

const URoleAdmin        string  = "admin"
const URoleUser         string  = "user"

type User struct {
    Login       string      `json:"login"`
    Pass        string      `json:"pass"`
    Role        string      `json:"role"`
    State       string      `json:"state"`
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

type Entry struct {
    FilePath    string      `json:"filePath"`
    FileId      int64       `json:"fileId"`
    CreatedAt   int64       `json:"updatedAt"`
    UpdatedAt   int64       `json:"createdAt"`
}

func NewEntry() *Entry {
    var descr Entry
    return &descr
}

func UnpackEntry(descrBin []byte) (*Entry, error) {
    var err error
    var descr Entry
    err = json.Unmarshal(descrBin, &descr)
    return &descr, err
}

func (descr *Entry) Pack() ([]byte, error) {
    var err error
    descrBin, err := json.Marshal(descr)
    return descrBin, err
}


type File struct {
    FileId      int64       `json:"fileId"`
    BatchCount  int64       `json:"batchCount"`
    BatchSize   int64       `json:"batchSize"`
    BlockSize   int64       `json:"blockSize"`
    DataSize    int64       `json:"dataSize"`
    CreatedAt   int64       `json:"createdAt"`
    UpdatedAt   int64       `json:"updatedAt"`
}

func NewFile() *File {
    var descr File
    return &descr
}

func UnpackFile(descrBin []byte) (*File, error) {
    var err error
    var descr File
    err = json.Unmarshal(descrBin, &descr)
    return &descr, err
}

func (descr *File) Pack() ([]byte, error) {
    var err error
    descrBin, err := json.Marshal(descr)
    return descrBin, err
}



type Block struct {
    BlockId     int64       `json:"blockId"`
    BatchId     int64       `json:"batchId"`
    FileId      int64       `json:"fileId"`
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
    BatchId     int64       `json:"batchId"`
    FileId      int64       `json:"fileId"`
    BatchSize   int64       `json:"batchSize"`
    BlockSize   int64       `json:"blockSize"`
    CreatedAt   int64       `json:"createdAt"`
    UpdatedAt   int64       `json:"updatedAt"`
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
