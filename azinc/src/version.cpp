#include <azin/version.hpp>
#include <azin/colors.hpp>
#include <iostream>

/// version 0.0.1 (wow so groundbreaking)
int versionCommand(int argc, char* argv[]) {
    (void)argc; (void)argv;
    std::cout << azin::ux::color::cyan << "Azin 0.0.1" << azin::ux::color::reset << "\n";
    return 0;
}
