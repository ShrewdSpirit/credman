import React, { Component } from 'react'
import { view } from 'react-easy-state'
import styles from '../../styles/components/profiles.less'

class Profiles extends Component {
    render(): JSX.Element {
        return <div className={styles.container}>
            PROFILES!
        </div>
    }
}

export default view(Profiles)
