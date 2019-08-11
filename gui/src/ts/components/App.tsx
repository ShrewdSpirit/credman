import React, { Component } from 'react'
import { view, store } from 'react-easy-state'
import styles from '../../styles/components/App.less'
import { encryptionHandshake, test } from '../interface'
import SplashScreen from './SplashScreen'
import MainScreen from './MainScreen'

const splashScreenState = store({ visible: true })

class App extends Component {
    async componentDidMount() {
        try {
            await encryptionHandshake()
            console.log('Encryption enabled')

            setTimeout(() => {
                splashScreenState.visible = false
            }, 1000)
        } catch (e) {
            console.log('Failed to handshake:', e)
        }

        try {
            console.log('TEST', await test())
        } catch (e) {
            console.log('faillllll', e)
        }
    }

    render(): JSX.Element {
        return <div className={styles.container}>
            {splashScreenState.visible ? <SplashScreen /> : <MainScreen />}
        </div>
    }
}

export default view(App)
