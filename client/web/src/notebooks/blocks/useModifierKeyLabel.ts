import { useMemo } from 'react'

import { isMacPlatform as isMacPlatformFunction } from '@sourcegraph/common'

export const useModifierKeyLabel = (): string => {
    const isMacPlatform = useMemo(() => isMacPlatformFunction(), [])
    return useMemo(() => (isMacPlatform ? '⌘' : 'Ctrl'), [isMacPlatform])
}
