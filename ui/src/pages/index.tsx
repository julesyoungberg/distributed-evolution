import Head from 'next/head'
import { useState } from 'react'

import useWebsocket from '../hooks/useWebsocket'

export default function Home() {
    const [targetImage, setTargetImage] = useState<string | undefined>(undefined)

    const onMessage = (data: Record<string, any>) => {
        if (data.targetImage) setTargetImage(data.targetImage)
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

                <div>
                    <h2>Target Image</h2>
                    {targetImage && <img src={`data:image/jpg;base64, ${targetImage}`} />}
                </div>
            </main>
        </div>
    )
}
