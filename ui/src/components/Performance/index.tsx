/** @jsx jsx */
import { jsx } from '@emotion/core'
import styled from '@emotion/styled'
import { Label, Select } from '@rebass/forms'
import download from 'downloadjs'
import { useCallback, useEffect, useState } from 'react'
import { Box, Flex, Text } from 'rebass'
import { LineChart, Line, XAxis, YAxis } from 'recharts'

import useAppState from '../../hooks/useAppState'
import useAutosave from '../../hooks/useAutosave'
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

    const { complete, completedAt, fitness, generation, jobID, startedAt, status } = state

    const duration = generation ? getDuration(startedAt, complete ? completedAt : undefined) : undefined

    useEffect(() => {
        if (complete || !(fitness && generation && status === 'active')) {
            return
        }

        setData([...data, { fitness, generation, time: duration.join('.') }])
    }, [fitness, generation, complete])

    const onClear = () => setData([])

    const onSave = useCallback(() => {
        if (!generation) {
            return false
        }

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

        download(jsonString, `${jobID}-${generation}-${twoDecimals(fitness)}.json`)
        return true
    }, [data, duration, state])

    const [rate, onRateChange] = useAutosave(onSave)

    return (
        <Box css={{ marginTop: 50, marginBottom: 50 }}>
            <Box css={{ marginBottom: 40 }}>
                <Text fontSize={[2, 3, 4]}>Performance</Text>
                {generation > 0 && (
                    <StyledText>
                        <b>Generation:</b> {generation}
                    </StyledText>
                )}
                {fitness > 0 && (
                    <StyledText>
                        <b>Fitness:</b> {twoDecimals(fitness)}
                    </StyledText>
                )}
                {duration && (
                    <StyledText>
                        <b>Duration:</b> {formatDuration(duration)}
                    </StyledText>
                )}
                <Box>
                    <Button css={{ marginRight: '10px' }} disabled={data.length === 0} onClick={onClear}>
                        Clear Data
                    </Button>
                    <Button css={{ marginRight: '10px' }} disabled={data.length === 0} onClick={onSave}>
                        Save Data
                    </Button>
                    <Box css={{ display: 'inline-block', width: '100px' }}>
                        <Label htmlFor='dataSaveRate'>Save Rate</Label>
                        <Select id='dataSaveRate' name='dataSaveRate' value={rate} onChange={onRateChange}>
                            {['none', '1m', '5m', '10m', '15m', '30m', '60m', '90m', '120m'].map((type) => (
                                <option key={type}>{type}</option>
                            ))}
                        </Select>
                    </Box>
                </Box>
            </Box>

            {data.length > 0 && (
                <Flex>
                    <Box>
                        <Chart color={theme.colors.blue} data={data} dataKey='fitness' label='Fitness' />
                    </Box>

                    <Box>
                        <Chart color={theme.colors.primary} data={data} dataKey='generation' label='Generation' />
                    </Box>
                </Flex>
            )}
        </Box>
    )
}
