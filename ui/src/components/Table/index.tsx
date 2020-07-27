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
        name: 'Job ID',
        selector: 'jobID',
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
]

export default function Table() {
    const { state } = useAppState()

    const now = new Date().getTime()

    return (
        <DataTable
            title='Workers'
            columns={columns}
            customStyles={{
                header: {
                    style: {
                        paddingLeft: 0,
                    }
                },
            }}
            data={Object.values(state.tasks || {}).map((task) => ({
                ...task,
                lastUpdate: `${(now - new Date(task.lastUpdate).getTime()) / 1000} seconds ago`,
            }))}
        />
    )
}
