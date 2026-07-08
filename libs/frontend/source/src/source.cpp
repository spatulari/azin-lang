#include <azin/source.hpp>

#include <cstddef>
#include <filesystem>
#include <fstream>
#include <ios>
#include <iterator>
#include <string_view>
#include <utility>

namespace source {

manager::manager(std::filesystem::path path)
    : m_path(std::move(path)) {
}

auto manager::load() -> bool {
    std::ifstream file(m_path, std::ios::binary);

    if (!file) {
        return false;
    }

    m_buffer.assign(std::istreambuf_iterator<char>{file}, std::istreambuf_iterator<char>{});

    m_position = 0;
    return true;
}

auto manager::text() const noexcept -> std::string_view {
    return m_buffer;
}

auto manager::current() const noexcept -> char {
    return peek();
}

auto manager::peek(std::size_t offset) const noexcept -> char {
    auto const pos = m_position + offset;

    if (pos >= m_buffer.size()) {
        return '\0';
    }

    return m_buffer[pos]; // NOLINT(cppcoreguidelines-pro-bounds-avoid-unchecked-container-access)
}

auto manager::advance() noexcept -> void {
    if (!eof()) {
        ++m_position;
    }
}

auto manager::eof() const noexcept -> bool {
    return m_position >= m_buffer.size();
}

auto manager::remaining() const -> std::string_view {
    return std::string_view{m_buffer}.substr(m_position);
}

auto manager::position() const noexcept -> std::size_t {
    return m_position;
}

auto manager::reset() noexcept -> void {
    m_position = 0;
}

auto manager::path() const noexcept -> std::filesystem::path const & {
    return m_path;
}

auto manager::file_name() const noexcept -> std::string {
    return m_path.filename().string();
}


} // namespace source
