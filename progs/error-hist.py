#!/usr/bin/env python3
import matplotlib.pyplot as plt, sys, os.path, seaborn as sns

def error(*args, end='\n'):
    sys.stderr.write(''.join(args) + end)

def usage():
    error(f"usage: {os.path.basename(__file__)} <expected> <received>")
    error("  where <expected> and <received> are Vovpal Wabbit input format")
    exit(1)

if len(sys.argv) != 3:
    error("[ERROR] Expected filename, got nothing" if len(sys.argv) == 1
            else "[ERROR] Two filenames are required")
    usage()


def get_prediction(lines : list[str]) -> list[float]:
    return [float(row[0])
        for row in map(lambda x: x.split(' '), lines)
        if len(row) > 0]

with open(sys.argv[1]) as exp, open(sys.argv[2]) as recv:
    expected, received = get_prediction(exp.readlines()), get_prediction(recv.readlines())
    data = [abs(e - r) for e, r in zip(expected, received)]

    base, _ = os.path.splitext(sys.argv[1])
    sns.displot(data).savefig(base + '.png')
