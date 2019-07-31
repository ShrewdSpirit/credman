import ezsec from 'ezsec'

declare const development: boolean // global defined by webpack

export function urlBuilder(path: string): string {
    return development ? `http://localhost:14201/${path}` : `/${path}`
}

export async function handshake() {

}

export function invoke<R>(method: string, args: any[] = [], encrypt: boolean = true): Promise<R> {
    return new Promise<R>((resolve, reject) => {
        let data = JSON.stringify({ method, args })

        fetch(urlBuilder('invoke'), {
            method: 'POST',
            mode: 'cors',
            cache: 'no-cache',
            credentials: 'same-origin',
            headers: {
                'Content-Type': 'application/json',
                'Encrypted': encrypt ? 'true' : 'false'
            },
            redirect: 'follow',
            referrer: 'no-referrer',
            body: data,
        }).then(async resp => {
            const result = await resp.text()
            try {
                const jsonResult = JSON.parse(result)
                resolve(jsonResult as R)
            } catch (e) {
                reject(result)
            }
        }).catch(err => reject(err))
    })
}
