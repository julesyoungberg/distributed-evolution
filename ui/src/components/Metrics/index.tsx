/** @jsx jsx */
import { jsx } from '@emotion/core'
import DataTable from 'react-data-table-component'

import useAppState from '../../hooks/useAppState'

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
    }
]

export default function Metrics() {
    const { state } = useAppState()

    const now = new Date().getTime()

    return (
        <DataTable
            title="Workers"
            columns={columns}
            data={Object.values(state.tasks || {}).map(task => ({
                ...task,
                lastUpdate: `${(now - (new Date(task.lastUpdate)).getTime()) / 1000} seconds ago`
            }))}
        />
    )
}