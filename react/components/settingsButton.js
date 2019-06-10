import React from 'react'
import {Button, Tooltip} from 'antd'

const styles = {
    button: {
        marginBottom: '10px',
        marginTop: '10px',
        marginRight: '10px',
        marginLeft: '10px'
    }
}

export default class SettingsButton extends React.Component {
    constructor(props) {
        super(props)
    }
    render(){
        return (
            <Button style={styles.button} size="large" type="default" shape="circle" icon="setting" onClick={this.props.onClick}/>
        )
    }
    
}



