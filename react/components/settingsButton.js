import React from 'react'
import {Button} from 'antd'

const styles = {
    button: {
        marginBottom: '10px',
        marginTop: '10px',
        marginRight: '10px',
        marginLeft: '10px'
    }
}

export default function SettingsButton(props) {
    return (
        <Button style={styles.button} size="large" type="primary" shape="circle" icon="setting" onClick={props.onClick}/>
    )
}
