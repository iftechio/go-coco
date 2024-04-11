package custom

import (
	"testing"

	"github.com/pkg/errors"
	"github.com/stretchr/testify/assert"
	"google.golang.org/grpc/status"
)

func TestError(t *testing.T) {
	t.Run("is error", func(t *testing.T) {
		err := errors.WithStack(ErrNoValidUser)
		assert.True(t, errors.Is(err, ErrNoValidUser))
	})
	t.Run("as error", func(t *testing.T) {
		err := errors.WithStack(ErrNoValidUser)
		target := &Error{}
		assert.True(t, errors.As(err, &target))
	})
}

func TestError_clone(t *testing.T) {
	t.Run("test nil", func(t *testing.T) {
		assert.Nil(t, (*Error)(nil).clone())
	})

	t.Run("test changing field", func(t *testing.T) {
		err1 := &Error{}
		err2 := err1.clone()

		err1.Toast = "123"
		err2.Toast = "abc"
		assert.NotEqual(t, err1.Toast, err2.Toast)
	})
}

func TestToGRPCError(t *testing.T) {
	_, ok := status.FromError(&Error{})
	assert.True(t, ok)
}
