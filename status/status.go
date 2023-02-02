package status

import (
	"github.com/orchestd/servicereply/types"
)

var statusMap = map[types.ReplyType]Status{
	types.BadRequestErrorReplyType:      InvalidStatus,
	types.DbErrorReplyType:              ErrorStatus,
	types.InternalServiceErrorReplyType: ErrorStatus,
	types.IoErrorReplyType:              ErrorStatus,
	types.NetworkErrorReplyType:         ErrorStatus,

	types.ServiceAuthErrorReplyType: UnauthorizedStatus,

	types.RejectedReplyType: RejectedStatus,
	types.NoMatchReplyType:  NoMatchStatus,
	types.SuccessReplyType:  SuccessStatus,
}

type Status string

const (
	UnauthorizedStatus Status = "unauthorized"
	RejectedStatus     Status = "rejected"
	NoMatchStatus      Status = "noMatch"
	ErrorStatus        Status = "error"
	InvalidStatus      Status = "invalid"
	SuccessStatus      Status = "success"
)

func GetStatus(et *types.ReplyType) Status {
	if et == nil {
		return SuccessStatus
	}
	return statusMap[*et]
}

var typesMap = map[Status]types.ReplyType {
	InvalidStatus : types.BadRequestErrorReplyType,

	ErrorStatus : types.InternalServiceErrorReplyType,

	UnauthorizedStatus : types.ServiceAuthErrorReplyType,

	RejectedStatus : types.RejectedReplyType,
	NoMatchStatus : types.NoMatchReplyType,
	SuccessStatus : types.SuccessReplyType,
}

func GetTypeByStatus(s Status) types.ReplyType {
	return typesMap[s]
}
