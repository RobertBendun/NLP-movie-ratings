#include <algorithm>
#include <array>
#include <cassert>
#include <charconv>
#include <cmath>
#include <fstream>
#include <iomanip>
#include <iostream>
#include <libgen.h>
#include <tuple>
#include <vector>

#include <fmt/format.h>

using namespace fmt::literals;

int main(int, char **argv)
{
	bool print_equal = getenv("EQUAL");

	char const* expected = *++argv;
	char const* received = *++argv;
	if (!expected || !received) {
		std::cerr << "usage: " << basename(argv[-2]) << " <expected> <received>" << std::endl;
		return 1;
	}

	std::ifstream exp_file(expected), recv_file(received);
	if (!exp_file)  { std::cerr << "Cannot open file: " << expected << std::endl; std::exit(1); }
	if (!recv_file) { std::cerr << "Cannot open file: " << received << std::endl; std::exit(1); }

	double mse = 0, min_diff = 10, max_diff = 0, avg_diff = 0;
	unsigned count = 0;

	for (std::string exp_line, recv_line;
		std::getline(exp_file, exp_line) && std::getline(recv_file, recv_line);
	) {
		double exp, recv;
		std::from_chars(exp_line.data(), exp_line.data() + exp_line.size(), exp);
		std::from_chars(recv_line.data(), recv_line.size() + recv_line.data(), recv);
		auto diff = std::abs(exp - recv);

		if (print_equal && diff <= std::numeric_limits<double>::epsilon())
			fmt::print("Equal:\nExpected: {}\nReceived: {}\n", exp_line, recv_line);

		mse += diff * diff;
		min_diff = std::min(min_diff, diff);
		max_diff = std::max(max_diff, diff);
		avg_diff += diff;
		++count;
	}

	fmt::print("MSE: {}\n", mse / count);
	fmt::print(R"(Difference
  Min: {}
  Max: {}
  Avg: {}
)", min_diff, max_diff, avg_diff / count);
}
