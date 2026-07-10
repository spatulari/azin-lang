add_library(frontend_token INTERFACE)

add_library(Azin::Frontend::Token ALIAS frontend_token)

target_include_directories(frontend_token
                           INTERFACE
                           $<BUILD_INTERFACE:${CMAKE_CURRENT_SOURCE_DIR}/include>
                           $<INSTALL_INTERFACE:include>
)

target_link_libraries(frontend_token
                      INTERFACE
                      Azin::CompilerFlags
)