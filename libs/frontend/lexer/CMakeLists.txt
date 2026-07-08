add_library(frontend_lexer STATIC
            src/lexer.cpp
            include/azin/lexer.hpp
)

add_library(Azin::Frontend::Lexer ALIAS frontend_lexer)

target_include_directories(frontend_lexer
                           PUBLIC
                           $<BUILD_INTERFACE:${CMAKE_CURRENT_SOURCE_DIR}/include>
                           $<INSTALL_INTERFACE:include>
)

target_link_libraries(frontend_lexer
                      PUBLIC
                      Azin::CompilerFlags
                      Azin::Frontend::Token
                      Azin::Frontend::Diagnostic
                      cppcoro
                      fmt::fmt
)

azin_enable_analyzers(frontend_lexer)