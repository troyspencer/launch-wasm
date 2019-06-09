import React from 'react'
import { Label } from 'semantic-ui-react'

export default class Launches extends React.Component {
    constructor(props) {
        super(props)
        this.state = {
            launches: 0,
        }
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
        this.setState({
            launches: event.launches,
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

