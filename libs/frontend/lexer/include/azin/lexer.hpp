#pragma once

#include <azin/diagnostic_engine.hpp>
#include <azin/token.hpp>

#include <optional>
#include <string_view>

namespace azc::frontend {

class lexer {
public:
    explicit lexer(std::string_view source, std::string_view filename,
                   diagnostic_engine &diagnostics);

    /**
     * @brief Scans and returns the next token in the source stream.
     * Returns an EOF token when the end of the file is reached.
     */
    [[nodiscard]] auto next_token() -> token;

    /**
     * @brief Reconstructs the source text for a given token.
     */
    [[nodiscard]] auto get_lexeme(token const &t) const noexcept -> std::string_view;

private:
    std::string_view m_source;
    std::string_view m_filename;
    diagnostic_engine &m_diagnostics;

    std::uint32_t m_position{0};
    std::uint32_t m_line{1};
    std::uint32_t m_column{1};

    [[nodiscard]] auto is_at_end() const noexcept -> bool;
    [[nodiscard]] auto peek() const noexcept -> char;
    [[nodiscard]] auto peek_next() const noexcept -> char;

    auto advance() noexcept -> char;
    auto match(char expected) noexcept -> bool;
    auto skip_whitespace() noexcept -> void;
    auto recover_to(char delimiter) noexcept -> void;

    [[nodiscard]] auto scan_token() -> std::optional<token>;
    [[nodiscard]] auto identifier() -> token;
    [[nodiscard]] auto number() -> token;
    [[nodiscard]] auto character() -> std::optional<token>;
    [[nodiscard]] auto string() -> std::optional<token>;

    [[nodiscard]] auto make_token(token_kind kind, std::uint32_t start, std::uint32_t line,
                                  std::uint32_t column) const -> token;

    [[nodiscard]] static auto identifier_kind(std::string_view lexeme) noexcept -> token_kind;
};

} // namespace azc::frontend
