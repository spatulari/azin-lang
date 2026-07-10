add_library(frontend_diagnostic STATIC
            src/diagnostic_engine.cpp
)

add_library(Azin::Frontend::Diagnostic ALIAS frontend_diagnostic)

target_include_directories(frontend_diagnostic
                           PUBLIC
                           $<BUILD_INTERFACE:${CMAKE_CURRENT_SOURCE_DIR}/include>
                           $<INSTALL_INTERFACE:include>
)

target_link_libraries(frontend_diagnostic
                      PRIVATE
                      Azin::CompilerFlags
                      fmt::fmt
)

azin_enable_analyzers(frontend_diagnostic)