package lever

import (
	"fmt"

	"github.com/hashicorp/go-multierror"
)

// MultiErrors helps handling multiple errors.
// It aims to be used especially when validating some input that may contains multiple errors
// One may choose to return an error (bad request) when the first validation error occurs.
// But Multierrors helps stacking the validation errors and return them all at once.
type MultiErrors struct {
	errs multierror.Error
}

// MultiErrorResponse represents a typical api response in case of error.
// In this case, the response looks like :
//  {
// 	 "errors": [
// 	 	"first error",
// 	 	"second error"
// 	 ]
//  }
// and should be associated with the right status code (ex : 400, bad request)
type MultiErrorResponse struct {
	Err []string `json:"errors,omitempty"`
}

// NewMultiErrors creates a Multierrors
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

// Append append a new error to the error stack already contained in the Multierror.
// A typical use of the Append function may look likes this :
// errs := lever.NewMultiErrors().Append(errors.New("an error"))
func (merr *MultiErrors) Append(errs ...error) *MultiErrors {
	merr.errs.Errors = append(merr.errs.Errors, errs...)
	return merr
}

// Error return a concatenation of the error stack as a string.
func (merr *MultiErrors) Error() string {
	return merr.errs.Error()
}

// ListAsString returns the list of stacked error as a slice of strings.
func (merr *MultiErrors) ListAsString() []string {
	if len(merr.errs.Errors) == 0 {
		return []string{}
	}

	var err []string
	for _, e := range merr.errs.Errors {
		err = append(err, e.Error())
	}
	return err
}

// MultiErrorsJSON generates a MultiErrorResponse ready to be marshalled as JSON
// It should typically be used when an API returns a error code (ex : 400 bad request) to
// list all the errors in the server response.
//
// MultiErrorsJSON supports both error and MultiErrors type. Thus, it can be used right before
// returning the HTTP response without any consideration of either there are multiple errors to return
// or not. Example :
// err := Validate(myVar)
// ctx.StopWithJSON(400, lever.MultiErrorsJSON(err))
func MultiErrorsJSON(err error) MultiErrorResponse {
	switch e := err.(type) {
	case *MultiErrors:
		return MultiErrorResponse{Err: e.ListAsString()}
	case error:
		return MultiErrorResponse{Err: []string{e.Error()}}
	default:
		return MultiErrorResponse{Err: []string{e.Error()}}
	}
}

// MultiErrorsSlice returns a representation of the stacked errors as a slice of string.
// It is made to generate a slice of string from an error either it is a basic error or a MultiErrors.
// 
// MultiErrorsSlice is usefull when returning a list of error in the API response without relying
// on the MultiErrorResponse response format.
// Typically, it can be used to generate a list of validation error in a Problem response
// (rfc 7807 : https://datatracker.ietf.org/doc/html/rfc7807)
func MultiErrorsSlice(err error) []string {
	switch e := err.(type) {
	case *MultiErrors:
		return e.ListAsString()
	case error:
		return []string{e.Error()}
	default:
		return []string{e.Error()}
	}
}
