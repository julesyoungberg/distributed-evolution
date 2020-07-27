/** @jsx jsx */
import { jsx } from '@emotion/core'
import { Box, Flex, Text } from 'rebass'

import useAppState from '../../hooks/useAppState'

import Image from './Image'

export default function Status() {
    const { state } = useAppState()
    const { error, generation, jobID, nextTargetImage, output, status, targetImage } = state

    console.log(state)

    if (['active', 'editing'].includes(status)) {
        let targetSrc = status === 'editing' ? nextTargetImage : targetImage

        return (
            <Flex css={{ paddingBottom: 20 }}>
                <Box width={1 / 2}>
                    <Text fontSize={[3, 4, 5]} fontWeight='bold'>
                        Target
                    </Text>
                    <Image src={targetSrc} />
                </Box>

                {['active', 'editing'].includes(status) && jobID > 0 ? (
                    <Box width={1 / 2}>
                        <Text fontSize={[3, 4, 5]} fontWeight='bold'>
                            Output - Generation: {generation}
                        </Text>
                        <Image src={output} />
                    </Box>
                ) : (typeof error === 'string' && (
                    <Box width={1 / 2}>
                        <Flex
                            css={{ height: '100%' }}
                            alignItems='center'
                            justifyContent='center'
                        >
                            <Box>{error}</Box>
                        </Flex>
                    </Box>
                ))}
            </Flex>
        )
    }

    let msg = ''

    if (status === 'idle') {
        msg = 'No active job'
    } else if (status === 'error' && error) {
        msg = 'Error: ' + error
    } else {
        msg = 'Disconnected from cluster'
    }

    return (
        <Box
            css={{
                paddingBottom: 20,
                width: '100%',
                height: 0,
                marginTop: '60px',
                paddingTop: '50%',
                position: 'relative',
            }}
        >
            <Flex
                css={{
                    position: 'absolute',
                    top: 0,
                    bottom: 0,
                    left: 0,
                    right: 0,
                }}
                alignItems='center'
                justifyContent='center'
            >
                <Box>{msg}</Box>
            </Flex>
        </Box>
    )
}
