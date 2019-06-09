import React from 'react'
import { Segment } from 'semantic-ui-react'
import Launches from './launches'

export default class Stats extends React.Component {
    constructor(props) {
        super(props)
    }
    render(){
        return (
            <Segment hidden={!this.props.showStats} floated="right" vertical={true}>
                <Launches />
            </Segment>
        )
    }
    
}


