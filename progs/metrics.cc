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

#include <fmt/format.h>

using u64 = unsigned long long;

std::string human(u64 v)
{
	char buf[std::numeric_limits<u64>::digits10] = {};
	auto res = std::to_chars(buf, sizeof(buf) + buf, v);
	assert(res.ptr != buf);
	auto const len = res.ptr - buf;

	std::string result;
	result.resize(len * 4 / 3 - (len % 3 == 0), ' ');
	for (unsigned i = 0; i < len; ++i) {
		result[result.size() - i - 1 - i/3] = buf[len - i - 1];
	}
	return result;
}

unsigned positive = 0, total = 0;

// Index [1,10]
u64 true_positive[11] = {};
u64 false_positive[11] = {};
u64 false_negative[11] = {};
u64 count[11] = {};

double div(double a, double b, double if_zero) {
	return b == 0 ? if_zero : a / b;
}

// https://medium.com/apprentice-journal/evaluating-multi-class-classifiers-12b2946e755b
int main(int, char **argv)
{
	char const* expected = *++argv;
	char const* received = *++argv;
	if (!expected || !received) {
		std::cerr << "usage: " << basename(argv[-2]) << " <expected> <received>" << std::endl;
		return 1;
	}

	std::ifstream exp_file(expected), recv_file(received);
	if (!exp_file)  { std::cerr << "Cannot open file: " << expected << std::endl; std::exit(1); }
	if (!recv_file) { std::cerr << "Cannot open file: " << received << std::endl; std::exit(1); }

	for (std::string exp_line, recv_line;
		std::getline(exp_file, exp_line) && std::getline(recv_file, recv_line);
	) {
		unsigned exp;
		float recv_f;
		unsigned recv;

		std::from_chars(exp_line.data(), exp_line.data() + exp_line.size(), exp);
		std::from_chars(recv_line.data(), recv_line.size() + recv_line.data(), recv_f);
		recv = std::round(recv_f);

		++total;
		++count[exp];
		if (exp == recv) {
			++positive;
			true_positive[exp]++;
		} else {
			false_positive[recv]++;
			false_negative[exp]++;
		}
	}

	fmt::print("Accuracy: {:.2f}% (equal: {}, total: {})\n", 100.0 * positive / total, human(positive), human(total));

	const std::array<std::pair<char const*, u64*>, 2> per_class_metrics = {{
		{ "Precision", false_positive },
		{ "Recall", false_negative }
	}};


	for (auto metric : per_class_metrics) {
		fmt::print("{}:\n", metric.first);
		double avg = 0;
		for (unsigned i = 1; i <= 10; ++i) {
			auto const v = div(true_positive[i], true_positive[i] + metric.second[i], 0);
			fmt::print("  {}: {:.2f}% (equal: {}, total: {})\n", i, 100.0 * v, human(true_positive[i]), human(count[i]));
			avg += v;
		}
		fmt::print("  Average: {:.2f}%\n", 100 * avg / 10);
	}
}
