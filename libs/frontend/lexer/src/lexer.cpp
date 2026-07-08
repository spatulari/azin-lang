#include <azin/diagnostic.hpp>
#include <azin/diagnostic_engine.hpp>
#include <azin/lexer.hpp>
#include <azin/token.hpp>

#include <cctype>
#include <fmt/format.h>
#include <string_view>

namespace azc::frontend {

lexer::lexer(std::string_view const source, std::string_view const filename,
             diagnostic_engine &diagnostics)
    : m_source(source)
    , m_filename(filename)
    , m_diagnostics(diagnostics) {
}

auto lexer::get_lexeme(token const &t) const noexcept -> std::string_view {
    return m_source.substr(t.offset, t.length);
}

auto lexer::next_token() -> token {
    while (!is_at_end()) {
        skip_whitespace();

        if (is_at_end()) {
            break;
        }

        if (auto const tok = scan_token()) {
            return *tok;
        }
        // If scan_token returned nullopt, it means we hit an error.
        // We already reported the diagnostic, so we just loop to the next token.
    }

    return make_token(token_kind::eof, m_position, m_line, m_column);
}

auto lexer::is_at_end() const noexcept -> bool {
    return m_position >= m_source.size();
}

auto lexer::peek() const noexcept -> char {
    return is_at_end() ? '\0' : m_source[m_position];
}

auto lexer::peek_next() const noexcept -> char {
    return (m_position + 1 >= m_source.size()) ? '\0' : m_source[m_position + 1];
}

auto lexer::advance() noexcept -> char {
    char const c = m_source[m_position++];
    if (c == '\n') {
        ++m_line;
        m_column = 1;
    }
    else {
        ++m_column;
    }
    return c;
}

auto lexer::match(char const expected) noexcept -> bool {
    if (is_at_end() || peek() != expected) {
        return false;
    }
    advance();
    return true;
}

auto lexer::skip_whitespace() noexcept -> void {
    while (!is_at_end()) {
        switch (peek()) {
        case ' ':
        case '\t':
        case '\r':
        case '\n':
            advance();
            break;
        default:
            return;
        }
    }
}

auto lexer::recover_to(char const delimiter) noexcept -> void {
    while (!is_at_end()) {
        if (peek() == delimiter) {
            advance();
            return;
        }
        if (peek() == '\n') {
            return;
        }
        advance();
    }
}

auto lexer::scan_token() -> std::optional<token> {
    auto const start = m_position;
    auto const line = m_line;
    auto const column = m_column;

    char c = peek();

    if (std::isalpha(static_cast<unsigned char>(c)) || c == '_') {
        return identifier();
    }
    if (std::isdigit(static_cast<unsigned char>(c))) {
        return number();
    }
    if (c == '\'') {
        return character();
    }
    if (c == '"') {
        return string();
    }

    c = advance();

    switch (c) {
    case '+':
        return make_token(token_kind::plus, start, line, column);
    case '-':
        return make_token(match('>') ? token_kind::arrow : token_kind::minus, start, line, column);
    case '*':
        return make_token(token_kind::star, start, line, column);
    case '/': {
        if (match('/')) {
            while (!is_at_end() && peek() != '\n') {
                advance();
            }
            return std::nullopt;
        }
        if (match('*')) {
            while (!is_at_end()) {
                if (peek() == '*' && peek_next() == '/') {
                    advance();
                    advance();
                    break;
                }
                advance();
            }
            if (is_at_end()) {
                m_diagnostics.report({diagnostic_severity::error,
                                      fmt::format("{}:{}:{}: Unterminated block comment.",
                                                  m_filename, line, column)});
            }
            return std::nullopt;
        }
    }
        return make_token(token_kind::slash, start, line, column);
    case '=':
        return make_token(match('=') ? token_kind::equal_equal : token_kind::equal, start, line,
                          column);
    case '!':
        return make_token(match('=') ? token_kind::bang_equal : token_kind::bang, start, line,
                          column);
    case '<':
        return make_token(match('=') ? token_kind::less_equal : token_kind::less, start, line,
                          column);
    case '>':
        return make_token(match('=') ? token_kind::greater_equal : token_kind::greater, start, line,
                          column);
    case '(':
        return make_token(token_kind::left_paren, start, line, column);
    case ')':
        return make_token(token_kind::right_paren, start, line, column);
    case '{':
        return make_token(token_kind::left_brace, start, line, column);
    case '}':
        return make_token(token_kind::right_brace, start, line, column);
    case '[':
        return make_token(token_kind::left_bracket, start, line, column);
    case ']':
        return make_token(token_kind::right_bracket, start, line, column);
    case ',':
        return make_token(token_kind::comma, start, line, column);
    case ';':
        return make_token(token_kind::semicolon, start, line, column);
    case ':':
        return make_token(token_kind::colon, start, line, column);
    case '.':
        return make_token(token_kind::dot, start, line, column);
    case '%':
        return make_token(token_kind::modulo, start, line, column);
    case '^':
        return make_token(token_kind::caret, start, line, column);
    case '~':
        return make_token(token_kind::tilde, start, line, column);
    case '|':
        return make_token(match('|') ? token_kind::logical_or : token_kind::pipe, start, line,
                          column);
    case '&':
        return make_token(match('&') ? token_kind::logical_and : token_kind::ampersand, start, line,
                          column);
    default:
        m_diagnostics.report(
            {diagnostic_severity::error,
             fmt::format("{}:{}:{}: Unexpected character '{}'.", m_filename, line, column, c)});
        return std::nullopt;
    }
}

auto lexer::identifier() -> token {
    auto const start = m_position;
    auto const line = m_line;
    auto const column = m_column;

    while (std::isalnum(static_cast<unsigned char>(peek())) || peek() == '_') {
        advance();
    }

    std::string_view const text = m_source.substr(start, m_position - start);
    return make_token(identifier_kind(text), start, line, column);
}

auto lexer::number() -> token {
    auto const start = m_position;
    auto const line = m_line;
    auto const column = m_column;

    while (std::isdigit(static_cast<unsigned char>(peek()))) {
        advance();
    }

    token_kind kind = token_kind::integer_literal;
    if (peek() == '.' && std::isdigit(static_cast<unsigned char>(peek_next()))) {
        kind = token_kind::float_literal;
        advance(); // consume '.'
        while (std::isdigit(static_cast<unsigned char>(peek()))) {
            advance();
        }
    }
    return make_token(kind, start, line, column);
}

auto lexer::character() -> std::optional<token> {
    auto const start = m_position;
    auto const line = m_line;
    auto const column = m_column;

    advance(); // opening '
    if (is_at_end()) {
        m_diagnostics.report(
            {diagnostic_severity::error,
             fmt::format("{}:{}:{}: Unterminated character literal.", m_filename, line, column)});
        return std::nullopt;
    }

    advance(); // character
    if (peek() != '\'') {
        m_diagnostics.report(
            {diagnostic_severity::error,
             fmt::format("{}:{}:{}: Character literal must contain exactly one character.",
                         m_filename, line, column)});
        recover_to('\'');
        return std::nullopt;
    }

    advance(); // closing '
    return make_token(token_kind::character_literal, start, line, column);
}

auto lexer::string() -> std::optional<token> {
    auto const start = m_position;
    auto const line = m_line;
    auto const column = m_column;

    advance(); // opening "
    while (!is_at_end() && peek() != '"' && peek() != '\n') {
        advance();
    }

    if (is_at_end()) {
        m_diagnostics.report(
            {diagnostic_severity::error,
             fmt::format("{}:{}:{}: Unterminated string literal.", m_filename, line, column)});
        recover_to('"');
        return std::nullopt;
    }

    advance(); // closing "
    return make_token(token_kind::string_literal, start, line, column);
}

auto lexer::make_token(token_kind const kind, std::uint32_t const start, std::uint32_t const line,
                       std::uint32_t const column) const -> token {
    return {.kind = kind,
            .length = m_position - start,
            .offset = start,
            .line = line,
            .column = column};
}

auto lexer::identifier_kind(std::string_view const lexeme) noexcept -> token_kind {
    if (lexeme == "fn") {
        return token_kind::kw_fn;
    }
    if (lexeme == "var") {
        return token_kind::kw_var;
    }
    if (lexeme == "return") {
        return token_kind::kw_return;
    }
    if (lexeme == "end") {
        return token_kind::kw_end;
    }
    if (lexeme == "char") {
        return token_kind::kw_char;
    }
    if (lexeme == "int") {
        return token_kind::kw_int;
    }
    if (lexeme == "float") {
        return token_kind::kw_float;
    }
    if (lexeme == "string") {
        return token_kind::kw_string;
    }
    return token_kind::identifier;
}

} // namespace azc::frontend
