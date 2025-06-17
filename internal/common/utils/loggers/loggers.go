// SPDX-FileCopyrightText: 2025 BROKE DA EAR LLC <https://brokedaear.com>
//
// SPDX-License-Identifier: Apache-2.0

package loggers

// logger describes a custom logging implementation. For example, one might
// implement a library like logrus or Zap (like us) or use Go's slog library
// as their preferred logger. This interface abstracts that functionality.
//
// The methods on logger are kept intentionally minimal. This makes the caller
// think less when using logger (that's good for everyone involved) and gives
// the caller the option to format their own string, since there is no <level>f
// method on logger (such as "Infof" or "Errorf").
type logger interface {
	Info(msg string, args ...any)
	Debug(msg string, args ...any)
	Warn(msg string, args ...any)
	Error(msg string, args ...any)
	Sync() error
}
