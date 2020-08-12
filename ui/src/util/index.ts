export function formatDuration(d: number[]): string {
    const [hours, minutes, seconds] = d

    let duration = ''

    if (hours > 0) {
        duration = `${hours} hour${hours > 1 ? 's' : ''}, `
    }

    if (minutes > 0) {
        duration += `${minutes} minute${minutes > 1 ? 's' : ''}, `
    }

    duration += `${seconds} seconds`

    return duration
}

export function getDurationInSeconds(startedAt: string, completedAt?: string): number {
    const now = completedAt ? new Date(completedAt).getTime() : new Date().getTime()
    return (now - new Date(startedAt).getTime()) / 1000
}

export function getDuration(startedAt: string, completedAt?: string): number[] {
    let delta = getDurationInSeconds(startedAt, completedAt)

    const hours = Math.floor(delta / 3600) % 24
    delta -= hours * 3600

    const minutes = Math.floor(delta / 60) % 60
    delta -= minutes * 60

    const seconds = Math.floor(delta)

    return [hours, minutes, seconds]
}

export function twoDecimals(n: number): number {
    const log10 = n ? Math.floor(Math.log10(n)) : 0
    const div = log10 < 0 ? Math.pow(10, 1 - log10) : 100
    return Math.round(n * div) / div
}

export function enrichTasks(tasks: Record<string, Record<string, any>>): Record<string, any>[] {
    return Object.values(tasks || {}).map((task) => ({
        ...task,
        duration: getDuration(task.startedAt, task.complete ? task.completedAt : undefined).join('.'),
    }))
}
