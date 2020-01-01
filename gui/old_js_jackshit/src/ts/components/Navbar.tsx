import React, { Component } from 'react'
import { view, store } from 'react-easy-state'
import styles from '../../styles/components/Navbar.less'
import FlexLayout, { JustifyContent, Direction } from './FlexLayout';

class NavBar extends Component {
    render(): JSX.Element {
        return <FlexLayout className={styles.container}>
            <FlexLayout direction={Direction.Column} justifyContent={JustifyContent.Center}>
                <div className={styles.logoContainer}>
                    <label className={styles.logoTextCred}>Cred</label>
                    <label className={styles.logoTextMan}>Man</label>
                </div>
            </FlexLayout>

            <FlexLayout className={styles.itemsContainer} justifyContent={JustifyContent.SpaceBetween}>
                <FlexLayout>
                    PATH
                </FlexLayout>

                <FlexLayout>
                    Settings and shit
                </FlexLayout>
            </FlexLayout>
        </FlexLayout>
    }
}

export default view(NavBar)
