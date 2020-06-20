/** @jsx jsx */
import { jsx } from '@emotion/core'
import { Box, Flex, Text } from 'rebass'

import useAppState from '../../hooks/useAppState'

import Image from './Image'

export default function Status() {
    const { state } = useAppState()
    const { generation, output, target } = state

    return (
        <Flex css={{ paddingBottom: 20 }}>
            <Box width={1 / 2}>
                <Text fontSize={[3, 4, 5]} fontWeight='bold'>
                    Target
                </Text>
                <Image src={target ? `data:image/jpg;base64, ${target}` : undefined} />
            </Box>

            <Box width={1 / 2}>
                <Text fontSize={[3, 4, 5]} fontWeight='bold'>
                    Output - Generation: {generation}
                </Text>
                <Image src={output ? `data:image/png;base64, ${output}` : undefined} />
            </Box>
        </Flex>
    )
}
