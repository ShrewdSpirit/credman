import React, { Component } from 'react'
import { view } from 'react-easy-state'

export enum Direction {
    Row = 'row',
    RowReverse = 'row-reverse',
    Column = 'column',
    ColumnReverse = 'column-reverse'
}

export enum Wrap {
    NoWrap = 'nowrap',
    Wrap = 'wrap',
    WrapReverse = 'wrap-reverse'
}

export enum JustifyContent {
    FlexStart = 'flex-start',
    FlexEnd = 'flex-end',
    Center = 'center',
    SpaceBetween = 'space-between',
    SpaceAround = 'space-around',
    SpaceEvenly = 'space-evenly'
}

export enum AlignItems {
    FlexStart = 'flex-start',
    FlexEnd = 'flex-end',
    Center = 'center',
    Stretch = 'stretch',
    Baseline = 'baseline'
}

export enum AlignContent {
    FlexStart = 'flex-start',
    FlexEnd = 'flex-end',
    Center = 'center',
    SpaceBetween = 'space-between',
    SpaceAround = 'space-around',
    Stretch = 'stretch'
}

export interface FlexLayoutProps {
    className?: string
    direction: Direction
    wrap: Wrap
    justifyContent: JustifyContent
    alignItems: AlignItems
    alignContent: AlignContent
}

class FlexLayout extends Component<FlexLayoutProps, {}>{
    static defaultProps = {
        direction: Direction.Row,
        wrap: Wrap.NoWrap,
        justifyContent: JustifyContent.FlexStart,
        alignItems: AlignItems.Stretch,
        alignContent: AlignContent.Stretch,
    }

    render(): JSX.Element {
        const { children, className, direction, wrap, justifyContent, alignItems, alignContent } = this.props
        return (<div className={className}
            style={{
                display: 'flex',
                flexDirection: direction,
                flexWrap: wrap,
                justifyContent,
                alignItems,
                alignContent,
            }}>
            {children}
        </div>)
    }
}

export default view(FlexLayout)
