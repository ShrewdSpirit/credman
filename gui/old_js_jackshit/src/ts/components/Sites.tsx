import React, { Component } from 'react'
import { view } from 'react-easy-state'
import styles from '../../styles/components/Sites.less'

class Sites extends Component {
    render(): JSX.Element {
        return <div className={styles.container}>
            SITES!
        </div>
    }
}

export default view(Sites)
