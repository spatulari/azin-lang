include(CPM)

CPMAddPackage(
        NAME CLI11
        VERSION 2.5.0
        GITHUB_REPOSITORY CLIUtils/CLI11
)

CPMAddPackage(
        NAME fmt
        GIT_TAG 12.2.0
        GITHUB_REPOSITORY fmtlib/fmt
)

# Fix for MSVC
set_target_properties(fmt PROPERTIES INTERFACE_SYSTEM_INCLUDE_DIRECTORIES $<TARGET_PROPERTY:fmt,INTERFACE_INCLUDE_DIRECTORIES>)

CPMAddPackage(
        NAME cppcoro
        GITHUB_REPOSITORY andreasbuhr/cppcoro
        GIT_TAG 8642e98596a92be30a2b061d3ed306d959d3214e
)

if(TARGET cppcoro)
    if(CMAKE_CXX_COMPILER_ID MATCHES "GNU|Clang")
        target_compile_options(cppcoro INTERFACE -w)

    elseif(MSVC OR CMAKE_CXX_COMPILER_ID STREQUAL "ClangCL")
        target_compile_options(cppcoro INTERFACE /W0 /wd4265 /wd4625 /wd4626 /wd4996)
    endif()
endif()