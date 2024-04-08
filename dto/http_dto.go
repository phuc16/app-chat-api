package dto

import (
	"app/errors"
	"app/pkg/apperror"
	"app/pkg/utils"
	"app/repository"
	"fmt"
	"regexp"

	"github.com/gin-gonic/gin"
	"github.com/go-playground/validator/v10"
)

type HTTPResp struct {
	Code int    `json:"code"`
	Msg  string `json:"msg"`
}

func (r HTTPResp) FromErr(err *apperror.Error) *HTTPResp {
	return &HTTPResp{
		Code: err.Code,
		Msg:  err.Error(),
	}
}

type QueryParams struct {
	Page     int            `form:"page" json:"page"`
	PageSize int            `form:"page_size" json:"page_size"`
	Sort     string         `form:"sort" json:"sort"`
	SortType string         `form:"sort_type" json:"sort_type"`
	Search   string         `form:"search" json:"search"`
	Filter   map[string]any `form:"filter" json:"filter"`
}

func (p QueryParams) Bind(ctx *gin.Context) (*QueryParams, error) {
	err := ctx.ShouldBindQuery(&p)
	if err != nil {
		return nil, apperror.NewError(errors.CodeUnknownError, validationErrorToText(err))
	}
	p.Filter = map[string]any{}
	filter := ctx.QueryMap("filter")
	for k, v := range filter {
		p.Filter[k] = v
	}
	return &p, nil
}

func (p *QueryParams) Validate(dto interface{}) error {
	if p.PageSize < 1 || p.PageSize > 100 {
		p.PageSize = 100
	}
	if p.Page < 1 {
		p.Page = 1
	}
	tags := utils.GetTagList(dto)
	if p.Sort != "" && !utils.IsExist(tags, p.Sort) {
		return apperror.NewError(errors.CodeUnknownError, fmt.Sprintf("invalid sort field %s", p.Sort))
	}
	if p.SortType != "" && p.SortType != "asc" && p.SortType != "desc" {
		return apperror.NewError(errors.CodeUnknownError, "sort type must be either asc or desc")
	}
	if p.SortType == "" {
		p.SortType = "asc"
	}
	if p.Search != "" {
		p.Search = regexp.QuoteMeta(p.Search)
	}
	if p.Filter == nil {
		p.Filter = map[string]any{}
	}
	return nil
}

func (p *QueryParams) ToRepoQueryParams() *repository.QueryParams {
	sortType := 1
	if p.SortType == "desc" {
		sortType = -1
	}
	return &repository.QueryParams{
		Filter:    p.Filter,
		Limit:     int64(p.PageSize),
		Skip:      int64((p.Page - 1) * p.PageSize),
		SortField: p.Sort,
		SortType:  sortType,
		Search:    p.Search,
	}
}

func validationErrorToText(err error) string {
	if err, ok := err.(validator.ValidationErrors); ok {
		for _, e := range err {
			switch e.Tag() {
			case "required":
				return fmt.Sprintf("%s is required", e.Field())
			case "max":
				return fmt.Sprintf("%s cannot be longer than %s", e.Field(), e.Param())
			case "min":
				return fmt.Sprintf("%s must be longer than %s", e.Field(), e.Param())
			case "email":
				return fmt.Sprintf("Invalid email format")
			case "len":
				return fmt.Sprintf("%s must be %s characters long", e.Field(), e.Param())
			}
			return fmt.Sprintf("%s is not valid", e.Field())
		}
	}
	return err.Error()
}
