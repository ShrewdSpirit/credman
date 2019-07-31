import { invoke } from '..'

export interface GetInfoResult {
    version: string
    commithash: string
    motto: string
}

export async function getInfo(): Promise<GetInfoResult> {
    return await invoke<GetInfoResult>('getinfo', [], false)
}
