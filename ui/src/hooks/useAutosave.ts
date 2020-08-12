import { useEffect, useRef, useState } from 'react'

import useAppState from './useAppState'
import { secondsSince } from '../util'

export default function useAutosave(onSave: () => boolean): [string, (v: InputEvent) => void] {
    const { state } = useAppState()
    const lastSaved = useRef<number | undefined>(undefined)
    const [rate, setRate] = useState<string>('1m')

    const { complete } = state

    const duration = parseInt(rate, 10) * 60

    useEffect(() => {
        if (complete || isNaN(duration) || duration < 1) {
            return
        }

        if (lastSaved.current && secondsSince(lastSaved.current) < duration) {
            return
        }

        if (onSave()) {
            lastSaved.current = new Date().getTime()
        }
    }, [complete, duration, lastSaved.current, onSave])

    useEffect(() => {
        if (complete) {
            onSave()
            lastSaved.current = undefined
        }
    }, [complete])

    const onRateChange = (e: InputEvent) => {
        const target = e.target as HTMLInputElement
        setRate(target.value)
    }

    return [rate, onRateChange]
}
