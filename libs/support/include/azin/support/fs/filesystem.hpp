/**
 * @file filesystem.hpp
 * @brief Filesystem utilities for accessing Azin source files.
 */

#pragma once

#include <expected>
#include <filesystem>

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
 * @brief Reads an Azin source file entirely into memory.
 */
auto read_source_file(std::filesystem::path const&) -> std::expected<std::string, FileError>;

} // namespace azin::support::fs
