
/*
 * Copyright 2022 Oleg Borodin  <borodin@unix7.org>
 */

package fsapi

import (
    "dstore/dscomm/dsdescr"
)

const SaveFileMethod string = "saveFile"

type SaveFileParams struct {
    FilePath    string              `msgpack:"filePath"  json:"filePath"`
}

type SaveFileResult struct {
    File   *dsdescr.File            `msgpack:"file"    json:"file"`
}

func NewSaveFileResult() *SaveFileResult {
    return &SaveFileResult{}
}

func NewSaveFileParams() *SaveFileParams {
    return &SaveFileParams{}
}


const ListFilesMethod string = "listFiles"

type ListFilesParams struct {
    Pattern     string              `msgpack:"pattern"  json:"pattern"`
    Regular     string              `msgpack:"pegular"  json:"regular"`
    GPattern    string              `msgpack:"gPattern"     json:"gPattern"`
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
    File   *dsdescr.File            `msgpack:"file"    json:"file"`
}

func NewLoadFileResult() *LoadFileResult {
    return &LoadFileResult{}
}
func NewLoadFileParams() *LoadFileParams {
    return &LoadFileParams{}
}


const DeleteFileMethod string = "deleteFile"
type DeleteFileParams struct {
    FilePath    string              `msgpack:"filePath"  json:"filePath"`
}

type DeleteFileResult struct {
    File   *dsdescr.File            `msgpack:"file"    json:"file"`
}

func NewDeleteFileResult() *DeleteFileResult {
    return &DeleteFileResult{}
}

func NewDeleteFileParams() *DeleteFileParams {
    return &DeleteFileParams{}
}


const FileStatsMethod string = "fileStats"

type FileStatsParams struct {
    Pattern     string              `msgpack:"pattern"      json:"pattern"`
    Regular     string              `msgpack:"pegular"      json:"regular"`
    GPattern    string              `msgpack:"gPattern"     json:"gPattern"`
}

type FileStatsResult struct {
    Count      int64                `msgpack:"count"    json:"count"`
    Usage      int64                `msgpack:"usage"    json:"usage"`
}

func NewFileStatsResult() *FileStatsResult {
    return &FileStatsResult{}
}

func NewFileStatsParams() *FileStatsParams {
    return &FileStatsParams{}
}


const EraseFilesMethod string = "eraseFiles"

type EraseFilesParams struct {
    Pattern     string              `msgpack:"pattern"  json:"pattern"`
    Regular     string              `msgpack:"pegular"  json:"regular"`
    GPattern    string              `msgpack:"gPattern" json:"gPattern"`
    Erase       bool                `msgpack:"erase"    json:"erase"`
}

type EraseFilesResult struct {
    Files   []*dsdescr.File         `msgpack:"files"    json:"files"`
}

func NewEraseFilesResult() *EraseFilesResult {
    return &EraseFilesResult{}
}

func NewEraseFilesParams() *EraseFilesParams {
    return &EraseFilesParams{}
}


