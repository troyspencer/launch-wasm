import React from 'react'
import FlexView from 'react-flexview'
import { FormCheckbox } from "shards-react";

export default class StatsToggle extends React.Component {
    constructor(props) {
        super(props)
        this.handleChange = this.handleChange.bind(this)
    }

    handleChange() {
        this.props.onShowStatsChange(!this.props.showStats)
    }

    render(){
        return (
            <FlexView vAlignContent="center" hAlignContent='center'>
                <FlexView marginLeft="10" marginRight="10" vAlignContent="center" hAlignContent='center'>
                <FormCheckbox
                    checked={this.props.showStats}
                    onChange={this.handleChange}
                    >
                        {this.props.showStats ? "Stats Visible" : "Stats Hidden"}
                </FormCheckbox>
                </FlexView>
            </FlexView>  
        )
    }
    
}


