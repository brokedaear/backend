// SPDX-FileCopyrightText: 2025 BROKE DA EAR LLC <https://brokedaear.com>
//
// SPDX-License-Identifier: Apache-2.0

package server

import "io"

type GRPCServer interface {
	io.Closer
}

type grpcServer struct{}
