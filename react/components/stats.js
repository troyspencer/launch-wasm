import React from 'react'
import Launches from './launches'
import Timer from './timer'
import FlexView from 'react-flexview';

export default function Stats(props) {
    return (
        <div hidden={!props.showStats}>
            <FlexView vAlignContent='top'>
                <Launches />
                <Timer />
            </FlexView>
        </div>
    )
}
