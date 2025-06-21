# Changelog

All notable changes to this project will be documented in this file.

The format is based on [Keep a Changelog](https://keepachangelog.com/en/1.0.0/),
and this project adheres to [Semantic Versioning](https://semver.org/spec/v2.0.0.html).

## [Unreleased]

### Added
- Initial release of SelfRestart library
- Automatic process restart functionality
- Cross-platform support (Linux, macOS, Windows)
- Automatic Go installation detection and setup
- Modular architecture with internal packages
- Comprehensive logging system
- Process management utilities
- Platform detection capabilities
- Thread-safe operations
- Complete test suite
- Example application demonstrating usage
- Multi-language documentation (English/Portuguese)

### Features
- **Automatic Restart**: Restart applications preserving arguments and environment
- **Platform Detection**: Support for Linux, macOS and Windows
- **Go Installation**: Automatic Go installation if needed
- **Process Management**: Complete control over PIDs and signals
- **Integrated Logging**: Logging system with different levels
- **Modular Design**: Clean and well-organized architecture
- **Thread Safety**: Safe for use in concurrent applications

## [1.0.0] - 2025-06-21

### Added
- Initial public release
- Core restart functionality
- Documentation and examples
- MIT License

### Architecture
- `/selfrestart.go` - Main public API
- `/internal/install/` - Go installation management
- `/internal/platform/` - Platform detection
- `/internal/process/` - Process management
- `/internal/restart/` - Restart logic
- `/logger/` - Logging system
- `/example/` - Usage examples
