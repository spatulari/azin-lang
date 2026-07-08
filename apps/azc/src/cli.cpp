#include "azin/support/fs/filesystem.hpp"

#include <azin/diagnostic.hpp>
#include <azin/diagnostic_engine.hpp>
#include <azin/lexer.hpp>
#include <azin/source.hpp>
#include <azin/token.hpp>

#include <CLI/CLI.hpp>
#include <azc/cli.hpp>
#include <cppcoro/io_service.hpp>
#include <cppcoro/sync_wait.hpp>
#include <cppcoro/task.hpp>
#include <cppcoro/when_all_ready.hpp>
#include <filesystem>
#include <fmt/base.h>
#include <fmt/color.h>
#include <span>
#include <string_view>
#include <utility>

namespace cli = azc::cli;
namespace frontend = azc::frontend;
namespace fs = azin::support::fs;

namespace {

template <typename... Args>
auto errorf(frontend::diagnostic_severity const severity, fmt::format_string<Args...> fmt,
            Args &&...args) -> void {
    fmt::color foreground = fmt::color::white;
    switch (severity) {
    case frontend::diagnostic_severity::error: {
        foreground = fmt::color::red;
    } break;
    case frontend::diagnostic_severity::warning: {
        foreground = fmt::color::yellow;
    } break;
    default:
        break;
    }
    fmt::print(stderr, fg(foreground), fmt, std::forward<Args>(args)...);
    fmt::print(stderr, "\n");
}

auto print_diagnostics(std::span<frontend::diagnostic const> diagnostics) -> void {
    for (auto const &[severity, message] : diagnostics) {
        errorf(severity, "{}", message);
    }
}

auto process_events(cppcoro::io_service &io) -> cppcoro::task<> {
    io.process_events();
    co_return;
}

template <typename T>
auto run_task(cppcoro::task<T> task, cppcoro::io_service &io) -> T {
    auto wrapped_task = [](cppcoro::task<T> t, cppcoro::io_service &io_ref) -> cppcoro::task<T> {
        try {
            auto res = co_await t;
            io_ref.stop();
            co_return res;
        }
        catch (...) {
            io_ref.stop();
            throw;
        }
    }(std::move(task), io);

    auto [result, _] =
        cppcoro::sync_wait(cppcoro::when_all_ready(std::move(wrapped_task), process_events(io)));

    return result.result();
}

auto lex_file(cppcoro::io_service &io, std::filesystem::path input, bool printTokens)
    -> cppcoro::task<int> {
    auto file = co_await fs::read_source_file_async(io, std::move(input));

    if (!file) {
        errorf(frontend::diagnostic_severity::error, "{}", file.error().message);
        co_return 1;
    }

    source::manager source{std::move(*file), input};

    frontend::diagnostic_engine diagnostics;
    frontend::lexer lexer{source.text(), source.file_name(), diagnostics};

    print_diagnostics(diagnostics.diagnostics());

    if (printTokens) {
        for (auto const &token : lexer.tokens()) {
            fmt::println("{} \"{}\" ({}:{})", frontend::token_kind_to_string(token.kind),
                         lexer.get_lexeme(token), token.line, token.column);

            if (token.kind == frontend::token_kind::eof) {
                break;
            }
        }
    }

    co_return diagnostics.has_errors() ? 1 : 0;
}

} // namespace

auto cli::run(int argc, char const *const *argv) -> int {
    CLI::App app{"Azin Compiler"};

    bool version = false;
    bool printTokens = false;
    std::filesystem::path input;

    app.add_flag("--version", version, "Display the compiler version");
    app.add_flag("--print-tokens", printTokens, "Print the generated tokens");
    app.add_option("input", input, "Source file to compile");

    CLI11_PARSE(app, argc, argv);

    if (version) {
        fmt::println("azc {}", AZIN_COMPILER_VERSION);
        return 0;
    }

    if (input.empty()) {
        errorf(frontend::diagnostic_severity::error,
               "no input file specified\nUsage: azc <source>");
        return 1;
    }

    cppcoro::io_service io;
    return run_task(lex_file(io, std::move(input), printTokens), io);
}
