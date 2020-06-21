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

    return (
        <DataTable
            title="Workers"
            columns={columns}
            data={state.tasks}
        />
    )
}