import React from 'react'
import { Button } from 'semantic-ui-react'

export default class SettingsButton extends React.Component {
    render(){
        return (
            <Button onClick={this.props.onClick} circular color="grey" icon='settings' size='massive' />
        )
    }
    
}


