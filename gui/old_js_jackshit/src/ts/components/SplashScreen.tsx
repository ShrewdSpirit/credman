import React, { Component } from 'react'
import { view, store } from 'react-easy-state'
import styles from '../../styles/components/SplashScreen.less'
import FlexLayout, { Direction, JustifyContent, AlignItems } from './FlexLayout'
import { getInfo } from '../interface'
import { conditionalStr } from '../Utils'

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

                <div className={styles.mottoContainer}>
                    <label className={`${styles.mottoText} ${conditionalStr(infoState.motto, styles.mottoTextLoaded)}`}>
                        {conditionalStr(infoState.motto, infoState.motto, 'N/A')}
                    </label>
                </div>
            </FlexLayout>

            <FlexLayout className={styles.info} justifyContent={JustifyContent.Center} alignItems={AlignItems.Center}>
                <label className={styles.infoText}>version</label>
                <label className={`${styles.infoValue} ${conditionalStr(infoState.version, styles.infoValueLoaded)}`}>{infoState.version}</label>
                <label className={styles.infoText}></label>
                <label className={styles.infoText}>commit</label>
                <label className={`${styles.infoValue} ${conditionalStr(infoState.commitHash, styles.infoValueLoaded)}`}>{infoState.commitHash}</label>
            </FlexLayout>
        </FlexLayout>
    }
}

export default view(SplashScreen)
