package acontext

import (
	"context"

	"github.com/danielkrainas/gobag/api/errcode"
)

func WithErrors(ctx context.Context, errors errcode.Errors) context.Context {
	return WithValue(ctx, "errors", errors)
}

func AppendError(ctx context.Context, err error) context.Context {
	errors := GetErrors(ctx)
	errors = append(errors, err)
	return WithErrors(ctx, errors)
}

func GetErrors(ctx context.Context) errcode.Errors {
	if errors, ok := ctx.Value("errors").(errcode.Errors); errors != nil && ok {
		return errors
	}

	return make(errcode.Errors, 0)
}
