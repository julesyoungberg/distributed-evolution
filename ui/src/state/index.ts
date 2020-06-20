import React from 'react'

export type Status = 'active' | 'error' | 'idle'

export interface State {
    error?: string
    generation: number
    output?: string
    status?: Status
    target?: string
}

export const initialState: State = {
    generation: 0,
}

export type ActionType = 
    | 'clearOutput'
    | 'clearTarget' 
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
