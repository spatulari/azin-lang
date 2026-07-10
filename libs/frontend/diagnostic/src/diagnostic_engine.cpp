#include <azin/diagnostic_engine.hpp>

#include <span>

namespace azc::frontend {

void diagnostic_engine::report(diagnostic const &diag) {
    if (diag.severity == diagnostic_severity::error) {
        m_has_errors = true;
    }

    m_diagnostics.push_back(std::move(diag));
}

} // namespace azc::frontend
