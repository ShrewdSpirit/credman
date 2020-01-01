import { invoke } from '..'

export interface TestResult {
    result: number
}

export async function test(): Promise<TestResult> {
    return await invoke<TestResult>('test', [12])
}
