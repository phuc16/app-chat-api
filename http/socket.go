package http

import (
	"app/entity"

	"github.com/gin-gonic/gin"
)

func (s *Server) ServeWs(ctx *gin.Context) {
	user := entity.GetUserFromContext(ctxFromGin(ctx))
	s.SocketSvc.ServeWs(ctxFromGin(ctx), user.Name, ctx.Writer, ctx.Request)
}
