
pkg_search_module(OPENZWAVE REQUIRED libopenzwave)

file(GLOB OPENZWAVED_SRC *.cpp)
include_directories($GRPCPP_INCLUDE_DIRS} ${OPENZWAVE_INCLUDE_DIRS})

add_executable(openzwaved ${OPENZWAVED_SRC})
target_link_libraries(openzwaved proto ${GRPCPP_LIBRARIES} ${OPENZWAVE_LIBRARIES})

