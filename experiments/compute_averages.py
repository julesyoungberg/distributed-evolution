#!/usr/bin/env python3
import json
import re

experiments = ['target1', 'target2', 'target3']
scales = ['scale1', 'scale2', 'scale3']

def load_data(path, filename):
    name = path + filename
    try:
        with open(name, 'r') as f:
            data = json.load(f)
    except FileNotFoundError:
        return None
    return data

def read_summary(path):
    data = load_data(path, 'summary.json')
    if data == None:
        return data
    return data['metadata']

def read_stats(path):
    return load_data(path, 'stats.json')

def get_duration(seconds):
    hours = int(seconds / 3600)
    seconds -= hours * 3600
    minutes = int(seconds / 60)
    seconds -= minutes * 60
    s = '.'.join(map(str, [hours, minutes, int(seconds)]))
    return s

def get_total_memory(stats):
    count = 0
    total = 0

    for key in stats:
        m = re.search(r'(\d*\.?\d*)', stats[key]['memory'])
        if m == None:
            continue
        count += 1
        total += float(m.groups()[0])

    if count == 0:
        return 0

    return round(total / count, 2)

def compute_averages():
    averages = {}

    for scale in scales:
        total = 0
        sums = { 
            'duration': 0, 
            'fitness': 0, 
            'generationsPerSecond': 0,
            'memory': 0,
        }

        for experiment in experiments:
            path = experiment + '/' + scale + '/'

            summary = read_summary(path)
            stats = read_stats(path)

            if summary == None and stats == None:
                continue

            total += 1

            if summary:
                [hours, minutes, seconds] = summary['duration'].split('.')
                seconds = int(seconds) + int(minutes) * 60 + int(hours) * 3600
                sums['duration'] += seconds
                sums['generationsPerSecond'] += float(summary['generation'] / seconds)
                sums['fitness'] += float(summary['fitness'])
            if stats:
                sums['memory'] += get_total_memory(stats)

        if total == 0:
            continue

        averages[scale] = {}
        averages[scale]['duration'] = get_duration(round(sums['duration'] / total))
        averages[scale]['generationsPerSecond'] = round(sums['generationsPerSecond'] / total, 2)
        averages[scale]['fitness'] = sums['fitness'] / total
        averages[scale]['memory'] = str(round(sums['memory'] / total, 2)) + 'MiB'

    return averages

averages = compute_averages()

with open('averages.json', 'w') as file:
    json.dump(averages, file, indent=4)
