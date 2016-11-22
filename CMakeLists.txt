
project(jvs)
cmake_minimum_required(VERSION 2.8.11)

find_package(PkgConfig REQUIRED)
find_package(Protobuf REQUIRED)

pkg_search_module(GRPCPP REQUIRED grpc++)

include_directories(${PROTOBUF_INCLUDE_DIRS})

add_subdirectory(lib)
add_subdirectory(service)

