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
 * @brief Result of opening a source file.
 *
 * On success, contains an input stream opened for reading. On failure,
 * contains a ::FileError describing the reason the file could not be
 * opened.
 */
using FileResult = std::expected<std::ifstream, FileError>;

/**
 * @brief Opens an Azin source file.
 *
 * Before opening the file, this function verifies that:
 * - the path refers to an existing regular file,
 * - the file has a supported source extension (`.az`),
 * - the file can be opened for reading.
 *
 * The returned stream is opened in binary mode and positioned at the
 * beginning of the file.
 *
 * @param source_path Path to the source file.
 *
 * @return An input stream on success; otherwise a ::FileError describing
 *         why the operation failed.
 */
[[nodiscard]]
auto open_source_file(std::filesystem::path const &source_path) -> FileResult;

} // namespace azin::support::fs
