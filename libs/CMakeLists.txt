add_subdirectory(support)

add_library(azin_frontend INTERFACE)
add_library(Azin::Frontend ALIAS azin_frontend)
target_include_directories(azin_frontend INTERFACE
    $<BUILD_INTERFACE:${CMAKE_CURRENT_SOURCE_DIR}/..>
)
target_link_libraries(azin_frontend INTERFACE
    Azin::CompilerFlags
    Azin::Support::Filesystem
    Azin::Support::ANSI
)