#!/usr/bin/env python3
import matplotlib.pyplot as plt
import json

experiment = 'target3'
scales = ['scale2', 'scale3', 'scale4']

def get_scale_historical_data(scale):
    filename = experiment + '/' + scale + '/summary.json'
    try:
        with open(filename, 'r') as f:
            data = json.load(f)
    except FileNotFoundError:
        print(filename)
        return None
    return data

def get_plotable_minutes(duration):
    [hours, minutes, seconds] = duration.split('.')
    return int(hours) * 60 + int(minutes) + (int(seconds) / 60)

def transform_historical_data(data):
    def transform_data_point(point):
        transformed = point
        transformed['time'] = get_plotable_minutes(point['time'])
        return transformed
    return [transform_data_point(point) for point in data]

def get_historical_data():
    data = {}
    for scale in scales:
        scale_data = get_scale_historical_data(scale)
        if scale_data == None:
            continue
        data[scale] = transform_historical_data(scale_data['historicalData'])
    return data

def get_plot_data(series):
    time = [point['time'] for point in series]
    fitness = [point['fitness'] for point in series]
    return (time, fitness)

def plot_fitness():
    data = get_historical_data()

    labels = {
        'scale1': 'single system',
        'scale2': 'simulated',
        'scale3': 'distributed w/ 10s timeout',
        'scale4': 'distributed w/ 5s timeout'
    }

    for scale in data:
        (time, fitness) = get_plot_data(data[scale])
        plt.plot(time, fitness, label=labels[scale])

    plt.xlabel('Time (Minutes)')
    plt.ylabel('Fitness')
    plt.title('Experiment Fitness Over Time')
    plt.legend()
    plt.savefig('plots/fitness_by_time.png')

plot_fitness()
