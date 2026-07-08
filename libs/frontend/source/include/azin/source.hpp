/**
 * @file source.hpp
 * @brief Declares the source file manager.
 */

#pragma once

#include <cstddef>
#include <filesystem>
#include <string>
#include <string_view>

namespace source {

/**
 * @brief Owns the contents and metadata of a loaded source file.
 *
 * The Manager is a read-only container for the source code buffer
 * and its associated file path. Traversal and parsing logic are
 * handled externally by the lexer.
 */
class manager {
public:
    /**
     * @brief Constructs a source manager.
     *
     * @param buffer The loaded source text.
     * @param path Path to the source file.
     */
    explicit manager(std::string buffer, std::filesystem::path path);

    /**
     * @brief Returns the entire source text.
     *
     * @return View of the loaded source buffer.
     */
    [[nodiscard]]
    auto text() const noexcept -> std::string_view;

    /**
     * @brief Returns the path to the source file.
     *
     * @return Constant reference to the source path.
     */
    [[nodiscard]]
    auto path() const noexcept -> std::filesystem::path const &;

    /**
     * @brief Returns the source file name.
     *
     * @return File name without directory components.
     */
    [[nodiscard]]
    auto file_name() const -> std::string;

private:
    /// Contents of the loaded source file.
    std::string m_buffer;

    /// Path to the source file.
    std::filesystem::path m_path;
};

} // namespace source
