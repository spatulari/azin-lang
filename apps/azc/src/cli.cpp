#include "azin/support/fs/filesystem.hpp"

#include <azin/diagnostic.hpp>
#include <azin/diagnostic_engine.hpp>
#include <azin/lexer.hpp>
#include <azin/source.hpp>
#include <azin/token.hpp>

#include <CLI/CLI.hpp>
#include <azc/cli.hpp>
#include <cstdio> // NOLINT
#include <filesystem>

// Disable the unreachable code warning for MSVC
#if defined(_MSC_VER) && !defined(__llvm__)
    #pragma warning(push)
    #pragma warning(disable : 4702)
#endif

#include <fmt/base.h>
#include <fmt/color.h>
#include <span>

#if defined(_MSC_VER) && !defined(__llvm__)
    #pragma warning(pop)
#endif

#include <string_view>

namespace cli = azc::cli;

namespace {

template <typename... Args>
void errorf(fmt::format_string<Args...> fmt_str, Args &&...args) {
    fmt::print(stderr, fg(fmt::color::red), fmt_str, std::forward<Args>(args)...);
    fmt::print(stderr, "\n");
}

void errprintln(std::string_view const msg) {
    fmt::print(stderr, fg(fmt::color::red), "error: {}\n", msg);
}

auto print_diagnostics(std::span<azc::frontend::diagnostic const> diagnostics) -> void {
    for (auto const &[_, message] : diagnostics) {
        errorf("{}", message);
    }
}

} // namespace

auto cli::run(int const argc, char const *const *argv) -> int {
    CLI::App app{"Azin Compiler"};

    bool version{false};
    bool print_tokens{false};
    std::filesystem::path input;

    app.add_flag("--version", version, "Display the compiler's version");
    app.add_option("input", input, "Source file to compile");
    app.add_flag("--print-tokens", print_tokens, "Print the tokens");

    CLI11_PARSE(app, argc, argv);

    if (version) {
        fmt::println("azc {}", AZIN_COMPILER_VERSION);
        return 0;
    }

    if (input.empty()) {
        errprintln("no input file specified\nUsage: azc <source>");
        return 1;
    }

    auto file_result = azin::support::fs::read_source_file(input);
    if (!file_result) {
        errorf("{}", file_result.error().message);
        return 1;
    }

    source::manager source{std::move(*file_result), input};
    fmt::println("Loaded {} bytes", source.text().size());

    frontend::diagnostic_engine diagnostics;
    frontend::lexer lexer{source.text(), source.file_name(), diagnostics};

    if (print_tokens) {
        frontend::token t;
        do {
            t = lexer.next_token();
            fmt::println("{} \"{}\" ({}:{})", azc::frontend::token_kind_to_string(t.kind),
                         lexer.get_lexeme(t), t.line, t.column);
        }
        while (t.kind != frontend::token_kind::eof);
    }

    print_diagnostics(diagnostics.diagnostics());

    return diagnostics.has_errors() ? 1 : 0;
}
