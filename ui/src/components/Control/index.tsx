/** @jsx jsx */
import { jsx } from '@emotion/core'
import styled from '@emotion/styled'
import { Input, Label, Select } from '@rebass/forms'
import fetch from 'isomorphic-fetch'
import getConfig from 'next/config'
import { useState, FormEvent } from 'react'
import { Box, Button, Flex } from 'rebass'

import useAppState from '../../hooks/useAppState'

const { publicRuntimeConfig } = getConfig()

const StyledButton = styled(Button)`
    padding: 10px 20px;
    cursor: pointer;
    font-weight: 700;
    text-transform: uppercase;

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
}

const initialConfig: Config = Object.freeze({
    shapeType: 'polygons',
    numShapes: 200,
    shapeSize: 20,
    popSize: 50,
    poolSize: 10,
    mutationRate: 0.02,
    crossRate: 0.2,
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
        img.src = 'https://picsum.photos/900'

        img.onload = () => {
            dispatch({ type: 'update', payload: { targetImage: getBase64Image(img) } })
            setLoading(false)
        }
    }

    const onSubmit = (e: FormEvent) => {
        e.preventDefault()
        dispatch({ type: 'clearOutput' })

        fetch(`${publicRuntimeConfig.apiUrl}/job`, {
            method: 'POST',
            headers: { 'Content-Type': 'application/json' },
            body: JSON.stringify({ ...config, targetImage: state.target }),
        })
    }

    const fieldProps = (name: string) => ({
        id: name,
        name: name,
        value: config[name],
        onChange: (e: InputEvent) => {
            console.log(Object.keys(e.target))
            setConfig({
                ...config,
                [e.target.name!]: e.target.value!,
            })
        }
    })

    return (
        <form onSubmit={onSubmit}>
            <Flex justifyContent='space-between'>
                <StyledButton disabled={loading} onClick={getRangomTargetImage}>
                    Random Target Image
                </StyledButton>
                <StyledButton color='secondary' type='submit'>
                    Start
                </StyledButton>
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
            </Flex>
        </form>
    )
}
