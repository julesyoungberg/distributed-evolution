import styled from '@emotion/styled'
import { Button } from 'rebass'

import { Theme } from '../theme'

const StyledButton = styled(Button)<{ disabled: boolean }, Theme>`
    padding: 10px 20px;
    cursor: pointer;
    font-weight: 700;
    text-transform: uppercase;
    background-color: ${({ disabled, theme }) => (disabled ? theme.colors.lightgray : theme.colors.primary)};

    &:focus {
        outline: none;
    }
`

export default StyledButton
