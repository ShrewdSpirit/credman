import React, { Component } from 'react'
import { view } from 'react-easy-state'
import styles from '../../styles/components/Profiles.less'
import FlexLayout, { Direction, AlignItems, JustifyContent } from './FlexLayout';
import { FontAwesomeIcon } from '@fortawesome/react-fontawesome'
import { faPlus, faDownload } from '@fortawesome/free-solid-svg-icons'
import ProfileItem from './ProfileItem';

class Profiles extends Component {
    render(): JSX.Element {
        return <FlexLayout className={styles.container} direction={Direction.Column} alignItems={AlignItems.Stretch}>
            <FlexLayout direction={Direction.Row} className={styles.managePanel}>
                <FlexLayout className={styles.managePanelItems} justifyContent={JustifyContent.Center} alignItems={AlignItems.Center}>
                    <FontAwesomeIcon icon={faPlus} />
                    Add
                </FlexLayout>

                <FlexLayout className={styles.managePanelItems} justifyContent={JustifyContent.Center} alignItems={AlignItems.Center}>
                    <FontAwesomeIcon icon={faDownload} />
                    Import
                </FlexLayout>
            </FlexLayout>

            <FlexLayout className={styles.profilesContainer} direction={Direction.Column} alignItems={AlignItems.Stretch}>
                <ProfileItem></ProfileItem>
                <ProfileItem></ProfileItem>
                <ProfileItem></ProfileItem>
                <ProfileItem></ProfileItem>
                <ProfileItem></ProfileItem>
                <ProfileItem></ProfileItem>
                <ProfileItem></ProfileItem>
                <ProfileItem></ProfileItem>
                <ProfileItem></ProfileItem>
                <ProfileItem></ProfileItem>
                <ProfileItem></ProfileItem>
                <ProfileItem></ProfileItem>
                <ProfileItem></ProfileItem>
                <ProfileItem></ProfileItem>
                <ProfileItem></ProfileItem>
                <ProfileItem></ProfileItem>
                <ProfileItem></ProfileItem>
                <ProfileItem></ProfileItem>
                <ProfileItem></ProfileItem>
                <ProfileItem></ProfileItem>
                <ProfileItem></ProfileItem>
                <ProfileItem></ProfileItem>
                <ProfileItem></ProfileItem>
                <ProfileItem></ProfileItem>
                <ProfileItem></ProfileItem>
                <ProfileItem></ProfileItem>
                <ProfileItem></ProfileItem>
                <ProfileItem></ProfileItem>
            </FlexLayout>
        </FlexLayout>
    }
}

export default view(Profiles)
