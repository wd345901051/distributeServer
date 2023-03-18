package server

type ServerStatus int

const (
	Master   ServerStatus = 1
	Follower ServerStatus = 2
)

type Server struct {
	Address string
	Status  ServerStatus
}

func NewServer() *Server {
	s := &Server{
		Address: "0.0.0.0:8686",
		Status:  1,
	}
	return s
}

func (s *Server) Run() {

}
