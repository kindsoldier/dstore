
/*
 * Copyright 2022 Oleg Borodin  <borodin@unix7.org>
 */

package bsapi

import (
    "dstore/dsdescr"
)

const DeleteFileMethod string = "deleteFile"
type DeleteFileParams struct {
    FilePath    string              `msgpack:"filePath"  json:"filePath"`
}

type DeleteFileResult struct {
}

func NewDeleteFileResult() *DeleteFileResult {
    return &DeleteFileResult{}
}

func NewDeleteFileParams() *DeleteFileParams {
    return &DeleteFileParams{}
}

const ListFilesMethod string = "listFiles"

type ListFilesParams struct {
    DirPath     string              `msgpack:"dirPath"  json:"dirPath"`
}

type ListFilesResult struct {
    Files   []*dsdescr.File         `msgpack:"files"    json:"files"`
}

func NewListFilesResult() *ListFilesResult {
    return &ListFilesResult{}
}

func NewListFilesParams() *ListFilesParams {
    return &ListFilesParams{}
}

const LoadFileMethod string = "loadFile"

type LoadFileParams struct {
    FilePath    string              `msgpack:"filePath"  json:"filePath"`
}

type LoadFileResult struct {
}

func NewLoadFileResult() *LoadFileResult {
    return &LoadFileResult{}
}
func NewLoadFileParams() *LoadFileParams {
    return &LoadFileParams{}
}

const SaveFileMethod string = "saveFile"

type SaveFileParams struct {
    FilePath    string              `msgpack:"filePath"  json:"filePath"`
}

type SaveFileResult struct {
}

func NewSaveFileResult() *SaveFileResult {
    return &SaveFileResult{}
}

func NewSaveFileParams() *SaveFileParams {
    return &SaveFileParams{}
}
