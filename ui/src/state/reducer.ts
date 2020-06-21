import { Action, State } from '.'

function handleStatus(payload: Record<string, any>): Partial<State> {
    if (['active', 'idle'].includes(payload.state)) {
        return { status: payload.state }
    }

    return {
        error: payload.error || payload.statusCode,
        status: 'error',
    }
}

function handleUpdate(state: State, payload: Record<string, any>): Partial<State> {
    if (payload.error) {
        return { error: payload.error, status: state.status === 'editing' ? 'editing' : 'error' }
    }

    const nextState: Partial<State> = { status: state.status === 'editing' ? 'editing' : 'active' }
    
    const fields = ['generation', 'numWorkers', 'output', 'targetImage', 'tasks']

    fields.forEach(field => {
        if (payload[field]) nextState[field] = payload[field]
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
            return { ...state, nextTargetImage: action.payload.target }
        case 'start':
            return { 
                ...state, 
                nextTargetImage: undefined,
                output: undefined,
                targetImage: state.nextTargetImage,
                ...handleUpdate(state, action.payload),
                status: 'active',
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
