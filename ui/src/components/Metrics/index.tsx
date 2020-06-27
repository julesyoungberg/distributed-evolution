/** @jsx jsx */
import { jsx } from '@emotion/core'
import styled from '@emotion/styled'
import fetch from 'isomorphic-fetch'
import getConfig from 'next/config'
import DataTable from 'react-data-table-component'
import { Button } from 'rebass'

import useAppState from '../../hooks/useAppState'

const { publicRuntimeConfig } = getConfig()

const columns = [
    {
        name: 'Connection',
        selector: 'connection',
    },
    {
        name: 'ID',
        selector: 'ID',
        sortable: true,
    },
    {
        name: 'Worker ID',
        selector: 'workerID',
        sortable: true,
    },
    {
        name: 'Thread',
        selector: 'thread',
    },
    {
        name: 'Generation',
        selector: 'generation',
        sortable: true,
    },
    {
        name: 'Last Update',
        selector: 'lastUpdate',
        sortable: true,
    },
    {
        name: 'Status',
        selector: 'status',
    }
]

const ConnectionButton = styled(Button)<{ connected: boolean }>`
    background-color: ${({ connected }) => connected ? 'red' : 'green'};
    cursor: pointer;
`

export default function Metrics() {
    const { state } = useAppState()

    const now = new Date().getTime()

    const toggleConnection = (task) => async () => {
        const action = task.connected ? 'disconnect' : 'reconnect'
        const path = `tasks/${task.ID}/${action}`

        console.log(`${action}ing task ${task.ID}`)
        
        const response = await fetch(`${publicRuntimeConfig.apiUrl}/${path}`)

        console.log(response, response)
    }

    return (
        <DataTable
            title="Workers"
            columns={columns}
            data={Object.values(state.tasks || {}).map(task => ({
                ...task,
                lastUpdate: `${(now - (new Date(task.lastUpdate)).getTime()) / 1000} seconds ago`,
                connection: (
                    <ConnectionButton connected={task.connected} onClick={toggleConnection(task)}>
                        {task.connected ? 'DISCONNECT' : 'RECONNECT'}
                    </ConnectionButton>
                )
            }))}
        />
    )
}