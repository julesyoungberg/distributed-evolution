#!/usr/bin/env python3
import matplotlib.pyplot as plt
import json
import re

labels = {
    'scale1': 'single system',
    'scale2': 'simulated',
    'scale3': 'distributed',
}

def regroup(averages):
    data = {}
    for scale in averages:
        for metric in averages[scale]:
            val = averages[scale][metric]
            if metric in data:
                data[metric][scale] = val
            else:
                data[metric] = { scale: val }
    return data

def get_labels(scales):
    return [labels[scale] for scale in scales]

def get_duration_minutes(duration):
    [hours, minutes, seconds] = duration.split('.')
    return int(minutes) + int(hours) * 60 + (int(seconds) / 60)

def plot_duration(durations):
    scales = get_labels(durations.keys())
    durations = [get_duration_minutes(dur) for dur in durations.values()]

    fig = plt.figure()
    ax = fig.add_subplot(111)
    ax.set_title('Average Duration by Scale')
    ax.set_ylabel('Duration (minutes)')
    ax.bar(scales, durations)
    plt.savefig('plots/duration.png')

def plot_generations_per_second(gensPerSeconds):
    scales = get_labels(gensPerSeconds.keys())
    values = gensPerSeconds.values()

    fig = plt.figure()
    ax = fig.add_subplot(111)
    ax.set_title('Average Generations Per Second by Scale')
    ax.set_ylabel('Generations Per Second')
    ax.bar(scales, values)
    plt.savefig('plots/generations_per_second.png')

def plot_fitness(fitness):
    scales = get_labels(fitness.keys())
    values = fitness.values()

    fig = plt.figure()
    ax = fig.add_subplot(111)
    ax.set_title('Average Fitness by Scale')
    ax.set_ylabel('Fitness')
    ax.bar(scales, values)
    plt.savefig('plots/fitness.png')

def get_memory_value(memory):
    return float(re.search(r'(\d*\.?\d*)', memory).groups()[0])

def plot_memory(memory):
    scales = get_labels(memory.keys())
    values = [get_memory_value(mem) for mem in memory.values()]

    fig = plt.figure()
    ax = fig.add_subplot(111)
    ax.set_title('Average Memory by Scale')
    ax.set_ylabel('Memory (Megabytes)')
    ax.bar(scales, values)
    plt.savefig('plots/memory.png')

def plot_averages():
    with open('averages.json', 'r') as f:
        averages = json.load(f)

    data = regroup(averages)

    plot_duration(data['duration'])
    plot_generations_per_second(data['generationsPerSecond'])
    plot_fitness(data['fitness'])
    plot_memory(data['memory'])

plot_averages()
