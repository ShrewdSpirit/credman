import React, { Component } from 'react'
import { view, store } from 'react-easy-state'

import style from '../css/App.css'

class App extends Component {
    render() {
        return <div>
            <div className={style.headerContainer}>
                <div className={style.logoText}>
                    <label className={style.textCred}>Cred</label><label className={style.textMan}>Man</label>
                </div>
                <div className={style.info}>
                    <label className={style.infoText}>version</label>
                    <label className={style.infoValue}>{window.AppVersion}</label>
                    <label className={style.infoText}></label>
                    <label className={style.infoText}>commit</label>
                    <label className={style.infoValue}>{window.CommitHash}</label>
                </div>
            </div>
        </div>
    }
}

export default view(App)
