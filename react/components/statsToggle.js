import React from 'react'
import { Checkbox, Label } from 'semantic-ui-react'
import FlexView from 'react-flexview'

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
            <FlexView vAlignContent="center" hAlignContent='center'>
                <FlexView marginLeft="10" vAlignContent="center" hAlignContent='center'>
                    <Label color="grey" content={this.props.showStats ? "Stats Visible" : "Stats Hidden"} />   
                </FlexView>
                <FlexView marginLeft="10" marginRight="10" vAlignContent="center" hAlignContent='center'>
                    <Checkbox onChange={this.handleChange} />  
                </FlexView>
            </FlexView>  
        )
    }
    
}


