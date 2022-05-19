package dcfile

type File struct {
    batchs []Batch
}

func NewFile() *File {
    var file File
    file.batchs = make([]Batch, 0)
    return &file
}
