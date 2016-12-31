
file(GLOB BASIC_SRC *.cpp)

add_executable(devicemanager-cpp ${BASIC_SRC})
target_link_libraries(devicemanager-cpp devicemanager)
