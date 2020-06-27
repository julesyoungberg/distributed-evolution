/** @jsx jsx */
import { jsx } from '@emotion/core'
import styled from '@emotion/styled'
import { ThemeProvider } from 'emotion-theming'
import Head from 'next/head'
import { useReducer } from 'react'
import { Heading } from 'rebass'

import Control from '../components/Control'
import Metrics from '../components/Metrics'
import Status from '../components/Status'

import useChannel from '../hooks/useChannel'
import { initialState, StateContext } from '../state'
import reducer from '../state/reducer'
import theme from '../theme'

const Main = styled.main`
    font-family: 'system-ui', sans-serif;

    max-width: 1200px;
    margin: auto;
`

const StyledHeading = styled(Heading)`
    margin-bottom: 20px;
`

const B = styled.span`
    color: black;
`

export default function Home() {
    const [state, dispatch] = useReducer(reducer, initialState)

    useChannel(dispatch)

    return (
        <>
            <Head>
                <title>Distributed Evolution</title>
                <link rel='icon' href='/favicon.ico' />
            </Head>

            <ThemeProvider theme={theme}>
                <StateContext.Provider value={{ dispatch, state }}>
                    <Main>
                        <StyledHeading color='primary' fontSize={[5, 6, 7]} letterSpacing='-2px'>
                            D<B>i</B>st<B>r</B>ibu<B>t</B>ed Ev<B>o</B>lut<B>i</B>on
                        </StyledHeading>

                        <Status />

                        <Control />

                        <Metrics />
                    </Main>
                </StateContext.Provider>
            </ThemeProvider>
        </>
    )
}
