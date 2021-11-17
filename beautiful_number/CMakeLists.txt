cmake_minimum_required(VERSION 3.18)

project(beautiful_number VERSION 0.1 LANGUAGES CXX)

if(CMAKE_CXX_COMPILER_ID STREQUAL "Clang")
    add_compile_options(-fconstexpr-steps=1000000000)
endif()

add_executable(cpp11_constexpr cpp11_constexpr.cpp)
set_property(TARGET cpp11_constexpr PROPERTY CXX_STANDARD 11)

add_executable(cpp14_constexpr cpp14_constexpr.cpp)
set_property(TARGET cpp14_constexpr PROPERTY CXX_STANDARD 14)

add_executable(cpp11 cpp11.cpp)
set_property(TARGET cpp11 PROPERTY CXX_STANDARD 11)
