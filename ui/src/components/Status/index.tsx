/** @jsx jsx */
import { jsx } from '@emotion/core'
import { Box, Flex, Text } from 'rebass'

import useAppState from '../../hooks/useAppState'

import Image from './Image'

export default function Status() {
    const { state } = useAppState()
    const { error, fitness, generation, jobID, nextTargetImage, output, status, targetImage } = state

    if (['active', 'editing'].includes(status) && jobID > 0) {
        let targetSrc = status === 'editing' ? nextTargetImage : targetImage
   
        return (
            <Flex css={{ paddingBottom: 20 }}>
                <Box width={1 / 2}>
                    <Text css={{ marginTop: '38px' }} fontSize={[3, 4, 5]} fontWeight='bold'>
                        Target Image
                    </Text>
                    <Image src={targetSrc ? `data:image/jpg;base64, ${targetSrc}` : undefined} />
                </Box>

                {['active', 'editing'].includes(status) && (
                    <Box width={1 / 2}>
                        <Text fontSize={[3, 4, 5]} fontWeight='bold'>
                            Generation: {generation}<br/>Fitness: {fitness}
                        </Text>
                        <Image src={output ? `data:image/png;base64, ${output}` : undefined} />
                    </Box>
                )}
            </Flex>
        )
    }

    let msg = ''

    if (jobID === 0) {
        msg = 'No active job'
    } else if (status === 'error' && error) {
        msg = 'Error: ' + error
    } else {
        msg = 'Disconnected from cluster'
    }

    console.log(state)

    return (
        <Box
            css={{
                paddingBottom: 20,
                width: '100%',
                height: 0,
                marginTop: '76px',
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
