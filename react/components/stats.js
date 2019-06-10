import React from 'react'
import Launches from './launches'
import Timer from './timer'
import FlexView from 'react-flexview';

export default class Stats extends React.Component {
    constructor(props) {
        super(props)
    }
    render(){
        return (
            <div hidden={!this.props.showStats}>
                <FlexView column hAlignContent='left'>
                    <Launches />
                    <Timer />
                </FlexView>
            </div>
        )
    }
    
}


