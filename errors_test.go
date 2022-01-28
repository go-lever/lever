package lever

import (
	"errors"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestAppenErrorsAndPrint(t *testing.T) {
	tests := []struct {
		name string
		errors []string
		want string
	}{
		{
			name: "One error",
			errors: []string{"one single error"},
			want: "one single error",
		},
		{
			name: "multiple errors",
			errors: []string{"first error", "second error", "third error"},
			want: "first error;;second error;;third error",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := NewMultiErrors()
			for _, e := range tt.errors {
				err.Append(errors.New(e))
			}
			assert.Equal(t, tt.want, err.Error())

		})
	}
}
