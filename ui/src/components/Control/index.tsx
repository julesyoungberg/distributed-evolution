/** @jsx jsx */
import { jsx } from '@emotion/core'
import styled from '@emotion/styled'
import { Checkbox, Input, Label, Select } from '@rebass/forms'
import fetch from 'isomorphic-fetch'
import download from 'downloadjs'
import getConfig from 'next/config'
import { FormEvent, useCallback, useRef, useState } from 'react'
import { Box, Flex } from 'rebass'

import useAppState from '../../hooks/useAppState'
import useAutosave from '../../hooks/useAutosave'
import { twoDecimals } from '../../util'

import Button from '../Button'

const { publicRuntimeConfig } = getConfig()

const Field = styled(Box)`
    width: 250px;
    margin: 10px 0;
`

interface Config {
    crossRate: number
    detectEdges: boolean
    mutationRate: number
    numColors: number
    numGenerations: number
    numShapes: number
    overDraw: number
    paletteType: 'random' | 'targetImage' | 'targetImageRandomCenters'
    poolSize: number
    popSize: number
    quantization: number
    shapeSize: number
    shapeType: 'circles' | 'lines' | 'polygons' | 'triangles'
    targetImage?: string
}

const initialConfig: Config = Object.freeze({
    crossRate: 0.2,
    detectEdges: false,
    mutationRate: 0.02,
    numColors: 64,
    numGenerations: 0,
    numShapes: 7000,
    overDraw: 20,
    paletteType: 'random',
    poolSize: 10,
    popSize: 50,
    quantization: 50,
    shapeSize: 20,
    shapeType: 'polygons',
})

function getBase64Image(img: HTMLImageElement) {
    const canvas = document.createElement('canvas')
    canvas.width = img.width
    canvas.height = img.height

    const ctx = canvas.getContext('2d')
    ctx.drawImage(img, 0, 0)

    const dataURL = canvas.toDataURL('image/png')
    return dataURL.replace(/^data:image\/(png|jpg|jpeg);base64,/, '')
}

export default function Control() {
    const { dispatch, state } = useAppState()
    const fileInputRef = useRef<HTMLInputElement | undefined>(undefined)
    const [loading, setLoading] = useState<boolean>(false)
    const [config, setConfig] = useState<Config>(initialConfig)

    const { fitness, generation, jobID, output, status } = state

    const onFileInputChange = () => {
        setLoading(true)
        dispatch({ type: 'clearTarget' })

        if (!fileInputRef.current) throw new Error('File input ref has no current value')

        const img = new Image()
        img.crossOrigin = 'anonymous'

        img.onload = () => {
            dispatch({ type: 'setTarget', payload: { target: getBase64Image(img) } })
            setLoading(false)
        }

        const file = fileInputRef.current.files[0]
        const reader = new FileReader()
        reader.addEventListener(
            'load',
            () => {
                img.src = reader.result as string
            },
            false
        )

        reader.readAsDataURL(file)
    }

    const uploadTargetImage = (event: MouseEvent) => {
        event.preventDefault()
        if (!fileInputRef.current) throw new Error('File input ref has no current value')
        fileInputRef.current.click()
    }

    const getRangomTargetImage = () => {
        setLoading(true)
        dispatch({ type: 'clearTarget' })

        const img = new Image()
        img.crossOrigin = 'anonymous'
        img.src = `https://picsum.photos/900?now=${Date.now()}`

        img.onload = () => {
            dispatch({ type: 'setTarget', payload: { target: getBase64Image(img) } })
            setLoading(false)
        }
    }

    const onStart = async (e: FormEvent) => {
        e.preventDefault()
        setLoading(true)

        const body = { ...config }
        body.targetImage = state.nextTargetImage || state.targetImage

        console.log('starting task body', body)

        const response = await fetch(`${publicRuntimeConfig.apiUrl}/job`, {
            method: 'POST',
            headers: { 'Content-Type': 'application/json' },
            body: JSON.stringify(body),
        })

        console.log('response', response)

        let data: any

        try {
            if (response.status === 200) {
                data = await response.json()
            } else {
                data = await response.text()
            }
        } catch (e) {
            data = e
        }

        dispatch({
            type: 'start',
            payload: { data, statusCode: response.status },
        })

        setLoading(false)
    }

    const onSave = useCallback(
        (event?: MouseEvent) => {
            if (!output) {
                return false
            }

            event?.preventDefault()
            download(
                `data:image/png;base64,${output}`,
                `${jobID}-${generation}-${twoDecimals(fitness)}.png`,
                'image/png'
            )
            return true
        },
        [output, jobID, generation, fitness]
    )

    const [rate, onRateChange] = useAutosave(onSave)

    const fieldProps = (name: string) => ({
        id: name,
        name: name,
        value: config[name],
        onChange: (e: InputEvent) => {
            const target = e.target as HTMLInputElement

            const key = target.name
            let value: any = target.value

            if (typeof initialConfig[key] === 'number') {
                if (Number.isInteger(initialConfig[key])) {
                    value = parseInt(value, 10)
                } else {
                    value = parseFloat(value)
                }
            }

            setConfig({ ...config, [key]: value })
        },
    })

    const onCheckboxChange = (e: InputEvent) => {
        const target = e.target as HTMLInputElement
        setConfig({ ...config, [target.name]: target.checked })
    }

    const disableButtons = loading || ['disconnected', 'error'].includes(status)

    return (
        <form css={{ marginBottom: 20 }}>
            <Flex css={{ textAlign: 'center' }} justifyContent='space-around'>
                <Box width={1 / 2}>
                    <Input css={{ display: 'none ' }} onChange={onFileInputChange} ref={fileInputRef} type='file' />
                    <Button css={{ marginRight: 10 }} disabled={disableButtons} onClick={uploadTargetImage}>
                        Upload Target
                    </Button>
                    <Button disabled={disableButtons} onClick={getRangomTargetImage}>
                        Random Target
                    </Button>
                </Box>
                <Box width={1 / 2}>
                    <Button
                        css={{ marginRight: 10 }}
                        disabled={!['active', 'editing'].includes(status)}
                        onClick={onStart}
                        type='submit'
                    >
                        Start
                    </Button>
                    <Button css={{ marginRight: 10 }} disabled={!output} onClick={onSave}>
                        Save
                    </Button>
                    <Box css={{ display: 'inline-block', width: '100px' }}>
                        <Label htmlFor='outputSaveRate'>Save Rate</Label>
                        <Select id='outputSaveRate' name='outputSaveRate' value={rate} onChange={onRateChange}>
                            {['none', '1m', '5m', '10m', '15m', '30m', '60m', '90m', '120m'].map((type) => (
                                <option key={type}>{type}</option>
                            ))}
                        </Select>
                    </Box>
                </Box>
            </Flex>
            <Flex css={{ marginTop: '20px' }} flexWrap='wrap' justifyContent='space-between'>
                <Field>
                    <Label htmlFor='shapeType'>Shape Type</Label>
                    <Select {...fieldProps('shapeType')}>
                        {['circles', 'lines', 'polygons', 'triangles'].map((type) => (
                            <option key={type}>{type}</option>
                        ))}
                    </Select>
                </Field>
                <Field>
                    <Label htmlFor='numShapes'>Number of Colors</Label>
                    <Input type='number' step='8' min='8' max='1024' {...fieldProps('numColors')} />
                </Field>
                <Field>
                    <Label htmlFor='numShapes'>Number of Shapes</Label>
                    <Input type='number' step='10' min='10' max='10000' {...fieldProps('numShapes')} />
                </Field>
                <Field>
                    <Label htmlFor='shapeSize'>Shape Size</Label>
                    <Input type='number' step='5' min='5' max='200' {...fieldProps('shapeSize')} />
                </Field>
                <Field>
                    <Label htmlFor='quantization'>Shape Position Quantization</Label>
                    <Input type='number' step='5' min='0' max='500' {...fieldProps('quantization')} />
                </Field>
                <Field>
                    <Label htmlFor='numGenerations'>Generations</Label>
                    <Input type='number' step='1000' min='0' max='1000000000' {...fieldProps('numGenerations')} />
                </Field>
                <Field>
                    <Label htmlFor='popSize'>Population Size</Label>
                    <Input type='number' step='5' min='5' max='200' {...fieldProps('popSize')} />
                </Field>
                <Field>
                    <Label htmlFor='poolSize'>Breeding Pool Size</Label>
                    <Input type='number' step='5' min='5' max='100' {...fieldProps('poolSize')} />
                </Field>
                <Field>
                    <Label htmlFor='mutationRate'>Mutation Rate</Label>
                    <Input step='0.001' min='0.0' max='1.0' {...fieldProps('mutationRate')} />
                </Field>
                <Field>
                    <Label htmlFor='crossRate'>Cross Rate</Label>
                    <Input type='number' step='0.001' min='0.0' max='1.0' {...fieldProps('crossRate')} />
                </Field>
                <Field>
                    <Label htmlFor='overDraw'>Over Draw</Label>
                    <Input type='number' step='1' min='0' max='100' {...fieldProps('overDraw')} />
                </Field>
                <Field>
                    <Label htmlFor='paletteType'>Palette Type</Label>
                    <Select {...fieldProps('paletteType')}>
                        {['random', 'targetImage', 'targetImageRandomCenters'].map((type) => (
                            <option key={type}>{type}</option>
                        ))}
                    </Select>
                </Field>
                <Field>
                    <Label>
                        <Checkbox
                            checked={config.detectEdges}
                            id='detectEdges'
                            name='detectEdges'
                            onChange={onCheckboxChange}
                            type='checkbox'
                        />
                        Detect Edges
                    </Label>
                </Field>
            </Flex>
        </form>
    )
}
