package gsrv

import (
	"github.com/IamStubborN/filmtracker/router"
	"github.com/gin-gonic/gin"
)

type Server struct {
	Router *gin.Engine
}

func CreateServer() *Server {
	server := Server{}
	server.Router = router.CreateRouter()
	return &server
}

func (s *Server) Run() error {
	if err := s.Router.Run(); err != nil {
		return err
	}
	return nil
}
