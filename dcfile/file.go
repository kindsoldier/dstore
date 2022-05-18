
package dcfile

//**********************************************************************************************
type File struct {
    batchs []Batch
}

func NewFile() *File {
    var file File
    file.batchs = make([]Batch, 0)
    return &file
}

func (file *File) Write(data []byte) (int, error) {
    var err error
    var written int
    return written, err
}

func (file *File) Read(data []byte) (int, error) {
    var err error
    var read int
    return read, err
}
