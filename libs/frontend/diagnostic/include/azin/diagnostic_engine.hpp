/**
 * @file diagnostic_engine.hpp
 * @brief Declares the diagnostic engine.
 */

#pragma once

#include <azin/diagnostic.hpp>

#include <span>
#include <vector>

namespace azc::frontend {

/**
 * @brief Collects and manages compiler diagnostics.
 */
class diagnostic_engine {
public:
    /**
     * @brief Reports a diagnostic.
     * @param diag Diagnostic to report.
     */
    void report(diagnostic const &diag);

    /**
     * @brief Returns whether any reported diagnostic is an error.
     */
    [[nodiscard]]
    auto has_errors() const noexcept -> bool {
        return m_has_errors;
    }

    /**
     * @brief Returns all reported diagnostics.
     */
    [[nodiscard]]
    auto diagnostics() const noexcept -> std::span<diagnostic const> {
        return m_diagnostics;
    }

private:
    std::vector<diagnostic> m_diagnostics;

    /// O(1) cache to track if any errors have been reported.
    bool m_has_errors = false;
};

} // namespace azc::frontend
