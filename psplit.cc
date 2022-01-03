#include <iostream>
#include <random>
#include <string_view>
#include <cassert>
#include <vector>
#include <charconv>
#include <cstring>
#include <fstream>
#include <algorithm>

#include <fmt/format.h>

std::vector<char const*> targets;
std::vector<unsigned> weights;

void ensure(bool condition, char const* error_message)
{
	if (!condition) {
		std::cerr << "error: " << error_message << std::endl;
		std::exit(1);
	}
}

int main(int, char **argv)
{
	ensure(*argv++, "No arguments were provided");

	for (; *argv; ++argv) {
		if (targets.size() == weights.size()) {
			weights.emplace_back();
			std::from_chars(*argv, *argv + std::strlen(*argv), weights.back());
		} else {
			targets.emplace_back(*argv);
		}
	}

	ensure(targets.size(), "No arguments were provided");
	ensure(targets.size() == weights.size(), "Missing filename for given probability");

	std::discrete_distribution<unsigned> dist(weights.begin(), weights.end());

	std::vector<std::ofstream> files;
	for (auto target : targets) {
		files.emplace_back(target, std::ios::trunc | std::ios::in);
		ensure((bool)files.back(), "Error opening file");
	}

	std::mt19937 rnd{std::random_device{}()};
	for (std::string line; std::getline(std::cin, line); ) {
		auto &stream = files[dist(rnd)];
		stream << line << '\n';
	}
}
