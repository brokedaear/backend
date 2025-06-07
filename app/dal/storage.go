// SPDX-FileCopyrightText: 2025 BROKE DA EAR LLC <https://brokedaear.com>
//
// SPDX-License-Identifier: Apache-2.0

package dal

type S3 struct{}

func NewS3Storage() *S3 {
	return &S3{}
}
