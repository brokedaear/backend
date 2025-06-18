<!--
SPDX-FileCopyrightText: 2025 BROKE DA EAR LLC <https://brokedaear.com>

SPDX-License-Identifier: Apache-2.0
-->

# Backend Monorepo

The backend at BdE is currently monolithic with considerations for a microservice architecture in mind.

## Tools and Dependencies

- Nix package manager

Nix is used to setup the development environment and build packages.

For your convenience, there are also some command line tools available in the nix environment, courtesy of da-flake:

- git
- fish
- ripgrep
- helix
- neovim
- jq
- yq

To format all files, use the command

```shell
nix fmt .
```

## Testing

This project uses its own assertion library, which is in the "assert" package (internal/common/tests/assert). Check the functions in that package to learn more.

There is also a struct base for table tests. You can find that in the "test" package (internal/common/tests/test).

For unit testing, black box tests are done in packages with the `_test` suffix. White box testing can be done by creating a file with the name of the original file and appending the `_internal_test` suffix. This file then uses the same package as the package being tested.

### Black box

```go
// pkg_test.go

package pkg_test

// Tests for public fields, methods, and functions...
```

### White box

```go
// pkg_internal_test.go

package pkg

// Tests for private fields, methods, and functions...
```

## Style Guide

Below is the style guide and choices for writing code. The focal point of these choices are readability.

### Shorthands

Prefer error verbosity over shorthands. Go has `err != nil` shorthands that allow the programmer to combine two lines of code into a single line. You must not follow this method. Instead, keep error declaration and nil checking in two lines. Here is an example:

YES:

```go
// Write it like this...
err := a.DoSomething()
if err != nil {
    return err
}
```

NO:

```go
// Do not do this...
if err := a.DoSomething(); err != nil {
    return err
}
```

Similarly, do the same for idioms like `_, ok := ...`.

### Comments

Spacing between comments and code are either are either zero or one. This project differentiates between zero or one spaces. One space means that the comment(s) describe a section of code. Zero spaces describes either a single function, variable, constant, or type. Zero space code blocks follow the recommended Go comment format, where the name of the function, variable, constant, or type comes first, then the description.

```go
// s describes a name.
var s string

// The code below checks if s is equal to anything significant.
// It returns if there is anything cool.

if s == "goat" {
    return
} else if s == "neo" {
    // This place is where things happen.

    fmt.Println("neoooo")
}
```

All comments must end in a period.

### Functions

No named returns. Named returns make code harder to read.

### Types

Completely abstain from using the `interface{}` type. This is archaic. As of Go 1.18, the type `any` is preferred for used instead of `interface{}`.
