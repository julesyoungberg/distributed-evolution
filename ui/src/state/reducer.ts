import { Action, State, Status } from '.'

function handleStatus(payload: Record<string, any>): Partial<State> {
    if (payload.state === 'active') {
        return { status: payload.state }
    }

    return {
        error: payload.error || payload.statusCode,
        status: 'error',
    }
}

function handleUpdate(state: State, payload: Record<string, any>): Partial<State> {
    if (payload.generation < state.generation) {
        return state
    }

    let status: Status = 'active'

    if (state.status === 'editing') {
        status = 'editing'
    } else if (payload.jobID === 0) {
        status = 'idle'
    }

    if (payload.error) {
        return { error: payload.error, status }
    }

    const nextState: Partial<State> = { status }
    const fields = ['fitness', 'generation', 'jobID', 'numWorkers', 'output', 'startedAt', 'targetImage', 'tasks']

    fields.forEach((field) => {
        if (payload[field] || payload[field] === 0) nextState[field] = payload[field]
    })

    return nextState
}

function reducer(state: State, action: Action): State {
    switch (action.type) {
        case 'status':
            return { ...state, ...handleStatus(action.payload) }
        case 'update':
            return { ...state, ...handleUpdate(state, action.payload) }
        case 'clearTarget':
            return { ...state, nextTargetImage: undefined, status: 'editing' }
        case 'setTarget':
            return { ...state, nextTargetImage: action.payload.target, status: 'editing' }
        case 'start':
            return {
                ...state,
                nextTargetImage: undefined,
                output: undefined,
                targetImage: state.nextTargetImage,
                ...handleUpdate(state, action.payload),
                status: 'active',
                error: action.payload.statusCode >= 400 ? action.payload.data : '',
            }
        case 'start':
        default:
            return state
    }
}

export default function(state: State, action: Action): State {
    const nextState = reducer(state, action)
    // if (action.type !== 'update') {
    //     console.log('PREV STATE', state)
    //     console.log('NEXT STATE', nextState)
    // }
    return nextState
}
