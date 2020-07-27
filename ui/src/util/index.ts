export function twoDecimals(n: number): number {
    const log10 = n ? Math.floor(Math.log10(n)) : 0
    const div = log10 < 0 ? Math.pow(10, 1 - log10) : 100
    return Math.round(n * div) / div
}
