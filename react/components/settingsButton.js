import React from 'react'
import {Button, Tooltip} from 'antd'

export default class SettingsButton extends React.Component {
    constructor(props) {
        super(props)
    }
    render(){
        return (
            <Tooltip placement="right" title="Settings">
                <Button size="large" type="default" shape="circle" icon="setting" onClick={this.props.onClick}/>
            </Tooltip>
        )
    }
    
}



