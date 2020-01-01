import React, { Component } from 'react'
import { view } from 'react-easy-state'
import styles from '../../styles/components/SiteView.less'

class SiteView extends Component {
    render(): JSX.Element {
        return <div className={styles.container}>
            SITEVIEW!
        </div>
    }
}

export default view(SiteView)
