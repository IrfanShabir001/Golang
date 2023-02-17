package shared

import (
	"context"
	"fmt"
	"io/ioutil"
	"net/http"
	"reflect"
	"time"

	"github.com/go-zoo/bone"
	"github.com/satori/go.uuid"
	"log"
)

type RequestId struct{}
type ResourceId struct{}
type RequestTimestamp struct{}
type RequestType struct{}
type UserId struct{}
type UserRoles struct{}

// Http Request
type HttpWebRequest struct{ Req *http.Request }

func (hwr HttpWebRequest) Raw() *http.Request        { return hwr.Req }
func (hwr HttpWebRequest) Target() string            { return hwr.Req.RequestURI }
func (hwr HttpWebRequest) Method() string            { return hwr.Req.Method }
func (hwr HttpWebRequest) Header(name string) string { return hwr.Req.Header.Get(name) }
func (hwr HttpWebRequest) Body() ([]byte, error)     { return ioutil.ReadAll(hwr.Req.Body) }
func (hwr HttpWebRequest) Param(name string) string {
	if v := hwr.Req.URL.Query().Get(name); len(v) > 0 {
		return v
	} else if v := bone.GetValue(hwr.Req, name); len(v) > 0 {
		return v
	} else {
		return ""
	}
}

/* response infos */
type ResponseOut struct {
	statusCode   int
	headers      map[string]string
	responseBody []byte
}

func NewResponseOut() *ResponseOut {
	return &ResponseOut{
		statusCode:   http.StatusOK,
		headers:      map[string]string{},
		responseBody: nil,
	}
}

func (ri *ResponseOut) GetStatus() int {
	return ri.statusCode
}

func (ri *ResponseOut) GetHeader(name string) string {
	return ri.headers[name]
}

func (ri *ResponseOut) GetBody() []byte {
	return ri.responseBody
}

func (ri *ResponseOut) Status(statusCode int) *ResponseOut {
	ri.statusCode = statusCode
	return ri
}

func (ri *ResponseOut) JsonHeader() *ResponseOut {
	ri.headers["Content-Type"] = "application/json"
	return ri
}

func (ri *ResponseOut) LocationHeader(location string) *ResponseOut {
	ri.headers["Location"] = location
	return ri
}

func (ri *ResponseOut) ETagHeader(version string) *ResponseOut {
	ri.headers["ETag"] = version
	return ri
}

func (ri *ResponseOut) Header(k, v string) *ResponseOut {
	ri.headers[k] = v
	return ri
}

func (ri *ResponseOut) Body(content []byte) *ResponseOut {
	ri.responseBody = content
	return ri
}

/* End point handlers */

func AuthHandler(next EndpointHandler) EndpointHandler {
	return func(req HttpWebRequest, ctx context.Context, config *MapPropertySource) (info *ResponseOut) {
		if Configs.GetBool("disable-auth") {
			return next(req, ctx, config)
		}

		return next(req, ctx, config)
	}
}

func ErrorRecovery(next EndpointHandler) EndpointHandler {
	return func(req HttpWebRequest, ctx context.Context, config *MapPropertySource) (info *ResponseOut) {
		defer func() {
			if r := recover(); r != nil {
				info = NewResponseOut()

				switch r.(type) {
				case *InvalidPathError:
					info.Status(http.StatusBadRequest)
					info.Body([]byte(
						fmt.Sprintf(
							errorTemplate,
							http.StatusBadRequest,
							r.(error).Error()),
					))
				case *InvalidParamGenericError:
					info.Status(http.StatusBadRequest)
					info.Body([]byte(
						fmt.Sprintf(
							errorTemplate,
							http.StatusBadRequest,
							r.(error).Error()),
					))
				case *InvalidParamError:
					info.Status(http.StatusBadRequest)
					info.Body([]byte(
						fmt.Sprintf(
							errorTemplate,
							http.StatusBadRequest,
							r.(error).Error()),
					))
				case *InvalidTypeError:
					info.Status(http.StatusBadRequest)
					info.Body([]byte(
						fmt.Sprintf(
							errorTemplate,
							http.StatusBadRequest,
							r.(error).Error()),
					))
				case *NoAttributeError:
					info.Status(http.StatusBadRequest)
					info.Body([]byte(
						fmt.Sprintf(
							errorTemplate,
							http.StatusBadRequest,
							r.(error).Error()),
					))

				case *MissingRequiredPropertyError:
					info.Status(http.StatusBadRequest)
					info.Body([]byte(
						fmt.Sprintf(
							errorTemplate,
							http.StatusBadRequest,
							r.(error).Error()),
					))

				case *DuplicateError:
					info.Status(http.StatusConflict)
					info.Body([]byte(
						fmt.Sprintf(
							errorTemplate,
							http.StatusConflict,
							"uniqueness",
							r.(error).Error()),
					))

				case *UnauthorisedError:
					info.Status(http.StatusUnauthorized)
					info.Body([]byte(
						fmt.Sprintf(
							errorTemplate,
							http.StatusUnauthorized,
							r.(error).Error()),
					))

				case *ForbiddenError:
					info.Status(http.StatusForbidden)
					info.Body([]byte(
						fmt.Sprintf(
							errorTemplate,
							http.StatusForbidden,
							r.(error).Error()),
					))

				case *UnverifiedDomain:
					info.Status(http.StatusNotModified)
					info.Body([]byte(
						fmt.Sprintf(
							errorTemplate,
							http.StatusNotModified,
							r.(error).Error()),
					))

				case *DatastoreError:
					info.Status(http.StatusServiceUnavailable)
					info.Body([]byte(
						fmt.Sprintf(
							errorTemplate,
							http.StatusServiceUnavailable,
							r.(error).Error()),
					))

				default:
					info.Status(http.StatusInternalServerError)
					info.Body([]byte(fmt.Sprintf(
						errorTemplate,
						http.StatusInternalServerError,
						r.(error).Error()),
					))
				}
			}

		}()
		return next(req, ctx, config)
	}
}

type EndpointHandler func(r HttpWebRequest, ctx context.Context, config *MapPropertySource) *ResponseOut

func Endpoint(next EndpointHandler, config *MapPropertySource) http.HandlerFunc {
	return http.HandlerFunc(func(rw http.ResponseWriter, req *http.Request) {
		ctx := context.Background()
		resp := next(HttpWebRequest{req}, ctx, config)
		for k, v := range resp.headers {
			rw.Header().Set(k, v)
		}
		rw.WriteHeader(resp.statusCode)
		_, _ = rw.Write(resp.responseBody)
	})
}

func InjectRequestScope(next EndpointHandler) EndpointHandler {
	return func(req HttpWebRequest, ctx context.Context, config *MapPropertySource) (info *ResponseOut) {
		if req.Header("X-Correlation-Id") != "" {
			ctx = context.WithValue(ctx, RequestId{}, req.Header("X-Correlation-Id"))
			ctx = context.WithValue(ctx, RequestTimestamp{}, time.Now().Unix())
		} else {
			uid := uuid.NewV4()
			ctx = context.WithValue(ctx, RequestId{}, uid.String())
			ctx = context.WithValue(ctx, RequestTimestamp{}, time.Now().Unix())
		}

		Info(ctx, req.Method()+" Requested ", req.Target(), " Headers ", req.Raw().Header)

		t := time.Now()
		resp := next(req, ctx,config)
		now := time.Now()
		Info(ctx,
			"Completed Req ", req.Target(), " time taken ( ", now.Sub(t), "s )",
			" resp: ", string(resp.responseBody))
		return resp
	}
}

func ErrorCheck(err error) {
	if err != nil {
		panic(err)
	}
}

func ErrorCheckNilThrowInvalidParam(obj interface{}) {
	if reflect.ValueOf(obj).IsNil() {
		panic(Error.InvalidParamGeneric())
	}
}

func ErrorCheckNil(obj interface{}) {
	if reflect.ValueOf(obj).IsNil() {
		panic(Error.Text("Internal Error Param issue "))
	}
}

func ExpectedResponseCode(ctx context.Context, resp *http.Response, expectedResponse int, text string) {
	if resp.StatusCode != expectedResponse {
		err := Error.Text(text)
		log.Printf("Unable to perform the request %d", resp.StatusCode)
		panic(err)
	}
}

func ExpectedResponseCodeParamValidation(ctx context.Context, resp *http.Response) {
	if resp.StatusCode == http.StatusNotFound {
		err := Error.InvalidParamGeneric()
		log.Printf("Unable to perform the request %d", resp.StatusCode)
		panic(err)
	}
}

func ErrorCheckForFalseAndThrowError(predicate bool, err error) {
	if !predicate {
		panic(err)
	}
}

func ErrorCheckForFalse(predicate bool, text string) {
	if !predicate {
		panic(Error.Text(text))
	}
}

func ErrorCheckForTrue(predicate bool, text string) {
	if !predicate {
		panic(Error.Text(text))
	}
}

func ErrorCheckForIntNotEqual(predicate int, notExpectedValue int, text string) {
	if predicate == notExpectedValue {
		panic(Error.InvalidParam(text, "Number", "Either count is not given or not proper value provided"))
	}
}

func ErrorCheckExpectedErrorCode(resp *http.Response, expected int, text string) {
	if resp.StatusCode != expected {
		panic(Error.Text(text))
	}
}

func ErrorCheckForEmptyStringAndReturn403(s string) {
	if s == "" {
		panic(Error.ForbiddenRequest())
	}
}