cmake_minimum_required(VERSION 3.22)
project(Azin VERSION 0.0.1 LANGUAGES CXX)

include(FetchContent)

FetchContent_Declare(
        CLI11
        GIT_REPOSITORY https://github.com/CLIUtils/CLI11.git
        GIT_TAG v2.5.0
)

FetchContent_MakeAvailable(CLI11)

list(APPEND CMAKE_MODULE_PATH "${CMAKE_CURRENT_SOURCE_DIR}/cmake")
include(Language)
include(ProjectOptions)
include(DetectCompiler)
include(CompilerOptions)
include(Analyzers)

add_subdirectory(libs)

add_subdirectory(apps/azc)
