
# At the moment we specify these here, since the protobuf code doesn't pass these checks.
# TODO: include these flags for all C++ code written in lib/ excluding protobuf code

#set(CMAKE_CXX_FLAGS "${CMAKE_CXX_FLAGS} -std=gnu++0x -Wall -Wextra -Wwrite-strings -Werror -pedantic -Wno-long-long")
#set(CMAKE_CXX_FLAGS "${CMAKE_CXX_FLAGS} -Wshadow -Wpointer-arith -Wcast-qual -Wmissing-braces -Wmissing-prototypes")
# From the Debian hardening guide
set(CMAKE_CXX_FLAGS "${CMAKE_CXX_FLAGS} -fstack-protector -Wformat-security")

include_directories(${PROJECT_SOURCE_DIR}/lib)

add_subdirectory(openzwaved)

