cmake_minimum_required(VERSION 3.22)
project(Azin VERSION 0.0.1 LANGUAGES CXX)

if(EXISTS "${CMAKE_CURRENT_SOURCE_DIR}/CMakeCache.txt" OR
   EXISTS "${CMAKE_CURRENT_SOURCE_DIR}/CMakeFiles")
	message(FATAL_ERROR
	        "\n"
	        "============================================================\n"
	        "Azin detected a previous in-source CMake build attempt.\n"
	        "\n"
	        "Please remove the generated files from the source directory:\n"
	        "\n"
	        "  rm -rf CMakeCache.txt CMakeFiles\n"
	        "\n"
	        "Then configure using:\n"
	        "\n"
	        "  cmake -S . -B build\n"
	        "  cmake --build build\n"
	        "\n"
	        "============================================================\n"
	)
endif()

get_filename_component(
		AZIN_SOURCE_DIR_REAL
		"${CMAKE_CURRENT_SOURCE_DIR}"
		REALPATH
)

get_filename_component(
		AZIN_BINARY_DIR_REAL
		"${CMAKE_BINARY_DIR}"
		REALPATH
)

if(AZIN_SOURCE_DIR_REAL STREQUAL AZIN_BINARY_DIR_REAL)
	message(FATAL_ERROR
	        "\n"
	        "============================================================\n"
	        "Azin does not support in-source builds.\n"
	        "\n"
	        "Remove generated files from the source directory and run:\n"
	        "\n"
	        "  cmake -S . -B build\n"
	        "  cmake --build build\n"
	        "\n"
	        "============================================================\n"
	)
endif()

list(APPEND CMAKE_MODULE_PATH "${CMAKE_CURRENT_SOURCE_DIR}/cmake")
include(Language)
include(ProjectOptions)
include(DetectCompiler)
include(CompilerOptions)
include(Analyzers)
include(Doxygen)

include(Dependencies)

add_subdirectory(libs)

add_subdirectory(apps/azc)
