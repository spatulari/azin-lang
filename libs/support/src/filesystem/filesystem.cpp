#include <azin/support/fs/filesystem.hpp>

#include <filesystem>
#include <fstream>
#include <ios>

namespace azin::support::fs {

namespace {

constexpr auto source_extension = ".az";

[[nodiscard]]
auto validate_source_file(std::filesystem::path const &source_path) -> Result {
    if (!std::filesystem::is_regular_file(source_path)) {
        return std::unexpected(FileError{"Source file does not exist: " + source_path.string()});
    }

    if (source_path.extension() != source_extension) {
        return std::unexpected(
            FileError{"Invalid source file extension: " + source_path.extension().string()});
    }

    return {};
}

} // namespace

auto open_source_file(std::filesystem::path const &source_path) -> FileResult {
    if (auto result = validate_source_file(source_path); !result) {
        return std::unexpected(result.error());
    }

    std::ifstream file(source_path, std::ios::binary);

    if (!file) {
        return std::unexpected(FileError{"Failed to open source file: " + source_path.string()});
    }

    return file;
}

} // namespace azin::support::fs
