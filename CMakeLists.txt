
project(jvs)
cmake_minimum_required(VERSION 2.8.11)

set(CMAKE_CXX_STANDARD 11)
set(CMAKE_CXX_EXTENSIONS ON)

find_package(PkgConfig REQUIRED)
find_package(Protobuf REQUIRED)
find_package(GRPC REQUIRED)

include_directories(${PROTOBUF_INCLUDE_DIRS})

add_subdirectory(lib)
add_subdirectory(service)
add_subdirectory(test)

