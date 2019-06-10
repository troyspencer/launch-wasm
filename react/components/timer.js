import React from 'react'
import {Button, Tooltip, Badge} from 'antd'

const styles = {
    button: {
        marginTop: '15px',
        marginLeft: '10px'
    }
}

export default class Timer extends React.Component {
    constructor(props) {
        super(props)
        this.state = {
            startTime: Date.now(),
            now: Date.now(),
        }
        this.handleResetTimer = this.handleResetTimer.bind(this)
        this.onResetTimer = this.onResetTimer.bind(this)
        this.updateNow = this.updateNow.bind(this)
    }

    componentDidMount() {
        // When the component is mounted, add your DOM listener to the "nv" elem.
        // (The "nv" elem is assigned in the render function.)
        window.document.addEventListener("resetTimer", this.handleResetTimer);
        this.interval = setInterval(this.updateNow, 100);
      }
    
    componentWillUnmount() {
        // Make sure to remove the DOM listener when the component is unmounted.
        window.document.removeEventListener("resetTimer", this.handleResetTimer);
        clearInterval(this.interval);
    }

    handleResetTimer(event) {
        this.onResetTimer()
    }

    updateNow() {
        this.setState({
            now: Date.now(),
        })
    }

    onResetTimer() {
        this.setState({
            now: Date.now(),
            startTime: Date.now(),
        })
    }

    render(){
        return (
            <Badge style={styles.button} count={parseFloat(Math.round((this.state.now - this.state.startTime)/100)/10).toFixed(1)}>
                <Tooltip placement="right" title="Elapsed Time">
                    <Button style={styles.button} icon="clock-circle" />
                </Tooltip>
            </Badge>
            
        )
    }
}

