import React, { Component } from 'react'
import { view, store } from 'react-easy-state'
import styles from '../../styles/components/Navbar.less'

export interface NavBarProps {
    className?: string
}

class NavBar extends Component<NavBarProps, {}> {
    render(): JSX.Element {
        return <div className={styles.container}>
            hey fuck
        </div>
    }
}

export default view(NavBar)
