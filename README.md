# Backend Monorepo

The backend at BdE is designed as a microservice architecture.

## `app`

App contains the main application code that interfaces with databases and the frontend.

## `logger`

Logger is a service that stores logs for system administration and data persistence.

## `monitor`

Monitor is a service that monitors environment and system health, along with the health of other microservices.

## `telemetry`

Telemetry is a service that communicates with plugin clients.

## Dependencies

- Nix package manager

Nix is used to setup the development environment and build packages.

## Style Guide

Below is the style guide and choices for writing code. The focal point of these choices are readability.

### Errors

Prefer error verbosity over shorthands. Go has `err != nil` shorthands that allow the programmer to combine two lines of code into a single line. You must not follow this method. Instead, keep error declaration and nil checking in two lines.

### Comments

Spacing between comments and code are either are either zero or one, as enforced by the Go linter. This project differentiates between zero or one spaces. One space means that the comment(s) describe a section of code. Zero spaces describes either a single function, variable, constant, or type. Zero space code blocks follow the recommended Go comment format, where the name of the function, variable, constant, or type comes first, then the description.
