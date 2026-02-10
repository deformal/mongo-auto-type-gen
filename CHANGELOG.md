# Changelog

All notable changes to this project will be documented in this file.

The format is based on [Keep a Changelog](https://keepachangelog.com/en/1.0.0/),
and this project adheres to [Semantic Versioning](https://semver.org/spec/v2.0.0.html).

## [Unreleased]

### Fixed

- **Duplicate type names resolved**: Embedded object types that conflict with root collection names are now automatically prefixed with their parent collection name (e.g., `DeliveryTripsShipments` instead of duplicate `Shipments`)
- **Removed unnecessary parentheses in array types**: Simple array types now render as `number[]` instead of `(number)[]`. Parentheses are only added for union types like `(string | number)[]`
- **Eliminated redundant type unions**: When both semantic type aliases (like `ISODateString` or `ObjectId`) and base `string` type are detected, the semantic type is preserved to maintain clarity (e.g., `ISODateString` instead of `ISODateString | string`)

### Changed

- Type generation now tracks all root collection type names to prevent naming conflicts with nested types
- Array type rendering optimized to only use parentheses when necessary
- Union type deduplication now prefers semantic type aliases over base types for better developer experience

## [0.1.4] - 2025-02-10

### Fixed

- Fixed created file name and extension issues

## [0.1.3] - 2025-02-10

### Changed

- Updated README file

## [0.1.2] - 2025-02-10

### Added

- Initial release with core MongoDB to TypeScript type generation functionality

[Unreleased]: https://github.com/deformal/mongo-auto-type-gen/compare/v0.1.4...HEAD
[0.1.4]: https://github.com/deformal/mongo-auto-type-gen/releases/tag/v0.1.4
[0.1.3]: https://github.com/deformal/mongo-auto-type-gen/releases/tag/v0.1.3
[0.1.2]: https://github.com/deformal/mongo-auto-type-gen/releases/tag/v0.1.2
