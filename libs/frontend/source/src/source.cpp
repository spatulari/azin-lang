#include <azin/source.hpp>

#include <filesystem>
#include <fstream>
#include <iterator>
#include <string_view>
#include <utility>

namespace source {

manager::manager(std::string buffer, std::filesystem::path path)
    : m_buffer(std::move(buffer))
    , m_path(std::move(path)) {
}

auto manager::text() const noexcept -> std::string_view {
    return m_buffer;
}

auto manager::path() const noexcept -> std::filesystem::path const & {
    return m_path;
}

auto manager::file_name() const -> std::string {
    return m_path.filename().string();
}

} // namespace source
