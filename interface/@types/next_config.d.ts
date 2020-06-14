declare interface PublicRuntimeConfig {
    apiUrl: string
    environment: 'development' | 'production'
}


declare interface NextConfig {
    publicRuntimeConfig: PublicRuntimeConfig
}

declare module 'next/config' {
    export default function getConfig(): NextConfig
}
