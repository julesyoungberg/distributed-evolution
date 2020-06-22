import getConfig from 'next/config'
import { useEffect } from 'react'

import { Action } from '../state'

const { publicRuntimeConfig } = getConfig()

// use a websocket channel to stay updated with the cluster's state
export default function useChannel(dispatch: (action: Action) => void) {
    useEffect(() => {
        let socket

        const onOpen = () => console.log('subscribed to server')

        const onMessage = (event: { data: string }) => {
            // ping pong
            if (event.data == 'keepalive') {
                socket?.send('keepalive')
                return 
            }

            let payload: Record<string, any>

            try {
                payload = JSON.parse(event.data)
            } catch (e) {
                payload = { error: e }
            }

            dispatch({ type: 'update', payload })
        }

        const onClose = (e: any) => {
            console.log('socket closed, reconnecting in 1 second. ', e.reason)
            // dispatch({ type: 'status', payload: { status: 'disconnected' } })
            setTimeout(listener, 1000)
        }

        const onError = (e: any) => {
            dispatch({ type: 'status', payload: { error: e.message } })
        }

        const listener = () => {
            socket = new WebSocket(publicRuntimeConfig.channelUrl)
            socket.addEventListener('open', onOpen)
            socket.addEventListener('message', onMessage)
            socket.addEventListener('close', onClose)
            socket.addEventListener('error', onError)
        }

        listener()
    }, [])
}
