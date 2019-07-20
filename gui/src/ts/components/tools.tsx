import React, { Component } from 'react'
import { view } from 'react-easy-state'
import styles from '../../styles/components/tools.less'

class Tools extends Component {
    render(): JSX.Element {
        return <div className={styles.container}>
            TOOLS
        </div>
    }
}

export default view(Tools)
