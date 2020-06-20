/** @jsx jsx */
import { jsx } from '@emotion/core'
import styled from '@emotion/styled'
import BeatLoader from 'react-spinners/BeatLoader'
import { Box, Flex, Image } from 'rebass'

import useTheme from '../../hooks/useTheme'
import { ReactElement } from 'react'

interface ImageProps {
    src?: string
}

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

export default function Img({ src }: ImageProps) {
    const theme = useTheme()

    let content: ReactElement | undefined

    if (src) {
        content = <Image src={src} />
    } else {
        content = (
            <Flex css={{ height: '100%' }} alignItems='center' justifyContent='center'>
                <Box>
                    <BeatLoader color={theme.colors?.primary} />
                </Box>
            </Flex>
        )
    }

    return (
        <Container>
            <Wrapper>{content}</Wrapper>
        </Container>
    )
}
