/** @jsx jsx */
import { jsx } from '@emotion/core'
import styled from '@emotion/styled'
import { useLayoutEffect, useRef, useState } from 'react'
import BeatLoader from 'react-spinners/BeatLoader'
import { Box, Flex, Image as RebassImage } from 'rebass'

import useTheme from '../hooks/useTheme'
import { ReactElement } from 'react'

interface ImageProps {
    src?: string
}

const StyledImage = styled(RebassImage)`
    max-width: 100%;
    max-height: 100%;
`

const Container = styled.div`
    width: 100%;
    height: 0;
    padding-top: 100%;
    position: relative;
`

const Wrapper = styled.div`
    position: absolute;
    top: 0;
    bottom: 0;
    left: 0;
    right: 0;
`

const getImgSrc = (data: string) => `data:image/png;base64,${data}`

interface ImgDimensions {
    height: number
    width: number
}

function getScaledDimensions(image: ImgDimensions, container: ImgDimensions): ImgDimensions {
    let width: number = 0
    let height: number = 0

    if (image.width > image.height) {
        width = container.width
        height = image.height * (container.width / image.width)
    } else {
        height = container.height
        width = image.width * (container.height / image.height)
    }

    return { width, height }
}

function scaleImage(input: string, width: number, height: number): Promise<string> {
    return new Promise((resolve) => {
        const canvas = document.createElement('canvas')
        const ctx = canvas.getContext('2d')

        canvas.width = width
        canvas.height = height

        const image = new Image()
        image.onload = () => {
            const dimensions = getScaledDimensions(image, canvas)
            ctx.drawImage(
                image,
                0, 0, image.width, image.height,
                (canvas.width - dimensions.width) / 2,
                (canvas.height - dimensions.height) / 2, 
                dimensions.width, dimensions.height
            )

            const scaled = canvas.toDataURL()
            resolve(scaled)
        }

        image.src = getImgSrc(input)
    })
}

export default function Img({ src }: ImageProps) {
    const theme = useTheme()
    const containerRef = useRef<HTMLDivElement | undefined>(undefined)
    const [scaled, setScaled] = useState<string | undefined>(undefined)

    useLayoutEffect(() => {
        if (!containerRef.current) {
            return;
        }

        const { clientHeight, clientWidth } = containerRef.current

        async function getScaled() {
            const result = await scaleImage(src, clientWidth, clientHeight)
            setScaled(result)
        }
        
        getScaled()
    }, [containerRef, src])

    let content: ReactElement | undefined
    if (scaled) {
        content = <StyledImage src={scaled} />
    } else {
        content = <BeatLoader color={theme.colors?.primary} />
    }

    return (
        <Container ref={containerRef}>
            <Wrapper>
                <Flex css={{ height: '100%' }} alignItems='center' justifyContent='center'>
                    <Box>{content}</Box>
                </Flex>
            </Wrapper>
        </Container>
    )
}
