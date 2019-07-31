import React, { Component } from 'react'
import { view, store } from 'react-easy-state'
import styles from '../../styles/components/splashscreen.less'
import FlexLayout, { Direction, JustifyContent, AlignItems } from './flexlayout'
import { getInfo } from '../interface'

const infoState = store({
    version: '',
    commitHash: '',
    motto: '',
})

class SplashScreen extends Component {
    async componentDidMount() {
        try {
            const result = await getInfo()
            infoState.version = result.version
            infoState.commitHash = result.commithash
            infoState.motto = result.motto
        } catch (e) {
            console.log(e)
        }
    }

    render() {
        return <FlexLayout direction={Direction.Column} className={styles.container} justifyContent={JustifyContent.Center} >
            <FlexLayout direction={Direction.Column} className={styles.headerContainer} justifyContent={JustifyContent.Center} alignItems={AlignItems.Center}>
                <FlexLayout justifyContent={JustifyContent.Center} alignItems={AlignItems.Center}>
                    <label className={styles.textCred}>Cred</label>
                    <label className={styles.textMan}>Man</label>
                </FlexLayout>

                <div className={styles.motto}>
                    <label>{infoState.motto}</label>
                </div>
            </FlexLayout>

            <FlexLayout className={styles.info} justifyContent={JustifyContent.Center} alignItems={AlignItems.Center}>
                <label className={styles.infoText}>version</label>
                <label className={styles.infoValue}>{infoState.version}</label>
                <label className={styles.infoText}></label>
                <label className={styles.infoText}>commit</label>
                <label className={styles.infoValue}>{infoState.commitHash}</label>
            </FlexLayout>
        </FlexLayout>
    }
}

export default view(SplashScreen)
