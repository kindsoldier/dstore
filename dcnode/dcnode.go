package main.go

func main() {
}

type Server struct {

}

func NewServer() *Server {
    var server Server
    return &server
}

func (server *Server) Config() error {
    var err error
    return err
}

func (server *Server) Fork() error {
    var err error
    return err

}

func (server *Server) Start() error {
    var err error
    return err
}
