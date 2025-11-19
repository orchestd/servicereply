package types

const (
	//500
	DbErrorReplyType              ReplyType = "db"
	NetworkErrorReplyType         ReplyType = "network"
	IoErrorReplyType              ReplyType = "io"
	InternalServiceErrorReplyType ReplyType = "internalServiceError"
	BadRequestErrorReplyType      ReplyType = "badRequest"

	//401
	ServiceAuthErrorReplyType ReplyType = "serviceAuth"

	//503
	ServiceUnavailableErrorReplyType ReplyType = "serviceUnavailable"

	//Logic
	NoMatchReplyType  ReplyType = "noMatch"
	SuccessReplyType  ReplyType = "success"
	RejectedReplyType ReplyType = "rejected"
)

type ReplyType string
