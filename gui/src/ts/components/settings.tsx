import React, { Component } from 'react'
import { view } from 'react-easy-state'
import styles from '../../styles/components/settings.less'

class Settings extends Component {
    render(): JSX.Element {
        return <div className={styles.container}>
            SETTINGS
        </div>
    }
}

export default view(Settings)
