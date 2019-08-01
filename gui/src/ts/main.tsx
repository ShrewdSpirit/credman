import React from 'react'
import ReactDOM from 'react-dom'
import '../styles/global.less'
import App from './components/app'
import { encryptionHandshake } from './interface'

async function doHandshake() {
    try {
        await encryptionHandshake()
        console.log('Encryption enabled')
    } catch (e) {
        console.log('Failed to handshake:', e)
    }
}
doHandshake()

ReactDOM.render(<App />, document.getElementById('root'))
