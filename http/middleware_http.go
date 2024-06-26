package http

import (
	"app/entity"
	"app/errors"
	"app/pkg/apperror"
	"app/pkg/trace"
	"app/pkg/utils"

	"github.com/gin-gonic/gin"
)

func (s *Server) Authenticate(ctx *gin.Context) {
	context, span := trace.Tracer().Start(ctxFromGin(ctx), utils.GetCurrentFuncName())
	defer span.End()

	bearerToken, ok := utils.GetBearerAuth(ctx)
	if !ok {
		abortWithStatusError(ctx, 401, apperror.NewError(errors.CodeTokenError, "empty token"))
		return
	}
	user, err := s.UserSvc.Authenticate(context, bearerToken)
	if err != nil {
		abortWithStatusError(ctx, 401, err)
		return
	}
	user.SetToContext(ctx)
	ctx.Next()
}

func (s *Server) CheckPermission(role string) func(ctx *gin.Context){
	return func(ctx *gin.Context) {
		context, span := trace.Tracer().Start(ctxFromGin(ctx), utils.GetCurrentFuncName())
		defer span.End()
	
		if role != entity.GetUserFromContext(context).Role {
			abortWithStatusError(ctx, 403, errors.PermissionDenied())
			return
		}
	
		ctx.Next()
	}
}
