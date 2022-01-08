#!/usr/bin/env python3
import matplotlib.pyplot as plt, sys, os.path, seaborn as sns

def error(*args, end='\n'):
    sys.stderr.write(''.join(args) + end)

def usage():
    error(f"usage: {os.path.basename(__file__)} <file>")
    error("  where <file> is Vovpal Wabbit input format")
    exit(1)

if len(sys.argv) != 2:
    error("[ERROR] Expected filename, got nothing" if len(sys.argv) == 1
            else "[ERROR] Only one filename is supported")
    usage()


filename = sys.argv[1]

with open(filename) as f:
    data = [float(row[0]) for row in map(lambda x: x.split(' '), f.readlines()) if len(row) > 0]
    base, _ = os.path.splitext(filename)
    sns.displot(data).savefig(base + '.png')
