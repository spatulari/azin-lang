#ifndef AZIN_HELPERS_HPP
    #define AZIN_HELPERS_HPP
    #include <fstream>


    namespace fs = std::filesystem;

    namespace azin::filesystem {
        int checkFileExists(const fs::path& source_path);
        int checkExtension(const fs::path& source_path);
        std::ifstream openSourceFile(const fs::path& source_path);
    }

#endif // AZIN_HELPERS_HPP