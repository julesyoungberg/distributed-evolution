/** @jsx jsx */
import { jsx } from '@emotion/core'
import styled from '@emotion/styled'
import download from 'downloadjs'
import { useEffect, useState } from 'react'
import { Box, Flex, Text } from 'rebass'
import { LineChart, Line, XAxis, YAxis } from 'recharts'

import useAppState from '../../hooks/useAppState'
import useTheme from '../../hooks/useTheme'
import { formatDuration, getDuration, twoDecimals } from '../../util'

import Button from '../Button'

const StyledText = styled(Text)`
    margin-top: 5px;
    margin-bottom: 5px;
`

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
            <Text css={{ textAlign: 'center' }} fontSize={[2, 3, 4]}>
                {label}
            </Text>

            <LineChart data={data} width={600} height={400}>
                <XAxis dataKey='time' height={40} interval='preserveStartEnd' minTickGap={20} tickCount={5} />
                <YAxis width={80} />
                <Line dataKey={dataKey} dot={false} type='monotone' stroke={color} />
            </LineChart>
        </Box>
    )
}

export default function Performance() {
    const [data, setData] = useState<DataPoint[]>([])
    const { state } = useAppState()
    const theme = useTheme()

    const { complete, completedAt, fitness, generation, startedAt, status } = state

    const duration = getDuration(startedAt, complete ? completedAt : undefined)

    useEffect(() => {
        if (complete || !(fitness && generation && status === 'active')) {
            return
        }

        setData([...data, { fitness, generation, time: duration.join('.') }])
    }, [fitness, generation, complete])

    if (!(fitness && generation && startedAt)) {
        return null
    }

    const onSave = () => {
        const json = {
            metadata: { duration: duration.join('.'), ...state },
            historicalData: data,
        }

        // delete the images
        delete json.metadata.nextTargetImage
        delete json.metadata.output
        delete json.metadata.palette
        delete json.metadata.targetImage
        delete json.metadata.targetImageEdges

        const jsonString = `data:text/json;charset=utf-8,${encodeURIComponent(JSON.stringify(json, null, 4))}`

        download(jsonString, 'distributed-evolution-data.json')
    }

    const onClear = () => setData([])

    return (
        <Box css={{ marginTop: 50, marginBottom: 50 }}>
            <Box css={{ marginBottom: 40 }}>
                <Text fontSize={[2, 3, 4]}>Performance</Text>
                <StyledText>
                    <b>Generation:</b> {generation}
                </StyledText>
                <StyledText>
                    <b>Fitness:</b> {twoDecimals(fitness)}
                </StyledText>
                {startedAt && (
                    <StyledText>
                        <b>Duration:</b> {formatDuration(duration)}
                    </StyledText>
                )}
                <Flex>
                    <Box css={{ marginRight: '10px' }}>
                        <Button onClick={onSave}>Save Data</Button>
                    </Box>
                    <Box>
                        <Button onClick={onClear}>Clear Data</Button>
                    </Box>
                </Flex>
            </Box>

            <Flex>
                <Box>
                    <Chart color={theme.colors.blue} data={data} dataKey='fitness' label='Fitness' />
                </Box>

                <Box>
                    <Chart color={theme.colors.primary} data={data} dataKey='generation' label='Generation' />
                </Box>
            </Flex>
        </Box>
    )
}
