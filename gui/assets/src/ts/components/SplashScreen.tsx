import React, { Component } from 'react'
import { view } from 'react-easy-state'
import styles from '../../styles/components/SplashScreen.less'
import {AppVersion, CommitHash} from '../Config'

class SplashScreen extends Component {
    render() {
        return <div className={styles.headerContainer}>
            <div className={styles.logoText}>
                <label className={styles.textCred}>Cred</label>
                <label className={styles.textMan}>Man</label>
            </div>

            <div className={styles.motto}>
                <label>Safeguard your credentials!</label>
            </div>

            <div className={styles.info}>
                <label className={styles.infoText}>version</label>
                <label className={styles.infoValue}>{AppVersion}</label>
                <label className={styles.infoText}></label>
                <label className={styles.infoText}>commit</label>
                <label className={styles.infoValue}>{CommitHash}</label>
            </div>
        </div>
    }
}

export default view(SplashScreen)
