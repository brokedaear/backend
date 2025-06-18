// SPDX-FileCopyrightText: 2025 BROKE DA EAR LLC <https://brokedaear.com>
//
// SPDX-License-Identifier: Apache-2.0

package errors_test

import (
	stdlibErrors "errors"
	"fmt"
	"testing"

	pkgErrors "github.com/pkg/errors"

	"backend.brokedaear.com/pkg/errors"
)

func BenchmarkNew(b *testing.B) {
	b.Run(
		"errors.New", func(b *testing.B) {
			for b.Loop() {
				_ = errors.New("benchmark error")
			}
		},
	)

	b.Run(
		"pkgErrors.New", func(b *testing.B) {
			for b.Loop() {
				_ = pkgErrors.New("benchmark error")
			}
		},
	)

	b.Run(
		"stdlibErrors.New", func(b *testing.B) {
			for b.Loop() {
				_ = stdlibErrors.New("benchmark error")
			}
		},
	)
}

func BenchmarkJoin(b *testing.B) {
	err1 := errors.New("error 1")
	err2 := errors.New("error 2")
	err3 := errors.New("error 3")

	b.Run(
		"errors.Join", func(b *testing.B) {
			for b.Loop() {
				_ = errors.Join(err1, err2, err3)
			}
		},
	)

	b.Run(
		"stdlibErrors.Join", func(b *testing.B) {
			for b.Loop() {
				_ = stdlibErrors.Join(err1, err2, err3)
			}
		},
	)
}

func BenchmarkWrap(b *testing.B) {
	base := errors.New("base error")

	b.Run(
		"errors.Wrap", func(b *testing.B) {
			for b.Loop() {
				_ = errors.Wrap(base, "wrapped")
			}
		},
	)

	b.Run(
		"pkgErrors.Wrap", func(b *testing.B) {
			for b.Loop() {
				_ = pkgErrors.Wrap(base, "wrapped")
			}
		},
	)
}

func BenchmarkWrapf(b *testing.B) {
	base := errors.New("base error")

	b.Run(
		"errors.Wrapf", func(b *testing.B) {
			for b.Loop() {
				_ = errors.Wrapf(base, "wrapped %d", b.N)
			}
		},
	)

	b.Run(
		"pkgErrors.Wrapf", func(b *testing.B) {
			for b.Loop() {
				_ = pkgErrors.Wrapf(base, "wrapped %d", b.N)
			}
		},
	)
}

func BenchmarkErrorf(b *testing.B) {
	b.Run(
		"errors.Errorf", func(b *testing.B) {
			for b.Loop() {
				_ = errors.Errorf("benchmark error %d", b.N)
			}
		},
	)

	b.Run(
		"pkgErrors.Errorf", func(b *testing.B) {
			for b.Loop() {
				_ = pkgErrors.Errorf("benchmark error %d", b.N)
			}
		},
	)

	b.Run(
		"fmt.Errorf", func(b *testing.B) {
			for b.Loop() {
				_ = fmt.Errorf("benchmark error %d", b.N)
			}
		},
	)
}

func BenchmarkWithStack(b *testing.B) {
	base := errors.New("base error")

	b.Run(
		"errors.WithStack", func(b *testing.B) {
			for b.Loop() {
				_ = errors.WithStack(base)
			}
		},
	)

	b.Run(
		"pkgErrors.WithStack", func(b *testing.B) {
			for b.Loop() {
				_ = pkgErrors.WithStack(base)
			}
		},
	)
}

func BenchmarkWithMessage(b *testing.B) {
	base := errors.New("base error")

	b.Run(
		"errors.WithMessage", func(b *testing.B) {
			for b.Loop() {
				_ = errors.WithMessage(base, "additional context")
			}
		},
	)

	b.Run(
		"pkgErrors.WithMessage", func(b *testing.B) {
			for b.Loop() {
				_ = pkgErrors.WithMessage(base, "additional context")
			}
		},
	)
}

func BenchmarkWithMessagef(b *testing.B) {
	base := errors.New("base error")

	b.Run(
		"errors.WithMessagef", func(b *testing.B) {
			for b.Loop() {
				_ = errors.WithMessagef(base, "context %d", b.N)
			}
		},
	)

	b.Run(
		"pkgErrors.WithMessagef", func(b *testing.B) {
			for b.Loop() {
				_ = pkgErrors.WithMessagef(base, "context %d", b.N)
			}
		},
	)
}

func BenchmarkIs(b *testing.B) {
	baseErr := errors.New("base error")
	wrappedErr := errors.Wrap(baseErr, "wrapped")
	differentErr := errors.New("different error")

	b.Run(
		"errors.Is_same", func(b *testing.B) {
			for b.Loop() {
				_ = errors.Is(wrappedErr, baseErr)
			}
		},
	)
	b.Run(
		"pkgErrors.Is_same", func(b *testing.B) {
			for b.Loop() {
				_ = pkgErrors.Is(wrappedErr, baseErr)
			}
		},
	)
	b.Run(
		"stdlibErrors.Is_same", func(b *testing.B) {
			for b.Loop() {
				_ = stdlibErrors.Is(wrappedErr, baseErr)
			}
		},
	)

	b.Run(
		"errors.Is_different", func(b *testing.B) {
			for b.Loop() {
				_ = errors.Is(wrappedErr, differentErr)
			}
		},
	)

	b.Run(
		"pkgErrors.Is_different", func(b *testing.B) {
			for b.Loop() {
				_ = pkgErrors.Is(wrappedErr, differentErr)
			}
		},
	)

	b.Run(
		"stdlibErrors.Is_different", func(b *testing.B) {
			for b.Loop() {
				_ = stdlibErrors.Is(wrappedErr, differentErr)
			}
		},
	)
}

func BenchmarkAs(b *testing.B) {
	customErr := customError{code: 404, msg: "not found"}
	wrappedErr := errors.Wrap(customErr, "wrapped")

	b.Run(
		"errors.As_success", func(b *testing.B) {
			for b.Loop() {
				var target customError
				_ = errors.As(wrappedErr, &target)
			}
		},
	)

	b.Run(
		"pkgErrors.As_success", func(b *testing.B) {
			for b.Loop() {
				var target customError
				_ = pkgErrors.As(wrappedErr, &target)
			}
		},
	)

	b.Run(
		"stdlibErrors.As_success", func(b *testing.B) {
			for b.Loop() {
				var target customError
				_ = stdlibErrors.As(wrappedErr, &target)
			}
		},
	)

	b.Run(
		"errors.As_failure", func(b *testing.B) {
			for b.Loop() {
				var target networkError
				_ = errors.As(wrappedErr, &target)
			}
		},
	)

	b.Run(
		"pkgErrors.As_failure", func(b *testing.B) {
			for b.Loop() {
				var target networkError
				_ = pkgErrors.As(wrappedErr, &target)
			}
		},
	)

	b.Run(
		"stdlibErrors.As_failure", func(b *testing.B) {
			for b.Loop() {
				var target networkError
				_ = stdlibErrors.As(wrappedErr, &target)
			}
		},
	)
}

func BenchmarkUnwrap(b *testing.B) {
	// Test with pkg/errors wrapped error
	pkgBase := errors.New("base error")
	pkgWrapped := errors.Wrap(pkgBase, "wrapped")

	// Test with stdlib wrapped error
	stdBase := errors.New("base error")
	stdWrapped := fmt.Errorf("wrapped: %w", stdBase)

	b.Run(
		"errors.PkgUnwrap_pkgWrapped", func(b *testing.B) {
			for b.Loop() {
				_ = errors.PkgUnwrap(pkgWrapped)
			}
		},
	)

	b.Run(
		"pkgErrors.Unwrap_pkgWrapped", func(b *testing.B) {
			for b.Loop() {
				_ = pkgErrors.Unwrap(pkgWrapped)
			}
		},
	)

	b.Run(
		"errors.StdUnwrap_pkgWrapped", func(b *testing.B) {
			for b.Loop() {
				_ = errors.StdUnwrap(pkgWrapped)
			}
		},
	)

	b.Run(
		"stdlibErrors.Unwrap_pkgWrapped", func(b *testing.B) {
			for b.Loop() {
				_ = stdlibErrors.Unwrap(pkgWrapped)
			}
		},
	)

	b.Run(
		"errors.PkgUnwrap_stdWrapped", func(b *testing.B) {
			for b.Loop() {
				_ = errors.PkgUnwrap(stdWrapped)
			}
		},
	)

	b.Run(
		"pkgErrors.Unwrap_stdWrapped", func(b *testing.B) {
			for b.Loop() {
				_ = pkgErrors.Unwrap(stdWrapped)
			}
		},
	)

	b.Run(
		"errors.StdUnwrap_stdWrapped", func(b *testing.B) {
			for b.Loop() {
				_ = errors.StdUnwrap(stdWrapped)
			}
		},
	)

	b.Run(
		"stdlibErrors.Unwrap_stdWrapped", func(b *testing.B) {
			for b.Loop() {
				_ = stdlibErrors.Unwrap(stdWrapped)
			}
		},
	)
}

func BenchmarkCause(b *testing.B) {
	baseErr := errors.New("root cause")
	wrappedOnce := errors.Wrap(baseErr, "first wrap")
	wrappedTwice := errors.Wrap(wrappedOnce, "second wrap")

	b.Run(
		"errors.Cause", func(b *testing.B) {
			for b.Loop() {
				_ = errors.Cause(wrappedTwice)
			}
		},
	)

	b.Run(
		"pkgErrors.Cause", func(b *testing.B) {
			for b.Loop() {
				_ = pkgErrors.Cause(wrappedTwice)
			}
		},
	)
}

func BenchmarkComplexOperations(b *testing.B) {
	b.Run(
		"errors.New_then_Wrap_then_WithStack", func(b *testing.B) {
			for b.Loop() {
				base := errors.New("base error")
				wrapped := errors.Wrap(base, "wrapped")
				_ = errors.WithStack(wrapped)
			}
		},
	)

	b.Run(
		"direct_calls_New_then_Wrap_then_WithStack", func(b *testing.B) {
			for b.Loop() {
				base := stdlibErrors.New("base error")
				wrapped := pkgErrors.Wrap(base, "wrapped")
				_ = pkgErrors.WithStack(wrapped)
			}
		},
	)

	b.Run(
		"errors.error_chain_traversal", func(b *testing.B) {
			base := errors.New("root")
			wrapped1 := errors.Wrap(base, "layer1")
			wrapped2 := errors.Wrap(wrapped1, "layer2")
			wrapped3 := errors.Wrap(wrapped2, "layer3")

			b.ResetTimer()
			for b.Loop() {
				_ = errors.Is(wrapped3, base)
				_ = errors.Cause(wrapped3)
			}
		},
	)

	b.Run(
		"direct_calls_error_chain_traversal", func(b *testing.B) {
			base := stdlibErrors.New("root")
			wrapped1 := pkgErrors.Wrap(base, "layer1")
			wrapped2 := pkgErrors.Wrap(wrapped1, "layer2")
			wrapped3 := pkgErrors.Wrap(wrapped2, "layer3")

			b.ResetTimer()
			for b.Loop() {
				_ = stdlibErrors.Is(wrapped3, base)
				_ = pkgErrors.Cause(wrapped3)
			}
		},
	)
}

func BenchmarkMemoryAllocation(b *testing.B) {
	b.Run(
		"errors.multiple_wraps", func(b *testing.B) {
			base := errors.New("base")
			b.ResetTimer()
			for b.Loop() {
				err := base
				for j := 0; j < 5; j++ {
					err = errors.Wrap(err, fmt.Sprintf("layer %d", j))
				}
			}
		},
	)

	b.Run(
		"direct_calls_multiple_wraps", func(b *testing.B) {
			base := stdlibErrors.New("base")
			b.ResetTimer()
			for b.Loop() {
				err := base
				for j := 0; j < 5; j++ {
					err = pkgErrors.Wrap(err, fmt.Sprintf("layer %d", j))
				}
			}
		},
	)
}
