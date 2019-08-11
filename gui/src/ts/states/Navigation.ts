import { store } from 'react-easy-state'

export enum Screen {
    Profiles,
    ProfileView,
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

type PathStack = Array<PathPart>

let currentPathStack = new PathStack()

export function getPathStack():

export const navigation = store({
    screen: Screen.Profiles,
})

export function navigate() {}

export function dialog() {}
