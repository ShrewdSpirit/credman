import ezsec from 'ezsec'

declare const development: boolean // global defined by webpack
const clientId = ezsec.uuid().value
const clientKey = ezsec.randomBuffer(32)

function buildFetchOptions(body: string, encrypt: boolean): RequestInit {
    return {
        method: 'POST',
        mode: 'cors',
        cache: 'no-cache',
        credentials: 'same-origin',
        headers: {
            'Content-Type': 'application/json',
            'Encrypted': encrypt ? 'true' : 'false',
            'Client-Id': clientId
        },
        redirect: 'follow',
        referrer: 'no-referrer',
        body
    }
}

function urlBuilder(path: string): string {
    return development ? `http://localhost:14201/${path}` : `/${path}`
}

interface ServerPublicKey {
    publickey: string
}

export async function encryptionHandshake(): Promise<any> {
    return new Promise<any>(async (resolve, reject) => {
        try {
            const serverPublic = await invoke<ServerPublicKey>('handshake_getkey', [], false)
            const encrypted = ezsec.RSAEncrypt(ezsec.ShaType.sha512, clientKey, serverPublic.publickey)
            if (encrypted.error) {
                reject('Failed to encrypt client key: ' + encrypted.error)
                return
            }
            await invoke('handshake_setkey', [clientId, encrypted.value], false)
            resolve()
        } catch (e) {
            reject(e)
        }
    })
}

export function invoke<R>(method: string, args: any[] = [], encrypt: boolean = true): Promise<R> {
    return new Promise<R>((resolve, reject) => {
        for (let i = 0; i < args.length; i++) {
            if (args[i] instanceof Uint8Array) {
                args[i] = Array.from(args[i])
            }
        }

        let data = JSON.stringify({ method, args })

        if (encrypt) {
            const encryptedData = ezsec.CFBEncrypt(ezsec.ShaType.sha512, data, clientKey)
            if (encryptedData.error) {
                reject(encryptedData.error)
                return
            }
            data = ezsec.uint8ArrayToString(encryptedData.value)
        }

        fetch(urlBuilder('invoke'), buildFetchOptions(data, encrypt))
            .then(async resp => {
                let result: string

                if (encrypt) {
                    const blob = await resp.blob()
                    const data = new Uint8Array(await new Response(blob).arrayBuffer())
                    const decryptedResult = ezsec.CFBDecrypt(ezsec.ShaType.sha512, data, clientKey)
                    if (decryptedResult.error) {
                        reject(decryptedResult.error)
                        return
                    }
                    result = ezsec.uint8ArrayToString(decryptedResult.value)
                } else {
                    result = await resp.text()
                }

                try {
                    const jsonResult = JSON.parse(result)
                    resolve(jsonResult as R)
                } catch (e) {
                    reject(result)
                }
            })
            .catch(err => reject(err))
    })
}
