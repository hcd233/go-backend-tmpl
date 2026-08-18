package ierr

import (
	"errors"
	"testing"

	"github.com/hcd233/aris-api-tmpl/internal/common/model"
)

func TestInternalError_Error(t *testing.T) {
	t.Parallel()

	tests := []struct {
		name     string
		err      error
		expected string
	}{
		{
			name:     "sentinel only",
			err:      ErrInternal,
			expected: "internal_error",
		},
		{
			name:     "with message",
			err:      New(ErrDBQuery, "select users"),
			expected: "db_query: select users",
		},
		{
			name:     "with cause",
			err:      Wrap(ErrDBQuery, errors.New("connection refused"), "select users"),
			expected: "db_query: select users: connection refused",
		},
		{
			name:     "with formatted message",
			err:      Wrapf(ErrJWTDecode, errors.New("token expired"), "user %d", 42),
			expected: "jwt_decode: user 42: token expired",
		},
		{
			name:     "newf without cause",
			err:      Newf(ErrValidation, "field %s is required", "name"),
			expected: "validation: field name is required",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()
			if got := tt.err.Error(); got != tt.expected {
				t.Errorf("Error() = %q, want %q", got, tt.expected)
			}
		})
	}
}

func TestInternalError_Is(t *testing.T) {
	t.Parallel()

	cause := errors.New("underlying error")
	wrapped := Wrap(ErrDBQuery, cause, "query failed")

	t.Run("same sentinel matches", func(t *testing.T) {
		t.Parallel()
		if !errors.Is(wrapped, ErrDBQuery) {
			t.Error("expected wrapped error to match ErrDBQuery sentinel")
		}
	})

	t.Run("different sentinel does not match", func(t *testing.T) {
		t.Parallel()
		if errors.Is(wrapped, ErrDBCreate) {
			t.Error("expected wrapped error NOT to match ErrDBCreate sentinel")
		}
	})

	t.Run("non-InternalError target does not match", func(t *testing.T) {
		t.Parallel()
		if errors.Is(wrapped, errors.New("some other error")) {
			t.Error("expected wrapped error NOT to match plain error")
		}
	})
}

func TestInternalError_Unwrap(t *testing.T) {
	t.Parallel()

	cause := errors.New("root cause")
	wrapped := Wrap(ErrDBQuery, cause, "query failed")

	t.Run("unwrap returns cause", func(t *testing.T) {
		t.Parallel()
		var ie *InternalError
		if !errors.As(wrapped, &ie) {
			t.Fatal("expected errors.As to succeed")
		}
		if ie.Unwrap() != cause {
			t.Errorf("Unwrap() = %v, want %v", ie.Unwrap(), cause)
		}
	})

	t.Run("no cause returns nil", func(t *testing.T) {
		t.Parallel()
		noCause := New(ErrDBQuery, "no cause")
		var ie *InternalError
		if !errors.As(noCause, &ie) {
			t.Fatal("expected errors.As to succeed")
		}
		if ie.Unwrap() != nil {
			t.Errorf("Unwrap() = %v, want nil", ie.Unwrap())
		}
	})
}

func TestInternalError_BizError(t *testing.T) {
	t.Parallel()

	tests := []struct {
		name         string
		sentinel     *InternalError
		expectedCode int
		expectedMsg  string
	}{
		{"internal", ErrInternal, 10000, "InternalError"},
		{"unauthorized", ErrUnauthorized, 10001, "Unauthorized"},
		{"no permission", ErrNoPermission, 10002, "NoPermission"},
		{"data not exists", ErrDataNotExists, 10003, "DataNotExists"},
		{"bad request", ErrBadRequest, 10006, "BadRequest"},
		{"jwt decode", ErrJWTDecode, 10001, "Unauthorized"},
		{"jwt encode", ErrJWTEncode, 10000, "InternalError"},
		{"db query", ErrDBQuery, 10000, "InternalError"},
		{"validation", ErrValidation, 10006, "BadRequest"},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()
			biz := tt.sentinel.BizError()
			if biz.Code != tt.expectedCode {
				t.Errorf("BizError().Code = %d, want %d", biz.Code, tt.expectedCode)
			}
			if biz.Message != tt.expectedMsg {
				t.Errorf("BizError().Message = %q, want %q", biz.Message, tt.expectedMsg)
			}
		})
	}
}

func TestToBizError(t *testing.T) {
	t.Parallel()

	fallback := model.NewError(99999, "Fallback")

	t.Run("extracts biz error from InternalError", func(t *testing.T) {
		t.Parallel()
		err := Wrap(ErrJWTDecode, errors.New("expired"), "token check")
		biz := ToBizError(err, fallback)
		if biz.Code != 10001 {
			t.Errorf("ToBizError().Code = %d, want 10001", biz.Code)
		}
	})

	t.Run("returns fallback for non-InternalError", func(t *testing.T) {
		t.Parallel()
		err := errors.New("plain error")
		biz := ToBizError(err, fallback)
		if biz != fallback {
			t.Errorf("ToBizError() = %v, want fallback %v", biz, fallback)
		}
	})
}
