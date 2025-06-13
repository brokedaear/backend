#!/usr/bin/env bash

# SPDX-FileCopyrightText: 2025 BROKE DA EAR LLC <https://brokedaear.com>
#
# SPDX-License-Identifier: Apache-2.0

set -e

figlet -f small BrokeDaCI

echo "Linting Go Files..."
golangci-lint run

echo "Linting Licenses..."
reuse lint

echo "Running Go tests..."
gotestsum --format testdox ./internal/...

figlet -f cricket Cherreh
