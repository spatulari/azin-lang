#pragma once

namespace cli {
    [[nodiscard]]
    auto run(int const argc, char const *const *argv) -> int;
} // namespace cli