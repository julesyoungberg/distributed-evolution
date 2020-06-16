/** @jsx jsx */
import { jsx } from '@emotion/core'
import styled from '@emotion/styled'
import Head from 'next/head'
import { useState } from 'react'

import useWebsocket from '../hooks/useWebsocket'

const Wrapper = styled.div`
    display: flex;

    img {
        max-width: 100%;
    }
`

export default function Home() {
    const [targetImage, setTargetImage] = useState<string | undefined>(undefined)
    const [output, setOutput] = useState<string | undefined>(undefined)
    const [generation, setGeneration] = useState<number>(0)

    const onMessage = (data: Record<string, any>) => {
        if (data.targetImage) setTargetImage(data.targetImage)
        if (data.currentGeneration) setGeneration(data.currentGeneration)
        if (data.output) setOutput(data.output)
    }

    useWebsocket(onMessage)

    return (
        <div className='container'>
            <Head>
                <title>Distributed Evolution</title>
                <link rel='icon' href='/favicon.ico' />
            </Head>

            <main>
                <h1>Distributed Evolution</h1>

                <Wrapper>
                    <div>
                        <h2>Target Image</h2>
                        {targetImage && <img src={`data:image/jpg;base64, ${targetImage}`} />}
                    </div>

                    <div>
                        <h2>Output - Generation: {generation}</h2>
                        {output && <img src={`data:image/png;base64, ${output}`} />}
                    </div>
                </Wrapper>
            </main>
        </div>
    )
}
