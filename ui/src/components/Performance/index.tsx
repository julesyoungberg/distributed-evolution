/** @jsx jsx */
import { jsx } from '@emotion/core'
import styled from '@emotion/styled'
import { useEffect, useState } from 'react'
import { Box, Flex, Text } from 'rebass'
import { LineChart, Line, XAxis, YAxis } from 'recharts'

import useAppState from '../../hooks/useAppState'
import useTheme from '../../hooks/useTheme'
import { twoDecimals } from '../../util'

const StyledText = styled(Text)`
    margin-top: 5px;
    margin-bottom: 5px;
`

function getDuration(startedAt: string): number[] {
    let delta = (new Date().getTime() - new Date(startedAt).getTime()) / 1000

    const hours = Math.floor(delta / 3600) % 24
    delta -= hours * 3600

    const minutes = Math.floor(delta / 60) % 60
    delta -= minutes * 60

    const seconds = Math.floor(delta)

    return [hours, minutes, seconds]
}

function formatDuration(d: number[]): string {
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

interface DataPoint {
    fitness: number
    generation: number
    time: string
}

interface ChartProps {
    color: string
    data: DataPoint[]
    dataKey: string
    label: string
}

function Chart({ color, data, dataKey, label }: ChartProps) {
    return (
        <Box>
            <Text css={{ textAlign: 'center' }} fontSize={[2, 3, 4]}>{label}</Text>
            
            <LineChart data={data} width={600} height={400}>
                <XAxis 
                    dataKey='time' 
                    height={40}
                    interval='preserveStartEnd'
                    minTickGap={20}
                    tickCount={5}
                />
                <YAxis width={80} />
                <Line 
                    dataKey={dataKey}
                    dot={false}
                    type='monotone'
                    stroke={color} 
                />
            </LineChart>
        </Box>
    )
}

export default function Performance() {
    const [data, setData] = useState<DataPoint[]>([])
    const { state } = useAppState()
    const theme = useTheme()

    const { fitness, generation, startedAt, status } = state

    const duration = getDuration(startedAt)

    useEffect(() => {
        if (!(fitness && generation && status === 'active')) {
            return
        }

        setData([...data, { fitness, generation, time: duration.join('.') }])
    }, [fitness, generation])

    if (!(fitness && generation && startedAt)) {
        return null
    }

    return (
        <Box css={{ marginTop: 50, marginBottom: 50 }}>
            <Box css={{ marginBottom: 40 }}>
                <Text fontSize={[2, 3, 4]}>
                    Performance
                </Text>
                <StyledText><b>Generation:</b> {generation}</StyledText> 
                <StyledText><b>Fitness:</b> {twoDecimals(fitness)}</StyledText>
                {startedAt && <StyledText><b>Duration:</b> {formatDuration(duration)}</StyledText>}
            </Box>

            <Flex>
                <Box>
                    <Chart
                        color={theme.colors.blue}
                        data={data}
                        dataKey='fitness'
                        label='Fitness'
                    />
                </Box>

                <Box>
                    <Chart
                        color={theme.colors.primary}
                        data={data}
                        dataKey='generation'
                        label='Generation'
                    />
                </Box>
            </Flex>
        </Box>
    )
}
