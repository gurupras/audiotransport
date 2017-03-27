package audiotransport

type ServerInterface interface {
	Bind(addr string) (err error)
	Listen(callback func(transport Transport)) (err error)
	Interrupt()
}

type Server struct {
	signalChan chan interface{}
}

func (server *Server) Interrupt() {
	server.signalChan <- struct{}{}
}
