
file(GLOB DEVICEMANAGER_SRC *.cpp)
include_directories($GRPCPP_INCLUDE_DIRS})

add_library(devicemanager ${DEVICEMANAGER_SRC})
target_link_libraries(devicemanager proto ${GRPCPP_LIBRARIES})
