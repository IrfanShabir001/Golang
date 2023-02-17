package shared

import (
	"context"
	"fmt"
	uuid "github.com/satori/go.uuid"
	"time"
)

func InjectServiceScope(t int64, next RoutineHandler) {
	ctx := context.Background()

	uid := uuid.NewV4()
	ctx = context.WithValue(ctx, RequestId{}, uid.String())
	ctx = context.WithValue(ctx, RequestTimestamp{}, time.Unix(t, 0).Unix())

	ErrorHandler(t, ctx, next)
}

type RoutineHandler func(t int64, ctx context.Context)

func ErrorHandler(t int64, ctx context.Context, next RoutineHandler) {
	defer GenericErrorHandler(ctx)
	next(t, ctx)
}

func GenericErrorHandler(ctx context.Context) {
	if r := recover(); r != nil {
		switch r.(type) {
		case *InvalidPathError:
			Fatal(ctx, fmt.Sprintf(
				errorTemplate,
				"InvalidPathError",
				r.(error).Error()))

		case *InvalidParamGenericError:
			Fatal(ctx, fmt.Sprintf(
				errorTemplate,
				"InvalidParamGenericError",
				r.(error).Error()))

		case *InvalidParamError:
			Fatal(ctx, fmt.Sprintf(
				errorTemplate,
				"InvalidParamError",
				r.(error).Error()))

		case *InvalidTypeError:
			Fatal(ctx, fmt.Sprintf(
				errorTemplate,
				"InvalidTypeError",
				r.(error).Error()))

		case *NoAttributeError:
			Fatal(ctx, fmt.Sprintf(
				errorTemplate,
				"NoAttributeError",
				r.(error).Error()))

		case *MissingRequiredPropertyError:
			Fatal(ctx, fmt.Sprintf(
				errorTemplate,
				"MissingRequiredPropertyError",
				r.(error).Error()))

		case *DuplicateError:
			Fatal(ctx, fmt.Sprintf(
				errorTemplate,
				"DuplicateError",
				r.(error).Error()))

		case *UnauthorisedError:
			Fatal(ctx, fmt.Sprintf(
				errorTemplate,
				"UnauthorisedError",
				r.(error).Error()))

		case *ForbiddenError:
			Fatal(ctx, fmt.Sprintf(
				errorTemplate,
				"ForbiddenError",
				r.(error).Error()))

		case *UnverifiedDomain:
			Fatal(ctx, fmt.Sprintf(
				errorTemplate,
				"UnverifiedDomain",
				r.(error).Error()))

		case *PaymentInvalidError:
			Fatal(ctx, fmt.Sprintf(
				errorTemplate,
				"PaymentInvalidError",
				r.(error).Error()))

		default:
			Fatal(ctx, fmt.Sprintf(
				errorTemplate,
				"Internal Error Unkown",
				r.(error).Error()))
		}
	}

}
