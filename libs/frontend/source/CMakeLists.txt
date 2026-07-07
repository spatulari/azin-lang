add_library(frontend_source STATIC
            src/source.cpp
            include/azin/source.hpp
)

add_library(Azin::Frontend::Source ALIAS frontend_source)

target_include_directories(frontend_source
                           PUBLIC
                           $<BUILD_INTERFACE:${CMAKE_CURRENT_SOURCE_DIR}/include>
                           $<INSTALL_INTERFACE:include>
)

target_link_libraries(frontend_source
                      PRIVATE
                      Azin::CompilerFlags
)

azin_enable_analyzers(frontend_source)