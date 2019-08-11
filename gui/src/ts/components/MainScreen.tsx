import React, { Component } from 'react'
import { view, store } from 'react-easy-state'
import styles from '../../styles/components/MainScreen.less'
import NavBar from './Navbar'

class MainScreen extends Component {
    render(): JSX.Element {
        return <div className={styles.container}>
            <NavBar></NavBar>
        </div>
    }
}

export default view(MainScreen)
