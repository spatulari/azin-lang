option(AZIN_ENABLE_ANALYZERS "Enable static analysis tools" ON)
option(AZIN_ENABLE_CLANG_TIDY "Enable clang-tidy" ON)
option(AZIN_ENABLE_CPPCHECK "Enable cppcheck" OFF)
option(AZIN_ENABLE_IWYU "Enable include-what-you-use" OFF)

option(AZIN_ANALYZERS_STRICT "Treat analyzer warnings as errors" OFF)

find_program(AZIN_CLANG_TIDY_EXE NAMES clang-tidy)
find_program(AZIN_CPPCHECK_EXE NAMES cppcheck)
find_program(AZIN_IWYU_EXE NAMES include-what-you-use)

set(AZIN_HAS_CLANG_TIDY OFF CACHE INTERNAL "")
set(AZIN_HAS_CPPCHECK OFF CACHE INTERNAL "")
set(AZIN_HAS_IWYU OFF CACHE INTERNAL "")

if (AZIN_ENABLE_ANALYZERS AND AZIN_ENABLE_CLANG_TIDY AND AZIN_CLANG_TIDY_EXE)
    set(_driver_mode "")
    if (CMAKE_CXX_COMPILER_ID STREQUAL "GNU")
        set(_driver_mode "--extra-arg-before=--driver-mode=g++")
    elseif (CMAKE_CXX_COMPILER_ID MATCHES "Clang")
        set(_driver_mode "--extra-arg-before=--driver-mode=clang++")
    endif ()

    set(_clang_tidy_args
            ${AZIN_CLANG_TIDY_EXE}
            ${_driver_mode}
            --extra-arg=-std=c++${AZIN_CXX_STANDARD}
            --extra-arg=-Wno-unknown-warning-option
    )

    if (AZIN_ANALYZERS_STRICT)
        list(APPEND _clang_tidy_args --warnings-as-errors=*)
    endif ()

    set(CMAKE_CXX_CLANG_TIDY "${_clang_tidy_args}" CACHE STRING "" FORCE)
    set(AZIN_HAS_CLANG_TIDY ON CACHE INTERNAL "" FORCE)

    message(STATUS "Analyzer: clang-tidy enabled")
elseif (AZIN_ENABLE_CLANG_TIDY AND NOT AZIN_CLANG_TIDY_EXE)
    message(STATUS "Analyzer: clang-tidy requested but not found")
endif ()

if (AZIN_ENABLE_ANALYZERS AND AZIN_ENABLE_CPPCHECK AND AZIN_CPPCHECK_EXE)

    set(_cppcheck_args
            ${AZIN_CPPCHECK_EXE}
            --enable=warning,performance,portability
            --inline-suppr
            --quiet
            --std=c++${AZIN_CXX_STANDARD}
    )

    if (AZIN_ANALYZERS_STRICT)
        list(APPEND _cppcheck_args --error-exitcode=1)
    endif ()

    set(CMAKE_CXX_CPPCHECK "${_cppcheck_args}" CACHE STRING "" FORCE)
    set(AZIN_HAS_CPPCHECK ON CACHE INTERNAL "" FORCE)

    message(STATUS "Analyzer: cppcheck enabled")
elseif (AZIN_ENABLE_CPPCHECK AND NOT AZIN_CPPCHECK_EXE)
    message(STATUS "Analyzer: cppcheck requested but not found")
endif ()

if (AZIN_ENABLE_ANALYZERS AND AZIN_ENABLE_IWYU AND AZIN_IWYU_EXE)
    set(CMAKE_CXX_INCLUDE_WHAT_YOU_USE
            ${AZIN_IWYU_EXE}
    )
    set(AZIN_HAS_IWYU ON CACHE INTERNAL "" FORCE)
    message(STATUS "Analyzer: include-what-you-use enabled")
elseif (AZIN_ENABLE_IWYU AND NOT AZIN_IWYU_EXE)
    message(STATUS "Analyzer: iwyu requested but not found")
endif ()

function(azin_enable_analysis target)
    if (NOT AZIN_ENABLE_ANALYZERS)
        return()
    endif ()

    if (AZIN_HAS_CLANG_TIDY)
        set_target_properties(${target} PROPERTIES
                CXX_CLANG_TIDY "${CMAKE_CXX_CLANG_TIDY}"
        )
    endif ()

    if (AZIN_HAS_CPPCHECK)
        set_target_properties(${target} PROPERTIES
                CXX_CPPCHECK "${CMAKE_CXX_CPPCHECK}"
        )
    endif ()

    if (AZIN_HAS_IWYU)
        set_target_properties(${target} PROPERTIES
                CXX_INCLUDE_WHAT_YOU_USE "${CMAKE_CXX_INCLUDE_WHAT_YOU_USE}"
        )
    endif ()
endfunction()