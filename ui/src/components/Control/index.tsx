/** @jsx jsx */
import { jsx } from '@emotion/core'
import styled from '@emotion/styled'
import { Input, Label, Select } from '@rebass/forms'
import fetch from 'isomorphic-fetch'
import getConfig from 'next/config'
import { useState, FormEvent } from 'react'
import { Box, Button, Flex } from 'rebass'

import useAppState from '../../hooks/useAppState'
import { Theme } from '../../theme'

const { publicRuntimeConfig } = getConfig()

const StyledButton = styled(Button)<{ disabled: boolean }, Theme>`
    padding: 10px 20px;
    cursor: pointer;
    font-weight: 700;
    text-transform: uppercase;
    background-color: ${({ disabled, theme }) => (disabled ? theme.colors.lightgray : theme.colors.primary)};

    &:focus {
        outline: none;
    }
`

const Field = styled(Box)`
    width: 250px;
    margin: 10px 0;
`

interface Config {
    shapeType: 'circles' | 'polygons' | 'triangles'
    numShapes: number
    shapeSize: number
    popSize: number
    poolSize: number
    mutationRate: number
    crossRate: number
    overDraw: number
}

const initialConfig: Config = Object.freeze({
    shapeType: 'polygons',
    numShapes: 7000,
    shapeSize: 30,
    popSize: 50,
    poolSize: 10,
    mutationRate: 0.02,
    crossRate: 0.2,
    overDraw: 20,
})

function getBase64Image(img: HTMLImageElement) {
    const canvas = document.createElement('canvas')
    canvas.width = img.width
    canvas.height = img.height

    const ctx = canvas.getContext('2d')
    ctx.drawImage(img, 0, 0)

    const dataURL = canvas.toDataURL('image/png')
    return dataURL.replace(/^data:image\/(png|jpg);base64,/, '')
}

export default function Control() {
    const { dispatch, state } = useAppState()
    const [loading, setLoading] = useState<boolean>(false)
    const [config, setConfig] = useState<Config>(initialConfig)

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

        const body = { ...config, targetImage: state.nextTargetImage }
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
        } catch(e) {
            data = e
        }

        console.log('data', data)

        dispatch({
            type: 'start', 
            payload: { data, statusCode: response.status },
        })

        setLoading(false)
    }

    const fieldProps = (name: string) => ({
        id: name,
        name: name,
        value: config[name],
        onChange: (e: InputEvent) => {
            const target = e.target as HTMLInputElement
            console.log(Object.keys(target))
            setConfig({
                ...config,
                [target.name]: target.value,
            })
        },
    })

    const disableButtons = loading || ['disconnected', 'error'].includes(state.status)

    return (
        <form css={{ marginBottom: 20 }}>
            <Flex css={{ textAlign: 'center' }} justifyContent='space-around'>
                <Box width={1 / 2}>
                    <StyledButton disabled={disableButtons} onClick={disableButtons ? undefined : getRangomTargetImage}>
                        Random Target Image
                    </StyledButton>
                </Box>
                <Box width={1 / 2}>
                    <StyledButton
                        disabled={disableButtons}
                        color='secondary'
                        onClick={disableButtons ? undefined : onStart}
                        type='submit'
                    >
                        Start
                    </StyledButton>
                </Box>
            </Flex>
            <Flex css={{ marginTop: '20px' }} flexWrap='wrap' justifyContent='space-between'>
                <Field>
                    <Label htmlFor='shapeType'>Shape Type</Label>
                    <Select {...fieldProps('shapeType')}>
                        {['circles', 'polygons', 'triangles'].map((type) => (
                            <option key={type}>{type}</option>
                        ))}
                    </Select>
                </Field>
                <Field>
                    <Label htmlFor='numShapes'>Number of Shapes per slice</Label>
                    <Input type='number' step='10' min='10' max='1000' {...fieldProps('numShapes')} />
                </Field>
                <Field>
                    <Label htmlFor='shapeSize'>Shape Size</Label>
                    <Input type='number' step='5' min='5' max='200' {...fieldProps('shapeSize')} />
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
            </Flex>
        </form>
    )
}
