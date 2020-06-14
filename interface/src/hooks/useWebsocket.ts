import getConfig from 'next/config'
import { useEffect } from 'react'

const { publicRuntimeConfig: { apiUrl } } = getConfig()

export default function useWebsocket(handleMessage: (data: Record<string, any>) => void) {
    useEffect(() => {
        const socket = new WebSocket(`ws://${apiUrl}/subscribe`)

        const onOpen = () => console.log('subscribed to server')

        const onMessage = (event: { data: string }) => {
            console.log('message from server: ', event.data)
        
            let json: Record<string, any> | undefined

            try {
                json = JSON.parse(event.data)
            } catch (e) {
                /* silence is golden */
            }

            if (json) handleMessage(json)
        }

        socket.addEventListener('open', onOpen)
        socket.addEventListener('message', onMessage)

        return () => {
            socket.removeEventListener('open', onOpen)
            socket.removeEventListener('message', onMessage)
        }
    })
}
