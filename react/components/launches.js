import React from 'react'
import { Label } from 'semantic-ui-react'

export default class Launches extends React.Component {
    constructor(props) {
        super(props)
        this.state = {
            launches: 0,
        }
        this.handleLaunchUpdate = this.handleLaunchUpdate.bind(this)
    }

    componentDidMount() {
        // When the component is mounted, add your DOM listener to the "nv" elem.
        // (The "nv" elem is assigned in the render function.)
        window.document.addEventListener("updateLaunches", this.handleLaunchUpdate);
      }
    
    componentWillUnmount() {
        // Make sure to remove the DOM listener when the component is unmounted.
        window.document.removeEventListener("updateLaunches", this.handleLaunchUpdate);
    }

    handleLaunchUpdate = (event) => {
        this.onLaunchUpdate(event.launches)
    }

    onLaunchUpdate(launches) {
        this.setState({
            launches: launches,
        })
    }

    render(){
        return (
            <Label color="grey">
                {this.state.launches}
            </Label>
        )
    }
    
}

