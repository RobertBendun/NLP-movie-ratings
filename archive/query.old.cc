#include <algorithm>
#include <cassert>
#include <charconv>
#include <filesystem>
#include <fstream>
#include <iostream>
#include <iterator>
#include <list>
#include <vector>

#include <fmt/format.h>
#include <range/v3/view.hpp>
#include <range/v3/algorithm.hpp>

namespace fs = std::filesystem;
using namespace fmt::literals;
using namespace std::string_view_literals;

template<typename T>
concept Stringable = requires (T const& a) { { a.to_string() }; };

namespace fmt
{
	template<Stringable T>
	struct formatter<T> : formatter<std::string_view>
	{
		auto format(T const& t, auto &ctx) { return format_to(ctx.out(), "{}", t.to_string()); }
	};
}

fs::path program_name;
fs::path database_file;
fs::path script_file;

void error_fatal(fs::path const& location, auto const& message)
{
	fmt::print("{}: error: {}", location.c_str(), message);
	std::exit(1);
}

void error_fatal(auto const& message)
{
	error_fatal(program_name, message);
}

void ensure_fatal(bool condition, auto const& message)
{
	if (condition) return;
	error_fatal(message);
}

void ensure_fatal(bool condition, fs::path const& location, auto const& message)
{
	if (condition) return;
	error_fatal(location, message);
}

namespace db
{
	enum class Field_Type
	{
		String,
		Float,
		Natural,
		__Last = Natural
	};

	static constexpr auto Field_Type_Strings = std::array {
		"string"sv,
		"float"sv,
		"nat"sv
	};

	struct Field_Descriptor
	{
		fs::path file;
		std::string name;
		Field_Type type;
		unsigned column = 0;

		auto to_string() const
		{
			return fmt::format("Field {}::{} : {}", file.c_str(), name, Field_Type_Strings[(int)type]);
		}
	};

	static std::string_view load_cell(std::string_view &line, std::string_view name)
	{
		ensure_fatal(!line.empty(), "Expected field: {}"_format(name));
		auto tab = line.find_first_of("\t\n");
		auto result = line.substr(0, tab);
		line.remove_prefix(tab+1);
		return result;
	}

	std::vector<Field_Descriptor> load_file(fs::path const& file)
	{
		std::ifstream src_file(file);
		if (!src_file) {
			error_fatal("cannot open database file `{}`"_format(file.c_str()));
		}

		// Skip (required) header
		{
			std::string header;
			ensure_fatal(bool(std::getline(src_file, header)), file, "Expected header");
		}

		std::vector<Field_Descriptor> descriptors;

		for (std::string linebuf; std::getline(src_file, linebuf); ) {
			auto &fd = descriptors.emplace_back();
			std::string_view line = linebuf;

			if (ranges::all_of(line, (int(*)(int))std::isspace))
				continue;

			fd.file = load_cell(line, "filename");
			fd.name = load_cell(line, "name");
			auto const type = load_cell(line, "type");

			static_assert((int)Field_Type::__Last == 2);
			bool success = false;
			for (auto [idx, name] : ranges::views::enumerate(Field_Type_Strings)) {
				if (name == type) {
					fd.type = (Field_Type)idx;
					success = true;
				}
			}
			ensure_fatal(success, "Unexpected value of field `type` = {}"_format(type));
		}

		std::unordered_map<std::string, unsigned> columns;

		for (auto &fd : descriptors) {
			auto const &file = fd.file.string();
			if (columns.contains(file))
				fd.column = ++columns[file];
			else {
				fd.column = columns[file] = 0;
			}
		}

		return descriptors;
	}
}

namespace lisp
{
	struct Context
	{
		std::vector<db::Field_Descriptor> fields;

		db::Field_Descriptor const& resolve(std::string_view path, std::string_view name) const
		{
			assert(!fields.empty());

			db::Field_Descriptor const* found = nullptr;

			for (auto const& fd : fields) {
				if (fd.file == path && fd.name == name)
					return fd;

				if ((path.empty() || std::string_view(fd.file.c_str()).find(path) != std::string_view::npos) && name == fd.name) {
					ensure_fatal(!found, script_file, "Symbol path=`{}` name=`{}` is ambigious"_format(path, name));
					found = &fd;
				}
			}

			ensure_fatal(found, script_file, "Cannot found field with path=`{}` name=`{}`"_format(path, name));
			return *found;
		}
	};

	struct Value
	{
		enum class Type
		{
			Nil,
			Symbol,
			List,
			Field
		};

		static Value nil() { return { Type::Nil }; }
		static Value symbol(std::string_view sv) { return { Type::Symbol, { sv.data(), sv.size() } }; }

		Type type = {};
		std::string sval = {};
		std::list<Value> list = {};
		db::Field_Descriptor const* fd = nullptr;

		std::string to_string() const
		{
			switch (type) {
			case Type::Nil:
				return "nil";
			case Type::Symbol:
				return sval;
			case Type::List:
				return fmt::format("({})", fmt::join(list, " "));
			case Type::Field:
				assert(fd);
				return fmt::format("<field {}/{}>", fd->file.c_str(), fd->name);
			}

			assert(false && "unreachable");
		}
	};

	Value read(std::string_view &source)
	{
		while (!source.empty()) {
			for (; !source.empty() && std::isspace(source.front()); source.remove_prefix(1)) {}
			if (source.starts_with('#')) {
				source.remove_prefix(source.find('\n'));
			} else {
				break;
			}
		}

		if (source.empty())
			return Value::nil();

		static constexpr std::string_view Valid_Symbol_Char = "_+-*/%$@!^&[]:;<>,.|=";
		if (std::isalpha(source.front()) || Valid_Symbol_Char.find(source.front()) != std::string_view::npos) {
			auto end = ranges::find_if(source.substr(1), [](char curr) {
				return !std::isalnum(curr) && Valid_Symbol_Char.find(curr) == std::string_view::npos;
			});
			auto symbol = Value::symbol({ std::cbegin(source), end });
			source.remove_prefix(std::distance(std::cbegin(source), end));
			return symbol;
		}

		if (source.starts_with('(')) {
			Value list, elem;
			list.type = Value::Type::List;
			source.remove_prefix(1);
			while ((elem = read(source)).type != Value::Type::Nil) {
				list.list.push_back(std::move(elem));
			}
			return list;
		}

		if (source.starts_with(')')) {
			source.remove_prefix(1);
		}

		return Value::nil();
	}

	Value eval(Context &ctx, Value value)
	{
		switch (value.type) {
		case Value::Type::Field:
		case Value::Type::Nil:
			return value;

		case Value::Type::Symbol:
			{
				auto symbol = std::string_view(value.sval);
				std::string_view file, field;
				if (auto pos = symbol.find('/'); pos != std::string_view::npos) {
					file = symbol.substr(0, pos);
					field = symbol.substr(pos+1);
				} else {
					field = symbol;
				}

				Value val;
				val.type = Value::Type::Field;
				val.fd = &ctx.resolve(file, field);
				return val;
			}
			break;

		case Value::Type::List:
			{
				if (value.list.empty())
					return Value::nil();

				auto func = value.list.front();
				if (func.type != Value::Type::Symbol)
					return value;

				if (func.sval == "list") {
					value.list.pop_front();
					for (auto &arg : value.list)
						arg = eval(ctx, arg);
					return value;
				}
			}
		}

		return Value::nil();
	}
}

void usage()
{
	std::cerr << "usage: " << program_name.c_str() << " [-d <database.tsv>] [-c <query>] [script]\n"
		"  where\n"
		"    -h,--help    Print this message\n"
		"    -c query     Uses query string provided in command line rather then script\n"
		"    -d database  Provides path to the database description file in TSV format\n"
		"     script      File containing query to execute\n"
		"\n"
		"  Debugging options\n"
		"    --schema     Print loaded fields from database schema and quit\n"
		"    --query      Print loaded query from script file or -c argument and quit\n";
	std::exit(1);
}

inline auto invalid_argument_check(std::convertible_to<bool> auto const& val, char const* message)
{
	ensure_fatal(bool(val), "Invalid argument: {}"_format(message));
	return val;
}

int main(int, char **argv)
{
	program_name = fs::path(*argv).filename();

	std::string_view source;

	bool print_db_fields = false;
	bool print_query = false;

	while (*++argv != nullptr) {
		if (*argv == "-h"sv || *argv == "--help"sv || *argv == "/?"sv)
			usage();
		if (*argv == "-d"sv) {
			database_file = invalid_argument_check(*++argv, "-d expects path to database!");
			continue;
		}
		if (*argv == "-c"sv) {
			script_file = "<arg>";
			source = invalid_argument_check(*++argv, "-c expects argument with a LISP query!");
			continue;
		}
		if (*argv == "--schema"sv) { print_db_fields = true; continue; }
		if (*argv == "--query"sv) { print_query = true; continue; }
		if (script_file.empty()) {
			script_file = *argv;
		} else {
			invalid_argument_check(false, "Multiple script file ware passed");
		}
	}
	if (source.empty() && script_file.empty()) {
		std::cerr << program_name.c_str() << ": error: Missing argument:"
			" Query must be passed either by providing script source or -c option.\n";
		return 1;
	}
	if (database_file.empty()) {
		std::cerr << program_name.c_str() << ": error: Missing argument:"
			" Missing database file argument (-d).\n";
		return 1;
	}

	auto database = db::load_file(database_file);
	if (print_db_fields) {
		fmt::print("Database schema\n  {}\n", fmt::join(database, "\n  "));
		return 0;
	}

	auto query = lisp::read(source);
	if (print_query) {
		fmt::print("Query\n  {}\n", query);
		return 0;
	}

	lisp::Context ctx;
	ctx.fields = database;
	auto result = lisp::eval(ctx, query);

	if (ranges::all_of(result.list, [](auto const& p) { return p.type == lisp::Value::Type::Field; })) {
		std::unordered_map<
	}
}
