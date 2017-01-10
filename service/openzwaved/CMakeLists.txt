
pkg_search_module(OPENZWAVE REQUIRED libopenzwave)

file(GLOB OPENZWAVED_SRC *.cpp)
include_directories(${OPENZWAVE_INCLUDE_DIRS})

add_executable(openzwaved ${OPENZWAVED_SRC})
target_link_libraries(openzwaved devicemanager ${OPENZWAVE_LIBRARIES})

