import getConfig from 'next/config'
import { useEffect } from 'react'

const { publicRuntimeConfig: { apiUrl } } = getConfig()

export default function useWebsocket(handleMessage: (data: Record<string, any>) => void) {
    useEffect(() => {
        const onOpen = () => console.log('subscribed to server')

        const onMessage = (event: { data: string }) => {
            console.log('message from server')
        
            let json: Record<string, any> | undefined

            try {
                json = JSON.parse(event.data)
            } catch (e) {
                /* silence is golden */
            }

            if (json) handleMessage(json)
        }

        const onClose = (e: any) => {
            console.log('socket closed, reconnecting in 1 second', e.reason)
            setTimeout(listener, 1000)
        }

        let socket

        const listener = () => {
            socket = new WebSocket(`ws://${apiUrl}/subscribe`)
            socket.addEventListener('open', onOpen)
            socket.addEventListener('message', onMessage)
            socket.addEventListener('close', onClose)
        }

        listener()
    }, [])
}
