package lever_test

import (
	"errors"
	"testing"

	"github.com/go-lever/lever"
	"github.com/stretchr/testify/assert"
)

func TestAppendErrorsAndPrint(t *testing.T) {
	tests := []struct {
		name   string
		errors []string
		want   string
	}{
		{
			name:   "One error",
			errors: []string{"one single error"},
			want:   "one single error",
		},
		{
			name:   "multiple errors",
			errors: []string{"first error", "second error", "third error"},
			want:   "first error;;second error;;third error",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := lever.NewMultiErrors()
			for _, e := range tt.errors {
				err.Append(errors.New(e))
			}
			assert.Equal(t, tt.want, err.Error())

		})
	}
}

func TestMultiErrors_ListAsString(t *testing.T) {
	tests := []struct {
		name string
		got  *lever.MultiErrors
		want []string
	}{
		{
			name: "One error as string",
			got:  lever.NewMultiErrors().Append(errors.New("an error")),
			want: []string{"an error"},
		},
		{
			name: "Multiple errors as string",
			got:  lever.NewMultiErrors().Append([]error{errors.New("an error"), errors.New("a second error")}...),
			want: []string{"an error", "a second error"},
		},
		{
			name: "no error",
			got:  lever.NewMultiErrors(),
			want: []string{},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			assert.Equal(t, tt.want, tt.got.ListAsString())
		})
	}
}

func TestMultiErrors_MultiErrorsJSON(t *testing.T) {
	tests := []struct {
		name string
		got  *lever.MultiErrors
		want lever.MultiErrorResponse
	}{
		{
			name: "One error as string",
			got:  lever.NewMultiErrors().Append(errors.New("an error")),
			want: lever.MultiErrorResponse{
				Err: []string{"an error"},
			},
		},
		{
			name: "Multiple errors as string",
			got:  lever.NewMultiErrors().Append([]error{errors.New("an error"), errors.New("a second error")}...),
			want: lever.MultiErrorResponse{
				Err: []string{"an error", "a second error"},
			},
		},
		{
			name: "no error",
			got:  lever.NewMultiErrors(),
			want: lever.MultiErrorResponse{
				Err: []string{},
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			assert.Equal(t, tt.want, lever.MultiErrorsJSON(tt.got))
		})
	}
}

func TestMultiErrors_MultiErrorsSlice(t *testing.T) {
	tests := []struct {
		name string
		got  error
		want []string
	}{
		{
			name: "One error as slice",
			got:  lever.NewMultiErrors().Append(errors.New("an error")),
			want: []string{"an error"},
		},
		{
			name: "Multiple errors as slice",
			got:  lever.NewMultiErrors().Append([]error{errors.New("an error"), errors.New("a second error")}...),
			want: []string{"an error", "a second error"},
		},
		{
			name: "no error",
			got:  lever.NewMultiErrors(),
			want: []string{},
		},
		{
			name: "One error that is not a multierror",
			got:  errors.New("basic error"),
			want: []string{"basic error"},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			errs := lever.MultiErrorsSlice(tt.got)
			assert.Equal(t, len(tt.want), len(errs))
			assert.Equal(t, tt.want, errs)
		})
	}
}
