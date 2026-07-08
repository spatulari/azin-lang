/**
 * @file filesystem.hpp
 * @brief Filesystem utilities for accessing Azin source files.
 */

#pragma once

#include <expected>
#include <filesystem>
#include <fstream>

namespace azin::support::fs {

/**
 * @brief Describes a filesystem operation failure.
 *
 * Filesystem utilities report recoverable failures through
 * `std::expected` rather than throwing exceptions.
 */
struct FileError {
    /// Human-readable description of the failure.
    std::string message;
};

/**
 * @brief Result of a filesystem operation.
 *
 * A successful result contains no value. On failure, the contained
 * ::FileError describes the reason the operation could not be completed.
 */
using Result = std::expected<void, FileError>;

/**
 * @brief Result of reading a source file.
 * Contains the full source string on success.
 */
using FileResult = std::expected<std::string, FileError>;

/**
 * @brief Reads an Azin source file entirely into memory.
 */
[[nodiscard]]
auto read_source_file(std::filesystem::path const &source_path) -> FileResult;

} // namespace azin::support::fs
