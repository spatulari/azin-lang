#pragma once

#include <azin/diagnostic_engine.hpp>
#include <azin/token.hpp>

#include <cppcoro/generator.hpp>
#include <fmt/format.h>
#include <optional>
#include <string_view>

namespace azc::frontend {

class lexer {
public:
    lexer(std::string_view source, std::string_view filename, diagnostic_engine &diagnostics);

    [[nodiscard]]
    auto tokens() -> cppcoro::generator<token>;

    [[nodiscard]]
    auto get_lexeme(token const &token) const noexcept -> std::string_view;

private:
    struct token_start {
        std::uint32_t offset;
        std::uint32_t line;
        std::uint32_t column;
    };

    //
    // Source
    //
    std::string_view m_source;
    std::string_view m_filename;
    diagnostic_engine &m_diagnostics;

    //
    // Cursor
    //
    std::uint32_t m_position{0};
    std::uint32_t m_line{1};
    std::uint32_t m_column{1};

    //
    // Cursor helpers
    //
    [[nodiscard]] auto is_at_end() const noexcept -> bool;
    [[nodiscard]] auto peek() const noexcept -> char;
    [[nodiscard]] auto peek_next() const noexcept -> char;

    auto advance() noexcept -> char;
    auto match(char expected) noexcept -> bool;

    [[nodiscard]]
    auto mark() const noexcept -> token_start;

    //
    // Scanning
    //
    auto skip_whitespace() noexcept -> void;
    auto skip_line_comment() noexcept -> void;
    auto skip_block_comment(token_start start) noexcept -> void;
    auto recover_to(char delimiter) noexcept -> void;


    auto scan_identifier(token_start start) -> token;
    auto scan_number(token_start start) -> token;
    auto scan_character(token_start start) -> std::optional<token>;
    auto scan_string(token_start start) -> std::optional<token>;
    auto scan_slash(token_start start) -> std::optional<token>;
    auto scan_operator(char c, token_start start) -> std::optional<token>;

    [[nodiscard]]
    auto scan_token() -> std::optional<token>;

    //
    // Token helpers
    //

    auto make_token(token_kind kind, token_start start) const noexcept -> token;

    [[nodiscard]]
    static auto identifier_kind(std::string_view lexeme) noexcept -> token_kind;

    //
    // Diagnostics
    //
    template <typename... Args>
    auto report_error(std::uint32_t line, std::uint32_t column, fmt::format_string<Args...> format,
                      Args &&...args) -> void {
        m_diagnostics.report({diagnostic_severity::error,
                              fmt::format("{}:{}:{}: {}", m_filename, line, column,
                                          fmt::format(format, std::forward<Args>(args)...))});
    }
};

} // namespace azc::frontend
