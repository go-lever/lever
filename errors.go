package lever

import (
	"fmt"

	"github.com/hashicorp/go-multierror"
)

type MultiErrors struct {
	errs multierror.Error
}

func NewMultiErrors() *MultiErrors {
	var err MultiErrors

	err.errs.ErrorFormat = func(es []error) string {
		var errs string

		for _, e := range es {

			if len(errs) == 0 {
				errs = e.Error()
				continue
			}
			errs = fmt.Sprintf("%s;;%s", errs, e.Error())
		}
		return errs
	}

	return &err
}

func (merr *MultiErrors) Append(err error) {
	merr.errs.Errors = append(merr.errs.Errors, err)
}

func (merr *MultiErrors) Error() string {
	return merr.errs.Error()
}