add_library(frontend_parser STATIC
            src/parser.cpp
            include/azin/parser.hpp
)

add_library(Azin::Frontend::Parser ALIAS frontend_parser)

target_include_directories(frontend_parser
                           PUBLIC
                           $<BUILD_INTERFACE:${CMAKE_CURRENT_SOURCE_DIR}/include>
                           $<INSTALL_INTERFACE:include>
)

target_link_libraries(frontend_parser
                      PRIVATE
                      Azin::CompilerFlags
)

azin_enable_analyzers(frontend_parser)