#include <algorithm>
#include <charconv>
#include <fstream>
#include <map>
#include <vector>
#include <random>

#include <fmt/format.h>

using namespace fmt::literals;

std::map<unsigned, std::vector<std::string>> scores;

int main(int, char **argv)
{
	char const* input_path = *++argv;
	char const* max_diff_str = *++argv;

	if (!input_path || !max_diff_str) {
		fmt::print(stderr, "usage: {} <input-file> <max-difference>\n", basename(argv[-2]));
		return 1;
	}

	double max_diff;
	std::from_chars(max_diff_str, max_diff_str + strlen(max_diff_str), max_diff);
	std::ifstream input(input_path);

	for (std::string row; std::getline(input, row); ) {
		float score;
		auto r = std::from_chars(row.data(), row.data() + row.size(), score);
		if (r.ptr == row.data())
			continue;

		scores[std::round(score)].push_back(std::move(row));
	}

	auto min = std::min_element(scores.cbegin(), scores.cend(), [](auto&& lhs, auto&& rhs) {
		return lhs.second.size() < rhs.second.size();
	});

	for (auto&& score : scores) {
		auto &vec = score.second;

		auto progress = (vec.size() - min->second.size()) * max_diff;
		auto new_length = std::min<unsigned>(min->second.size() + progress, vec.size());
		std::mt19937_64 rnd{std::random_device{}()};

		std::shuffle(vec.begin(), vec.end(), rnd);
		std::for_each(vec.begin(), vec.begin() + new_length, [](auto &line) { fmt::print("{}\n", line); });
	}
}
