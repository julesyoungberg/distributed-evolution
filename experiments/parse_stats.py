#!/usr/bin/env python3
import json
import re

experiments = ["target1", "target2", "target3"]
scales = ["scale1", "scale2"]


def normalized_memory(memory):
    m = re.match(r"(\d*\.?\d*)([a-zA-z]+)", memory)
    if m == None:
        return memory
    (value, unit) = m.groups()
    value = float(value)
    if unit == "GB":
        value *= 1000
    elif unit == "KB":
        value /= 1000
    elif unit == "B":
        value /= 1000000
    return str(value) + "MB"


def parse_line(line):
    parts = line.split()
    return {
        "name": parts[1].replace("distributed-evolution_", "").replace("_", "-"),
        "cpuPercent": parts[2],
        "memory": parts[3],
        "memoryPercent": parts[6],
        "netInput": normalized_memory(parts[7]),
        "netOutput": normalized_memory(parts[9]),
    }


def parse_file(filename):
    data = {}
    with open(filename, "r") as file:
        lines = file.readlines()

    for line in lines:
        if re.search("CONTAINER ID", line):
            continue

        line_data = parse_line(line)
        name = line_data["name"]

        if name in data:
            for key in line_data:
                if key == "name":
                    continue
                if line_data[key] > data[name][key]:
                    data[name][key] = line_data[key]
        else:
            data[name] = {}
            for key in line_data:
                if key == "name":
                    continue
                data[name][key] = line_data[key]

    return data


def process_experiment_stats(experiment, scale):
    path = experiment + "/" + scale + "/"
    data = parse_file(path + "stats")
    outFilenmae = path + "stats.json"
    with open(outFilenmae, "w") as file:
        json.dump(data, file, indent=4)


def parse_stats():
    for experiment in experiments:
        for scale in scales:
            try:
                process_experiment_stats(experiment, scale)
            except FileNotFoundError:
                continue


parse_stats()
