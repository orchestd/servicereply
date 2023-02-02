package http

import (
	"github.com/orchestd/servicereply/types"
	"net/http"
)

var httpErrors = map[types.ReplyType]int{
	types.BadRequestErrorReplyType:      http.StatusBadRequest,
	types.DbErrorReplyType:              http.StatusInternalServerError,
	types.InternalServiceErrorReplyType: http.StatusInternalServerError,
	types.IoErrorReplyType:              http.StatusInternalServerError,
	types.NetworkErrorReplyType:         http.StatusInternalServerError,

	types.ServiceAuthErrorReplyType: http.StatusUnauthorized,

	types.RejectedReplyType: http.StatusOK,
	types.NoMatchReplyType:  http.StatusOK,
	types.SuccessReplyType:  http.StatusOK,
}

func GetHttpCode(et *types.ReplyType) int {
	if et == nil {
		return http.StatusOK
	}
	return httpErrors[*et]
}
