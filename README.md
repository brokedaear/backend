# backend

Website backend API and application that interfaces with frontend and database. This code is based on guides by Alex Edwards in _Let's Go_ and _Let's Go Further_, as well as `n30w/Darkspace`, with several modifications to fit business requirements and specifications.

## Dependencies

- Go
- Nixpacks
- golang-lint

## Style Guide

Below is the style guide and choices for writing code. The focal point of these choices are readability.

### Errors

Prefer error verbosity over shorthands. Go has `err != nil` shorthands that allow the programmer to combine two lines of code into a single line. You must not follow this method. Instead, keep error declaration and nil checking in two lines.

### Comments

Spacing between comments and code are either are either zero or one, as enforced by the Go linter. This project differentiates between zero or one spaces. One space means that the comment(s) describe a section of code. Zero spaces describes either a single function, variable, constant, or type. Zero space code blocks follow the recommended Go comment format, where the name of the function, variable, constant, or type comes first, then the description.
