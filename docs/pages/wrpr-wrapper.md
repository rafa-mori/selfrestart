# SelfRestart CLI Wrapper Documentation

## Overview

The `cmd/wrpr.go` file implements a wrapper for the SelfRestart CLI, providing a structured interface for managing application commands and configurations.

## Main Structure

### `SelfRestart` Type

```go
type SelfRestart struct {
    parentCmdName string
    printBanner   bool
}
```

**Fields:**

- `parentCmdName`: Parent command name for example concatenation
- `printBanner`: Flag that controls whether the banner should be displayed

## Interface Methods

### Configuration Methods

- **`Alias()`**: Returns empty string (no alias defined)
- **`ShortDescription()`**: Short description of SelfRestart
- **`LongDescription()`**: Long description of SelfRestart  
- **`Usage()`**: Command usage pattern
- **`Examples()`**: List of usage examples
- **`Active()`**: Always returns `true` (active module)
- **`Module()`**: Returns "article" as module name

### Execution Methods

- **`Execute()`**: Executes the main command
- **`Command()`**: Builds and configures the main Cobra command

### Utility Methods

- **`SetParentCmdName()`**: Sets the parent command name
- **`concatenateExamples()`**: Concatenates examples with parent command name

## Main Functionalities

### 1. Cobra Command Configuration

The `Command()` method configures:

- Root command with use, aliases, and version
- Adds subcommands through `cc.ServiceCmdList()`
- Adds version command
- Sets usage definitions for all commands

### 2. Argument Processing

The code processes command line arguments to:

- Check if specific commands are being executed
- Configure short descriptions for commands without them

### 3. Environment Configuration

The `RegX()` function configures the instance based on environment variables:

- `ARTICLE_PRINT_BANNER`: Controls banner display (default: "true")

## Initialization Function

### `RegX()`

```go
func RegX() *SelfRestart
```

**Responsibilities:**

- Reads `GOFORGE_PRINT_BANNER` environment variable
- Creates and returns new `SelfRestart` instance
- Sets default configuration for banner display

## Dependencies

- **Cobra**: CLI framework for Go
- **Internal modules**:
  - `cc`: CLI commands
  - `gl`: Logging system
  - `vs`: Version management

## Design Patterns

- **Wrapper Pattern**: Encapsulates Cobra functionality
- **Factory Pattern**: `RegX()` function for instance creation
- **Interface Segregation**: Specific methods for different CLI aspects

## Notes

- Module name is hardcoded as "article"
- Support for configuration via environment variables
- Integration with logging system for debugging
- Automatic configuration of commands and subcommands

## Usage Example

```go
// Create SelfRestart instance
selfrestart := RegX()

// Execute CLI
if err := selfrestart.Execute(); err != nil {
    log.Fatal(err)
}
```

## Environment Variables

| Variable | Description | Default |
|----------|-------------|---------|
| `ARTICLE_PRINT_BANNER` | Controls banner display | `"true"` |

## Command Structure

```text
article
├── [service commands from cc.ServiceCmdList()]
└── version (from vs.CliCommand())
```
