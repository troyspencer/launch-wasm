import React from 'react'
import {Button} from 'antd'

export default class SettingsButton extends React.Component {
    constructor(props) {
        super(props)
    }
    render(){
        return (
            <Button size="large" type="default" shape="circle" icon="setting" onClick={this.props.onClick}/>
        )
    }
    
}


