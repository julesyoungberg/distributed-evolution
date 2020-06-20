import { useTheme } from 'emotion-theming'

import { Theme } from '../theme'

function useMyTheme(): Theme {
    return useTheme<Theme>()
}

export default useMyTheme
