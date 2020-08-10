/** @jsx jsx */
import { jsx } from '@emotion/core'
import { Box, Flex, Text } from 'rebass'

import Image from '../Image'
import useAppState from '../../hooks/useAppState'

export default function AuxImages() {
    const { state } = useAppState()
    const { palette, targetImageEdges } = state

    if (!(palette || targetImageEdges)) {
        return null
    }

    return (
        <Flex css={{ paddingBottom: 20 }}>
            {targetImageEdges && (
                <Box width={1 / 2}>
                    <Text fontSize={[2, 3, 4]}>
                        Target Edges
                    </Text>
                    <Image src={targetImageEdges} />
                </Box>
            )}

            {palette && (
                <Box width={1 / 2}>
                    <Text fontSize={[2, 3, 4]}>
                        Palette
                    </Text>
                    <Image src={palette} />
                </Box>
            )}
        </Flex>
    )
}
