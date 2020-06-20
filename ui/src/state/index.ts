import React from 'react'

export type Status = 'active' | 'disconnected' | 'editing' | 'error' | 'idle'

export interface State {
    error?: any
    generation: number
    nextTargetImage?: string
    numWorkers?: number
    output?: string
    status?: Status
    targetImage?: string
    tasks?: Record<string, any>[]
}

export const initialState: State = {
    generation: 0,
}

export type ActionType = 
    | 'clearTarget' 
    | 'setTarget'
    | 'start'
    | 'status'
    | 'update'

export interface Action {
    type: ActionType
    payload?: Record<string, any>
}

export interface StateContextType {
    dispatch: (action: Action) => void
    state: State
}

export const StateContext = React.createContext<StateContextType>({
    dispatch: (_: Action) => null,
    state: initialState,
})
