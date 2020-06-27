/** @jsx jsx */
import { jsx } from '@emotion/core'
import styled from '@emotion/styled'
import fetch from 'isomorphic-fetch'
import getConfig from 'next/config'
import DataTable from 'react-data-table-component'
import { Button } from 'rebass'

import useAppState from '../../hooks/useAppState'

import { Theme } from '../../theme'

const { publicRuntimeConfig } = getConfig()

const columns = [
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
    },
    {
        name: 'Connection',
        selector: 'connection',
    },
]

const ConnectionButton = styled(Button)<{ disabled: boolean }, Theme>`
    background-color: ${({ disabled, theme }) => disabled ? theme.colors.lightgray : theme.colors.blue};
    cursor: pointer;
`

export default function Table() {
    const { state } = useAppState()

    const now = new Date().getTime()

    const disconnect = (task) => async () => {
        const path = `tasks/${task.ID}/disconnect`

        console.log(`disconnecting task ${task.ID}`)
        
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
                    <ConnectionButton
                        disabled={!(task.connected && task.status == 'inprogress')}
                        onClick={disconnect(task)}
                    >
                        DISCONNECT
                    </ConnectionButton>
                )
            }))}
        />
    )
}