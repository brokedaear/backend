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

## Style Guide

Below is the style guide and choices for writing code. The focal point of these choices are readability.

### Errors

Prefer error verbosity over shorthands. Go has `err != nil` shorthands that allow the programmer to combine two lines of code into a single line. You must not follow this method. Instead, keep error declaration and nil checking in two lines. Here is an example:

YES:
```go
err := a.DoSomething()
if err != nil {
    return err
}
```

NO:
```go
if err := a.DoSomething(); err != nil {
return err
}
```


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
