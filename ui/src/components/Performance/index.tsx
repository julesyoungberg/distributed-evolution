/** @jsx jsx */
import { jsx } from '@emotion/core'
import styled from '@emotion/styled'
import { Box, Text } from 'rebass'

import useAppState from '../../hooks/useAppState'

const Line = styled(Text)`
    margin-top: 5px;
    margin-bottom: 5px;
`

function getDuration(startedAt: string): string {
    let delta = (new Date().getTime() - new Date(startedAt).getTime()) / 1000

    const hours = Math.floor(delta / 3600) % 24
    delta -= hours * 3600

    const minutes = Math.floor(delta / 60) % 60
    delta -= minutes * 60

    const seconds = Math.floor(delta)

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

export default function Performance() {
    const { state } = useAppState()
    const { fitness, generation, startedAt, status } = state

    return (
        <Box css={{ marginTop: 50, marginBottom: 50 }}>
            <Text fontSize={[2, 3, 4]}>
                Performance
            </Text>
            <Line><b>Generation:</b> {generation}</Line> 
            <Line><b>Fitness:</b> {fitness}</Line>
            {startedAt && generation > 0 && status === 'active' && (
                <Line><b>Duration:</b> {getDuration(startedAt)}</Line>
            )}
        </Box>
    )
}
