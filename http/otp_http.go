package http

import (
	"app/dto"

	"github.com/gin-gonic/gin"
)

// RequestOtp godoc
//
//	@Summary	RequestOtp
//	@Description
//	@Tags		otps
//	@Produce	json
//	@Param		request	body		dto.OtpReq	true	"request"
//	@Success	200		{object}	dto.HTTPResp
//	@Failure	400		{object}	dto.HTTPResp
//	@Failure	500		{object}	dto.HTTPResp
//	@Router		/api/otps/request [post]
func (s *Server) RequestOtp(ctx *gin.Context) {
	req, err := dto.OtpReq{}.Bind(ctx)
	if err != nil {
		abortWithStatusError(ctx, 400, err)
		return
	}
	_, err = s.OtpSvc.GenerateOtp(ctxFromGin(ctx), req.Email)
	if err != nil {
		abortWithStatusError(ctx, 400, err)
		return
	}
}
