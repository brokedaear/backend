#!/usr/bin/env bash

# SPDX-FileCopyrightText: 2025 BROKE DA EAR LLC <https://brokedaear.com>
#
# SPDX-License-Identifier: Apache-2.0

set -euo pipefail

printf "\n\n"
figlet -f chunky Broke
figlet -f chunky da
figlet -f chunky BENCH

# Colors for output.
readonly RED='\033[0;31m'
readonly GREEN='\033[0;32m'
readonly YELLOW='\033[1;33m'
readonly BLUE='\033[0;34m'
readonly NC='\033[0m' # No Color

# Print colored output.
print_info() {
  echo -e "${BLUE}ℹ️  $1${NC}"
}

print_success() {
  echo -e "${GREEN}✅ $1${NC}"
}

print_warning() {
  echo -e "${YELLOW}⚠️  $1${NC}"
}

print_error() {
  echo -e "${RED}❌ $1${NC}"
}

if [[ ! -f "go.mod" ]]; then
  print_error "go.mod not found. Please run this script from the project root."
  exit 1
fi

print_info "Running Go benchmarks..."

go test -run=^$ -bench=. -benchmem -v ./... | tee benchmark_results.txt

print_success "Benchmarks completed!"
print_info "Results saved to benchmark_results.txt"

# Generate a summary of the benchmarks.
print_info "Benchmark Summary:"
echo "===================="
grep "^Benchmark" benchmark_results.txt | head -20
if [[ $(grep -c "^Benchmark" benchmark_results.txt) -gt 20 ]]; then
  print_info "... and $(($(grep -c "^Benchmark" benchmark_results.txt) - 20)) more benchmarks"
fi
echo ""

# Check for any failures.
if grep -q "FAIL" benchmark_results.txt; then
  print_warning "Some benchmarks may have failed. Check benchmark_results.txt for details."
  exit 1
else
  print_success "All benchmarks completed successfully!"
fi

printf "\n\n"

figlet -f cricket hoFastAh? | dotacat
