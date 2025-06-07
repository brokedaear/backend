// SPDX-FileCopyrightText: 2025 BROKE DA EAR LLC <https://brokedaear.com>
//
// SPDX-License-Identifier: Apache-2.0

package infra

import "io"

func Teardown(closers ...io.Closer) error {
	for _, closer := range closers {
		err := closer.Close()
		if err != nil {
			return err
		}
	}

	return nil
}
