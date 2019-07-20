import React, { Component } from 'react'
import { view, store } from 'react-easy-state'
import { MemoryRouter as Router, Link, Route } from 'react-router-dom'
import styles from '../../styles/components/app.less'
import SplashScreen from './splashscreen'
import Profiles from './profiles'
import Tools from './tools'
import Settings from './settings'

const splashScreen = store({ visible: true })

class App extends Component {
    constructor(props: any) {
        super(props)
        setTimeout(() => {
            // splashScreen.visible = false
        }, 10000)
    }

    render(): JSX.Element {
        return <div className={styles.container}>
            {splashScreen.visible ?
                <SplashScreen /> :
                <Router>
                    <div>
                        <nav>
                            <Link to="/">Profiles</Link>
                            <Link to="/tools">Tools</Link>
                            <Link to="/settings">Settings</Link>
                        </nav>
                    </div>

                    <Route path="/" exact component={Profiles} />
                    <Route path="/tools" component={Tools} />
                    <Route path="/settings" component={Settings} />
                </Router>
            }
        </div>
    }
}

export default view(App)
