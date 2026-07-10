#include <azin/support/fs/filesystem.hpp>

#include <expected>
#include <filesystem>
#include <fstream>
#include <string>

namespace azin::support::fs {

namespace {
constexpr std::string_view source_extension = ".az";

[[nodiscard]]
auto validate_source_file(std::filesystem::path const &source_path)
    -> std::expected<void, FileError> {
    std::error_code ec;

    if (!std::filesystem::is_regular_file(source_path, ec)) {
        return std::unexpected(FileError{"Source file does not exist or is not a regular file: " +
                                         source_path.string()});
    }

    if (source_path.extension() != source_extension) {
        return std::unexpected(FileError{"Unsupported source file extension '" +
                                         source_path.extension().string() + "' (expected .az)."});
    }

    return {};
}
} // namespace

auto read_source_file(std::filesystem::path const &source_path)
    -> std::expected<std::string, FileError> {
    if (auto result = validate_source_file(source_path); !result) {
        return std::unexpected(result.error());
    }

    std::ifstream file(source_path, std::ios::binary);
    if (!file) {
        return std::unexpected(FileError{"Failed to open '" + source_path.string() + "'."});
    }

    file.seekg(0, std::ios::end);
    std::size_t size = static_cast<std::size_t>(file.tellg());
    file.seekg(0, std::ios::beg);

    std::string buffer(size, '\0');

    if (!file.read(buffer.data(), static_cast<std::streamsize>(size))) {
        return std::unexpected(FileError{"Failed to read '" + source_path.string() + "'."});
    }

    return buffer;
}

} // namespace azin::support::fs
