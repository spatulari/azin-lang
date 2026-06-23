cmake_minimum_required(VERSION 3.22)
project(Azin LANGUAGES CXX)

set(CMAKE_CXX_STANDARD 23)
set(CMAKE_CXX_STANDARD_REQUIRED ON)

add_subdirectory(cli)
add_subdirectory(azinc)

set(CMAKE_EXPORT_COMPILE_COMMANDS ON)

