import getConfig from 'next/config'
import Head from 'next/head'
import { useEffect, useState } from 'react'

export default function Home() {
    const [targetImage, setTargetImage] = useState<string | undefined>(undefined)

    useEffect(() => {
        const { publicRuntimeConfig } = getConfig()
    
        const socket = new WebSocket(`ws://${publicRuntimeConfig.apiUrl}/subscribe`)

        const onOpen = () => {
            console.log('subscribed to server')
        }

        const onMessage = (event) => {
            console.log('message from server ', event.data)
            
            let json: Record<string, any> | undefined

            try {
                json = JSON.parse(event.data)
            } catch (e) {
                // silence is golden
            }

            if (!json) return

            setTargetImage(json.targetImage)
        }

        socket.addEventListener('open', onOpen)
        socket.addEventListener('message', onMessage)

        return () => {
            socket.removeEventListener('open', onOpen)
            socket.removeEventListener('message', onMessage)
        }
    })

    return (
        <div className='container'>
            <Head>
                <title>Distributed Evolution</title>
                <link rel='icon' href='/favicon.ico' />
            </Head>

            <main>
                <h1>Welcome</h1>

                {targetImage && <img src={`data:image/jpg;base64, ${targetImage}`} />}
            </main>
        </div>
    )
}
