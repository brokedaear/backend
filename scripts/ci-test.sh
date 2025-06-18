#!/usr/bin/env bash

# SPDX-FileCopyrightText: 2025 BROKE DA EAR LLC <https://brokedaear.com>
#
# SPDX-License-Identifier: Apache-2.0

set -euo pipefail

printf "\n\n"
figlet -f chunky BrokeDa
figlet -f chunky goTESTS
echo "Running Go tests..."
gotestsum --format testdox ./...

printf "\n\n"

cat ./scripts/shaka.txt
