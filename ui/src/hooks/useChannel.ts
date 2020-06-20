import getConfig from 'next/config'
import { useEffect } from 'react'

import { Action } from '../state'

const { publicRuntimeConfig } = getConfig()

export default function useChannel(dispatch: (action: Action) => void) {
    useEffect(() => {
        const onOpen = () => console.log('subscribed to server')

        const onMessage = (event: { data: string }) => {
            console.log('message from server')

            let payload: Record<string, any>

            try {
                payload = JSON.parse(event.data)
            } catch (e) {
                payload = { error: e }
            }

            dispatch({ type: 'update', payload })
        }

        const onClose = (e: any) => {
            console.log('socket closed, reconnecting in 1 second', e.reason)
            setTimeout(listener, 1000)
        }

        let socket

        const listener = () => {
            socket = new WebSocket(publicRuntimeConfig.channelUrl)
            socket.addEventListener('open', onOpen)
            socket.addEventListener('message', onMessage)
            socket.addEventListener('close', onClose)
        }

        listener()
    }, [])
}
