package http

import (
	"app/dto"
	"app/entity"
	"app/errors"
	"app/pkg/apperror"
	"app/pkg/logger"
	"app/pkg/utils"
	"net/http"

	"github.com/gin-gonic/gin"
)

// Login godoc
//
//	@Summary	Login
//	@Description
//	@Tags		authentications
//	@Accept		json
//	@Produce	json
//	@Param		request	body		dto.UserLoginReq	true	"request"
//	@Success	200		{object}	dto.AccessToken
//	@Failure	400		{object}	dto.HTTPResp
//	@Failure	500		{object}	dto.HTTPResp
//	@Router		/api/auth/login [post]
func (s *Server) Login(ctx *gin.Context) {
	req, err := dto.UserLoginReq{}.Bind(ctx)
	if err != nil {
		abortWithStatusError(ctx, 400, err)
		return
	}
	accessToken, err := s.UserSvc.Login(ctxFromGin(ctx), req.ToUser(ctxFromGin(ctx)))
	if errors.UserNotFound().Is(err) {
		logger.For(ctxFromGin(ctx)).Error(apperror.As(err).StackTrace())
		abortWithStatusError(ctx, 400, errors.UserNotRegister())
		return
	}
	if err != nil {
		abortWithStatusError(ctx, 400, err)
		return
	}
	ctx.AbortWithStatusJSON(http.StatusOK, dto.AccessToken{
		AccessToken: accessToken,
	})
}

// Logout godoc
//
//	@Summary	Logout
//	@Description
//	@Tags		authentications
//	@Produce	json
//	@Param		Authorization	header	string	true	"Bearer token"
//	@Success	200
//	@Failure	400	{object}	dto.HTTPResp
//	@Failure	500	{object}	dto.HTTPResp
//	@Router		/api/auth/logout [get]
func (s *Server) Logout(ctx *gin.Context) {
	bearerToken, ok := utils.GetBearerAuth(ctx)
	if !ok {
		abortWithStatusError(ctx, 400, apperror.NewError(errors.CodeTokenError, "empty token"))
		return
	}
	err := s.UserSvc.UserLogout(ctxFromGin(ctx), bearerToken)
	if err != nil {
		abortWithStatusError(ctx, 400, err)
		return
	}
	ctx.AbortWithStatus(http.StatusOK)
}

// GetProfile godoc
//
//	@Summary	GetProfile
//	@Description
//	@Tags		users
//	@Produce	json
//	@Param		Authorization	header		string	true	"Bearer token"
//	@Success	200				{object}	dto.UserResp
//	@Failure	400				{object}	dto.HTTPResp
//	@Failure	500				{object}	dto.HTTPResp
//	@Router		/api/user/profile [get]
func (s *Server) GetProfile(ctx *gin.Context) {
	user := entity.GetUserFromContext(ctxFromGin(ctx))
	ctx.AbortWithStatusJSON(http.StatusOK, dto.UserResp{}.FromUser(&user))
}

// GetUserList godoc
//
//	@Summary	GetUserList
//	@Description
//	@Tags		users
//	@Produce	json
//	@Param		Authorization	header		string	true	"Bearer token"
//	@Param		page			query		int		false	"page of paging"
//	@Param		page_size		query		int		false	"size of page of paging"
//	@Param		sort			query		string	false	"name of field need to sort"
//	@Param		sort_type		query		string	false	"sort desc or asc"
//	@Param		search			query		string	false	"keyword to search in model"
//	@Success	200				{object}	dto.HTTPResp
//	@Failure	400				{object}	dto.HTTPResp
//	@Failure	500				{object}	dto.HTTPResp
//	@Router		/api/users [get]
func (s *Server) GetUserList(ctx *gin.Context) {
	params, err := dto.QueryParams{}.Bind(ctx)
	if err != nil {
		abortWithStatusError(ctx, 400, err)
		return
	}
	err = params.Validate(dto.UserResp{})
	if err != nil {
		abortWithStatusError(ctx, 400, err)
		return
	}
	users, total, err := s.UserSvc.GetUserList(ctxFromGin(ctx), params.ToRepoQueryParams())
	if err != nil {
		abortWithStatusError(ctx, 400, err)
		return
	}
	var list = []*dto.UserResp{}
	for _, u := range users {
		list = append(list, dto.UserResp{}.FromUser(u))
	}
	res := dto.UserListResp{
		Total:    total,
		Page:     params.Page,
		PageSize: params.PageSize,
		List:     list,
	}
	ctx.AbortWithStatusJSON(200, res)
}

// GetUser godoc
//
//	@Summary	GetUser
//	@Description
//	@Tags		users
//	@Produce	json
//	@Param		Authorization	header		string	true	"Bearer token"
//	@Success	200				{object}	dto.UserResp
//	@Failure	400				{object}	dto.HTTPResp
//	@Failure	500				{object}	dto.HTTPResp
//	@Router		/api/users/{id} [get]
func (s *Server) GetUser(ctx *gin.Context) {
	id := ctx.Param("id")
	user, err := s.UserSvc.GetUser(ctxFromGin(ctx), &entity.User{
		ID: id,
	})
	if err != nil {
		abortWithStatusError(ctx, 400, err)
		return
	}
	ctx.AbortWithStatusJSON(200, dto.UserResp{}.FromUser(user))
}

// CreateUser godoc
//
//	@Summary	CreateUser
//	@Description
//	@Tags		users
//	@Produce	json
//	@Param		request	body		dto.UserCreateReq	true	"request"
//	@Success	200		{object}	dto.HTTPResp
//	@Failure	400		{object}	dto.HTTPResp
//	@Failure	500		{object}	dto.HTTPResp
//	@Router		/api/users [post]
func (s *Server) CreateUser(ctx *gin.Context) {
	req, err := dto.UserCreateReq{}.Bind(ctx)
	if err != nil {
		abortWithStatusError(ctx, 400, err)
		return
	}
	_, err = s.UserSvc.CreateUser(ctxFromGin(ctx), req.ToUser(ctxFromGin(ctx)))
	if err != nil {
		abortWithStatusError(ctx, 400, err)
		return
	}
}

// ActiveUser godoc
//
//	@Summary	ActiveUser
//	@Description
//	@Tags		users
//	@Produce	json
//	@Param		request	body		dto.UserActiveReq	true	"request"
//	@Success	200		{object}	dto.HTTPResp
//	@Failure	400		{object}	dto.HTTPResp
//	@Failure	500		{object}	dto.HTTPResp
//	@Router		/api/users/active [put]
func (s *Server) ActiveUser(ctx *gin.Context) {
	req, err := dto.UserActiveReq{}.Bind(ctx)
	if err != nil {
		abortWithStatusError(ctx, 400, err)
		return
	}
	_, err = s.UserSvc.ActiveUser(ctxFromGin(ctx), req.ToUser(ctxFromGin(ctx)))
	if err != nil {
		abortWithStatusError(ctx, 400, err)
		return
	}
}

// ResetPassword godoc
//
//	@Summary	ResetPassword
//	@Description
//	@Tags		users
//	@Produce	json
//	@Param		request	body		dto.UserResetPasswordReq	true	"request"
//	@Success	200		{object}	dto.HTTPResp
//	@Failure	400		{object}	dto.HTTPResp
//	@Failure	500		{object}	dto.HTTPResp
//	@Router		/api/users/reset-password [put]
func (s *Server) ResetPassword(ctx *gin.Context) {
	req, err := dto.UserResetPasswordReq{}.Bind(ctx)
	if err != nil {
		abortWithStatusError(ctx, 400, err)
		return
	}
	_, err = s.UserSvc.ResetPassword(ctxFromGin(ctx), req.ToUser(ctxFromGin(ctx)))
	if err != nil {
		abortWithStatusError(ctx, 400, err)
		return
	}
}

// UpdateUser godoc
//
//	@Summary	UpdateUser
//	@Description
//	@Tags		users
//	@Produce	json
//	@Param		Authorization	header		string				true	"Bearer token"
//	@Param		request			body		dto.UserUpdateReq	true	"request"
//	@Success	200				{object}	dto.HTTPResp
//	@Failure	400				{object}	dto.HTTPResp
//	@Failure	500				{object}	dto.HTTPResp
//	@Router		/api/users [put]
func (s *Server) UpdateUser(ctx *gin.Context) {
	req, err := dto.UserUpdateReq{}.Bind(ctx)
	if err != nil {
		abortWithStatusError(ctx, 400, err)
		return
	}
	err = req.Validate()
	if err != nil {
		abortWithStatusError(ctx, 400, err)
		return
	}
	_, err = s.UserSvc.UpdateUser(ctxFromGin(ctx), req.ToUser(ctxFromGin(ctx)))
	if err != nil {
		abortWithStatusError(ctx, 400, err)
		return
	}
}

// DeleteUser godoc
//
//	@Summary	DeleteUser
//	@Description
//	@Tags		users
//	@Produce	json
//	@Param		Authorization	header		string				true	"Bearer token"
//	@Param		request			body		dto.UserDeleteReq	true	"request"
//	@Success	200				{object}	dto.HTTPResp
//	@Failure	400				{object}	dto.HTTPResp
//	@Failure	500				{object}	dto.HTTPResp
//	@Router		/api/users [delete]
func (s *Server) DeleteUser(ctx *gin.Context) {
	req, err := dto.UserDeleteReq{}.Bind(ctx)
	if err != nil {
		abortWithStatusError(ctx, 400, err)
		return
	}
	_, err = s.UserSvc.DeleteUser(ctxFromGin(ctx), req.ToUser(ctxFromGin(ctx)))
	if err != nil {
		abortWithStatusError(ctx, 400, err)
		return
	}
}
