#include <azin/support/fs/filesystem.hpp>

#include <cppcoro/io_service.hpp>
#include <cppcoro/read_only_file.hpp>
#include <cppcoro/task.hpp>
#include <expected>
#include <filesystem>
#include <string>

namespace azin::support::fs {

namespace {
constexpr std::string_view source_extension = ".az";

[[nodiscard]]
auto validate_source_file(std::filesystem::path const &source_path) -> Result {
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

// NOLINTNEXTLINE(cppcoreguidelines-avoid-reference-coroutine-parameters)
auto read_source_file_async(cppcoro::io_service &io_service,
                            std::filesystem::path const source_path) -> FileTask {
    if (auto result = validate_source_file(source_path); !result) {
        co_return std::unexpected(result.error());
    }

    auto const file = cppcoro::read_only_file::open(io_service, source_path);

    std::string buffer(file.size(), '\0');

    std::size_t const bytes_read = co_await file.read(0, buffer.data(), buffer.size());
    if (bytes_read != buffer.size()) {
        co_return std::unexpected(
            FileError{"Unexpected end of file while reading '" + source_path.string() + "'."});
    }

    co_return buffer;
}

} // namespace azin::support::fs
