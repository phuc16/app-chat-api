package http

import (
	"app/entity"
	"net/http"

	"app/dto"

	"github.com/gin-gonic/gin"
)

func (s *Server) ServeWs(ctx *gin.Context) {
	user := entity.GetUserFromContext(ctxFromGin(ctx))
	s.ConversationSvc.ServeWs(ctxFromGin(ctx), user.ID, ctx.Writer, ctx.Request)
}

// GetConversation godoc
//
//	@Summary	GetConversation
//	@Description
//	@Tags		conversations
//	@Produce	json
//	@Param		Authorization	header		string	true	"Bearer token"
//	@Success	200				{object}	dto.ConversationInfoResp
//	@Failure	400				{object}	dto.HTTPResp
//	@Failure	500				{object}	dto.HTTPResp
//	@Router		/api/conversations/{id} [get]
func (s *Server) GetConversation(ctx *gin.Context) {
	id := ctx.Param("id")
	conversation, err := s.ConversationSvc.GetConversation(ctxFromGin(ctx), &entity.Conversation{
		ID: id,
	})
	if err != nil {
		abortWithStatusError(ctx, 400, err)
		return
	}
	ctx.AbortWithStatusJSON(200, dto.ConversationInfoResp{}.FromConversation(conversation))
}

// GetChatList godoc
//
//	@Summary	GetChatList
//	@Description
//	@Tags		conversations
//	@Produce	json
//	@Param		Authorization	header		string	true	"Bearer token"
//	@Param		page			query		int		false	"page of paging"
//	@Param		page_size		query		int		false	"size of page of paging"
//	@Param		sort_type		query		string	false	"sort desc or asc"
//	@Param		search			query		string	false	"keyword to search in model, search by msg"
//	@Success	200				{object}	dto.ChatListResp
//	@Failure	400				{object}	dto.HTTPResp
//	@Failure	500				{object}	dto.HTTPResp
//	@Router		/api/conversations/{id}/chats [get]
func (s *Server) GetChatList(ctx *gin.Context) {
	id := ctx.Param("id")
	params, err := dto.QueryParams{}.Bind(ctx)
	if err != nil {
		abortWithStatusError(ctx, 400, err)
		return
	}
	err = params.Validate(dto.ChatResp{})
	if err != nil {
		abortWithStatusError(ctx, 400, err)
		return
	}
	conversations, total, err := s.ConversationSvc.GetChatList(ctxFromGin(ctx), id, params.ToRepoQueryParams())
	if err != nil {
		abortWithStatusError(ctx, 400, err)
		return
	}
	var list = []*dto.ChatResp{}
	for _, u := range conversations {
		list = append(list, dto.ChatResp{}.FromChat(u))
	}
	res := dto.ChatListResp{
		Total:    total,
		Page:     params.Page,
		PageSize: params.PageSize,
		List:     list,
	}
	ctx.AbortWithStatusJSON(200, res)
}

// UpdateMessage godoc
//
//	@Summary	UpdateMessage
//	@Description
//	@Tags		conversations
//	@Produce	json
//	@Param		Authorization	header		string					true	"Bearer token"
//	@Param		request			body		dto.MessageUpdateReq	true	"request"
//	@Success	200				{object}	dto.HTTPResp
//	@Failure	400				{object}	dto.HTTPResp
//	@Failure	500				{object}	dto.HTTPResp
//	@Router		/api/conversations/{id}/chats [put]
func (s *Server) UpdateMessage(ctx *gin.Context) {
	id := ctx.Param("id")
	req, err := dto.MessageUpdateReq{
		ConversationID: id,
	}.Bind(ctx)
	if err != nil {
		abortWithStatusError(ctx, 400, err)
		return
	}
	_, err = s.ConversationSvc.UpdateMessage(ctxFromGin(ctx), req.ToConversation(ctxFromGin(ctx)))
	if err != nil {
		abortWithStatusError(ctx, 400, err)
		return
	}
	ctx.AbortWithStatusJSON(http.StatusOK, map[string]string{"msg": "ok"})
}
