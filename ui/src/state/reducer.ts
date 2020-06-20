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

function handleUpdate(payload: Record<string, any>): Partial<State> {
    if (payload.error) {
        return { error: payload.error, status: 'error' }
    }

    const nextState: Partial<State> = {}

    // TODO make this a loop
    if (payload.targetImage) nextState.target = payload.targetImage
    if (payload.currentGeneration) nextState.generation = payload.currentGeneration
    if (payload.output) nextState.output = payload.output

    return nextState
}

export default function reducer(state: State, action: Action): State {
    switch (action.type) {
        case 'status':
            return { ...state, ...handleStatus(action.payload) }
        case 'update':
            return { ...state, ...handleUpdate(action.payload) }
        case 'clearTarget':
            return { ...state, target: undefined }
        case 'clearOutput':
            return { ...state, output: undefined }
        default:
            return state
    }
}
