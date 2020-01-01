import React, { Component } from 'react'
import { view } from 'react-easy-state'
import styles from '../../styles/components/MainScreen.less'
import { currentScreen, Screen } from '../states'
import NavBar from './Navbar'
import Profiles from './Profiles'
import Sites from './Sites'
import SiteView from './SiteView'
import FlexLayout, { Direction, AlignItems, JustifyContent } from './FlexLayout';

class MainScreen extends Component {
    render(): JSX.Element {
        let screen: JSX.Element

        switch (currentScreen()) {
            case Screen.Profiles:
                screen = <Profiles></Profiles>
                break
            case Screen.Sites:
                screen = <Sites></Sites>
                break
            case Screen.SiteView:
                screen = <SiteView></SiteView>
                break
        }

        return <FlexLayout className={styles.container} direction={Direction.Column} alignItems={AlignItems.Stretch}>
            <NavBar></NavBar>
            <FlexLayout className={styles.screenContainer} justifyContent={JustifyContent.Center}>
                {screen}
            </FlexLayout>
        </FlexLayout>
    }
}

export default view(MainScreen)
