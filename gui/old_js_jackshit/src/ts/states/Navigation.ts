import { store } from 'react-easy-state'

export enum Screen {
    Profiles,
    Sites,
    SiteView,
}

export enum Dialog {
    Settings,
    NewSite,
    NewProfile,
    FileOpen,
    FileSave,
    YesNo,
}

export interface PathPart {
    screen: Screen
    name: string
}

const navigation = store({
    screen: Screen.Profiles,
})

let currentPathStack = new Array<PathPart>()

export const getPathStack = () => currentPathStack

export const currentScreen = () => navigation.screen

export function navigate() {}

export function back() {}

export function dialog() {}
