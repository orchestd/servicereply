package servicereply

import (
	"bytes"
	"encoding/json"
	"fmt"
	"github.com/orchestd/servicereply/commonError"
	"github.com/orchestd/servicereply/status"
	"github.com/orchestd/servicereply/types"
	"github.com/pkg/errors"
	"runtime"
	"strings"
	"text/template"
)

type ServiceReply interface {
	error
	WithError(error) ServiceReply
	GetError() error
	WithReplyValues(ValuesMap) ServiceReply
	GetReplyValues() ValuesMap
	WithLogMessage(string) ServiceReply
	GetLogMessage() *string
	WithLogValues(ValuesMap) ServiceReply
	GetLogValues() ValuesMap
	GetActionLog() string
	GetSource() string
	GetUserError() string
	GetErrorType() *types.ReplyType
	IsSuccess() bool
}

const statusError = "error"

type BaseServiceError struct {
	source      string
	logMessage  *string
	actionLog   string
	logValues   ValuesMap
	err         error
	errorType   *types.ReplyType
	userMessage string
	extraData   ValuesMap
}
type Message struct {
	Id     string                 `json:"id"`
	Values map[string]interface{} `json:"values"`
}

type BaseResponse struct {
	Status  status.Status `json:"status"`
	Message *Message      `json:"message,omitempty"`
}

type IResponse interface {
	GetMessageId() *string
	GetMessageValues() *map[string]interface{}
	GetStatus() status.Status
}
type Response struct {
	BaseResponse
	Data interface{} `json:"data,omitempty"`
}

func (r Response) GetStatus() status.Status {
	return r.Status
}
func (r Response) GetMessageId() string {
	if r.Message != nil {
		return r.Message.Id
	}
	return ""
}
func (r Response) GetMessageValues() *map[string]interface{} {
	if r.Message != nil {
		return &r.Message.Values
	}
	return nil
}

type ValuesMap map[string]interface{}

var tmplError = "Templating error"

func (se *BaseServiceError) WithError(err error) ServiceReply {
	if se.err != nil {
		se.err = errors.Wrap(se.err, err.Error())
	} else {
		se.err = err
	}
	return se
}

func (se *BaseServiceError) GetError() error {
	return se.err
}
func (se *BaseServiceError) Error() string {
	sr := Response{
		BaseResponse: BaseResponse{
			Status: statusError,
			Message: &Message{
				Id:     se.GetUserError(),
				Values: se.GetReplyValues(),
			},
		},
		Data: nil,
	}
	seBytesArr, _ := json.Marshal(sr)
	return string(seBytesArr)
}
func (se *BaseServiceError) Parse(err string) (Response, error) {
	parsedSr := Response{}
	if err := json.Unmarshal([]byte(err), &parsedSr); err != nil {
		return Response{}, err
	}
	return parsedSr, nil
}
func (se *BaseServiceError) WithReplyValues(extraData ValuesMap) ServiceReply {
	se.extraData = extraData
	return se
}

func (se *BaseServiceError) GetReplyValues() ValuesMap {
	return se.extraData
}
func (se *BaseServiceError) IsSuccess() bool {
	if se == nil || se.GetErrorType() == nil {
		return true
	}
	return *se.GetErrorType() == types.SuccessReplyType
}

func (se *BaseServiceError) WithLogMessage(logMessage string) ServiceReply {
	se.logMessage = &logMessage
	return se
}

func (se *BaseServiceError) GetLogMessage() *string {
	if se.logValues != nil && se.logMessage != nil {
		tmpl, err := template.New("logMessage").Parse(*se.logMessage)
		if err != nil {
			return &tmplError
		}
		buf := &bytes.Buffer{}
		if err := tmpl.Execute(buf, se.logValues); err != nil {
			return &tmplError
		}
		s := buf.String()
		return &s
	}
	return se.logMessage
}

func (se *BaseServiceError) WithLogValues(logValues ValuesMap) ServiceReply {
	se.logValues = logValues
	return se
}

func (se *BaseServiceError) GetLogValues() ValuesMap {
	return se.logValues
}

func (se *BaseServiceError) GetUserError() string {
	return se.userMessage
}

func (se *BaseServiceError) GetErrorType() *types.ReplyType {
	return se.errorType
}

func (se *BaseServiceError) GetActionLog() string {
	return se.actionLog
}
func (se *BaseServiceError) GetSource() string {
	return se.source
}

func NewServiceError(errType *types.ReplyType, err error, userMessage string, runTimeCaller int) ServiceReply {
	sr, ok := err.(ServiceReply)
	if ok {
		return sr
	}

	runTimeCaller += 1
	pc, fn, line, _ := runtime.Caller(runTimeCaller)
	sourceArr := strings.Split(fn, "/")
	if len(sourceArr) >= 2 {
		sourceArr = sourceArr[len(sourceArr)-2:]
	}

	formattedAction := fmt.Sprintf("error in %s", runtime.FuncForPC(pc).Name())
	source := fmt.Sprintf("%s:%d", strings.Join(sourceArr, "/"), line)

	return &BaseServiceError{
		source:      source,
		actionLog:   formattedAction,
		err:         err,
		errorType:   errType,
		userMessage: userMessage,
	}
}

func NewBadRequestError(userMessage string) ServiceReply {
	et := types.BadRequestErrorReplyType
	return NewServiceError(&et, nil, userMessage, 1)
}

type ValidationErrors map[string]string

func NewValidationError(validationErrors ValidationErrors) ServiceReply {
	valueMap := ValuesMap{}
	for field, ve := range validationErrors {
		val := make(map[string]string)
		val["key"] = ve
		valueMap[field] = val
	}
	return NewRejectedReply("validation").WithReplyValues(valueMap)
}

func NewInternalServiceError(err error) ServiceReply {
	et := types.InternalServiceErrorReplyType
	return NewServiceError(&et, err, commonError.InternalServiceError, 1)
}

func NewDbError(err error) ServiceReply {
	et := types.DbErrorReplyType
	return NewServiceError(&et, err, commonError.InternalServiceError, 1)
}

func NewIoError(err error) ServiceReply {
	et := types.IoErrorReplyType
	return NewServiceError(&et, err, commonError.InternalServiceError, 1)
}

func NewNetworkError(err error) ServiceReply {
	et := types.NetworkErrorReplyType
	return NewServiceError(&et, err, commonError.InternalServiceError, 1)
}

func NewServiceAuthError(userMessage string) ServiceReply {
	et := types.ServiceAuthErrorReplyType
	return NewServiceError(&et, nil, userMessage, 1)
}

func NewRejectedReply(userMessage string) ServiceReply {
	et := types.RejectedReplyType
	return NewServiceError(&et, nil, userMessage, 1)
}

func NewNoMatchReply(userMessage string) ServiceReply {
	et := types.NoMatchReplyType
	return NewServiceError(&et, nil, userMessage, 1)
}

func NewMessage(userMessage string) ServiceReply {
	et := types.SuccessReplyType
	return NewServiceError(&et, nil, userMessage, 1)
}

func NewValidationMandatoryRejectedReply(fields []string) ServiceReply {
	return NewBadRequestError("IsMandatory").WithReplyValues(ValuesMap{"fields": fields})
}
func NewWithReplyHeadersValues(values map[string]string) ServiceReply {
	return NewNil().WithReplyValues(map[string]interface{}{"replyHeadersValues": values})
}

func NewNil() ServiceReply {
	return &BaseServiceError{}
}
