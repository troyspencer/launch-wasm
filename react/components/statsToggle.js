import React from 'react'
import { Checkbox, Segment, Label, Divider, Header } from 'semantic-ui-react'

export default class StatsToggle extends React.Component {
    constructor(props) {
        super(props)
        this.handleChange = this.handleChange.bind(this)
    }

    handleChange(event, data) {
        this.props.onShowStatsChange(data.checked)
    }

    render(){
        return (
            <Segment>
                <Header as='h2' floated='right'>
                Settings
                </Header>
                <Divider clearing />
                <Checkbox toggle onChange={this.handleChange} />      
                <Label color="grey" content="Show Stats" /> 
            </Segment>
        )
    }
    
}


