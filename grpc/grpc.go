package grpc

import (
	"github.com/orchestd/servicereply/types"
	"google.golang.org/grpc/codes"
)

var grpcErrors = map[types.ReplyType]codes.Code{
	types.BadRequestErrorReplyType:      codes.InvalidArgument,
	types.DbErrorReplyType:              codes.Internal,
	types.InternalServiceErrorReplyType: codes.Internal,
	types.IoErrorReplyType:              codes.Internal,
	types.NetworkErrorReplyType:         codes.Internal,

	types.ServiceAuthErrorReplyType: codes.Unauthenticated,

	types.RejectedReplyType: codes.OK,
	types.NoMatchReplyType:  codes.OK,
	types.SuccessReplyType:  codes.OK,

}

func GetGrpcCode(et types.ReplyType) codes.Code{
	return grpcErrors[et]
}
