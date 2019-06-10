import React from 'react'
import Launches from './launches'

export default class Stats extends React.Component {
    constructor(props) {
        super(props)
    }
    render(){
        return (
            <div hidden={!this.props.showStats}>
                <Launches />
            </div>
        )
    }
    
}


