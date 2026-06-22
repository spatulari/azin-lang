#include <azin/new.hpp>
#include <azin/colors.hpp>

#include <fstream>
#include <iostream>
#include <filesystem>

namespace fs = std::filesystem;

/// writes a string to a file and cries if it fails
static bool write_file(const fs::path& path, const std::string& content) {
    std::ofstream file(path);
    if (!file.is_open()) {
        return false;
    }
    file << content;
    return file.good();
}

/// creates the project structure for the new language
static int create_project(const std::string& name) {
    fs::path root = fs::current_path() / name;

    if (fs::exists(root)) {
        std::cout << azin::ux::color::red
                  << "error: directory '" << name << "' already exists"
                  << azin::ux::color::reset << "\n";
        return 1;
    }

    fs::create_directory(root);
    fs::create_directory(root / "src");
    fs::create_directory(root / "bin");

    std::string toml =
        "[project]\n"
        "name = \"" + name + "\"\n"
        "version = \"0.1.0\"\n"
        "azin-version = \"0.0.1\"\n";

    std::string main_az =
        "fn main() {\n"
        "    // your azin code here\n"
        "}\n";

    if (!write_file(root / "azin.toml", toml)) {
        std::cout << azin::ux::color::red
                  << "error: failed to create azin.toml"
                  << azin::ux::color::reset << "\n";
        return 1;
    }

    if (!write_file(root / "src" / "main.az", main_az)) {
        std::cout << azin::ux::color::red
                  << "error: failed to create src/main.az"
                  << azin::ux::color::reset << "\n";
        return 1;
    }

    std::cout << azin::ux::color::green
              << "Created new project '" << name << "'"
              << azin::ux::color::reset << "\n";

    return 0;
}


int newCommand(int argc, char* argv[]) {
    if (argc < 2) {
        std::cout << azin::ux::color::green
                  << "Usage: azin new <project-name>"
                  << azin::ux::color::reset << "\n";
        return 1;
    }

    return create_project(argv[1]);
}
