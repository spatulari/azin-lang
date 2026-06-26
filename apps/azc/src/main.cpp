#include "ast/nodes/nodes.hpp"
#include "frontend/lexer/lexer.hpp"
#include "frontend/parser/parser.hpp"
#include "frontend/source/source.hpp"
#include "sema/checker/checker.hpp"

#include <azin/support/ansi/styled_view.hpp>
#include <azin/support/fs/filesystem.hpp>

#include <exception>
#include <filesystem>
#include <format>
#include <iostream>
#include <iterator>
#include <string>
#include <utility>

namespace ansi = azin::support::ansi;
namespace fs = azin::support::fs;

namespace {

auto print_error(std::string const &msg) -> void {
    std::cerr << ansi::red(msg) << '\n';
}

auto format_loc(azin::frontend::Span const &s) -> std::string {
    return std::format("{}:{}", s.begin.line, s.begin.col);
}

auto preload_qualified(azin::ast::Stmt const &stmt, azin::sema::Checker &checker) -> void {
    if (auto const *es = dynamic_cast<azin::ast::ExprStmt const *>(&stmt)) {
        if (auto const *qc = dynamic_cast<azin::ast::QualifiedCallExpr const *>(es->expr.get())) {
            checker.ensure_std_namespace(qc->path, qc->span);
        }
    }
}

auto run(int const argc, char const *const *argv) -> int {
    if (argc < 2) {
        std::cout << ansi::green(
            std::format("Usage: {} <source.az>\n", argc > 0 ? argv[0] : "azc")); // NOLINT(cppcoreguidelines-pro-bounds-pointer-arithmetic)
        return 1;
    }

    std::filesystem::path const src_path{argv[1]}; // NOLINT(cppcoreguidelines-pro-bounds-pointer-arithmetic)

    if (auto r = fs::check_file_exists(src_path); !r) { print_error(r.error().message); return 1; }
    if (auto r = fs::check_extension(src_path);   !r) { print_error(r.error().message); return 1; }

    auto file_r = fs::open_source_file(src_path);
    if (!file_r) { print_error(file_r.error().message); return 1; }

    std::string const src{std::istreambuf_iterator<char>(*file_r), std::istreambuf_iterator<char>()};
    std::string const src_name = src_path.string();
    azin::frontend::SourceFile const source{.name = src_name, .text = src};

    auto tokens = azin::frontend::Lexer{source}.tokenize();
    auto mod    = azin::frontend::Parser{std::move(tokens)}.parse_module();

    azin::sema::Checker checker;
    checker.declare_module(mod);
    checker.resolve_imports(mod);

    for (auto const &fn : mod.fns) {
        for (auto const &stmt : fn.body) {
            preload_qualified(*stmt, checker);
        }
    }

    for (auto const &fn : mod.fns) {
        for (auto const &stmt : fn.body) {
            if (auto const *vd = dynamic_cast<azin::ast::VarDecl const *>(stmt.get())) {
                if (auto err = azin::sema::check_ambiguous_var(*vd, checker.fn_table())) {
                    print_error(std::format("{}:{}: error: {}", src_path.string(), format_loc(err->span), err->what()));
                    return 1;
                }
            }
        }
    }

    checker.check_module(mod);

    return 0;
}

} // namespace

auto main(int const argc, char const *const *argv) noexcept(false) -> int {
    try {
        return run(argc, argv);
    } catch (azin::sema::SemaError const &e) {
        std::cerr << ansi::red(std::format("error at {}: {}\n", format_loc(e.span), e.what()));
        return 1;
    } catch (azin::frontend::ParseError const &e) {
        std::cerr << ansi::red(std::format("parse error at {}: {}\n", format_loc(e.span), e.what()));
        return 1;
    } catch (std::exception const &e) {
        std::cerr << ansi::red(std::format("error: {}\n", e.what()));
        return 1;
    } catch (...) {
        return 1;
    }
}
