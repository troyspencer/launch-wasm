import React, { useState, useEffect } from 'react'
import {Button, Tooltip, Badge} from 'antd'

const styles = {
    button: {
        marginTop: '15px',
        marginLeft: '10px'
    }
}

export default function Timer(props) {
    const initialNow = Date.now()
    const [startTime,setStartTime] = useState(initialNow)
    const [now,setNow] = useState(initialNow)
    const [elapsedTime,setElapsedTime] = useState(0)
    const [pausedSeconds,setPausedSeconds] = useState(0)
    const [pauseStartedTime,setPauseStartedTime] = useState(initialNow)

    const generateElapsedTime = () => {
        const cleanSeconds = Math.round((now - startTime - pausedSeconds)/1000)
        
        var days   = Math.floor(cleanSeconds / 86400);
        var hours   = Math.floor((cleanSeconds - (days * 86400)) / 3600);
        var minutes = Math.floor((cleanSeconds - (days * 86400) - (hours * 3600)) / 60);
        var seconds = cleanSeconds - (days * 86400) - (hours * 3600) - (minutes * 60);
    
        if (minutes == 0 && hours == 0 && days == 0) {
            return seconds
        }
        if (seconds < 10) {seconds = "0"+seconds;}

        if (hours == 0 && days == 0) {
            return minutes+':'+seconds
        }
        if (minutes < 10) {minutes = "0"+minutes;}

        if (days == 0) {
            return hours+':'+minutes+':'+seconds;
        }
        if (hours < 10) {hours = "0"+hours;}

        return days+':'+hours+':'+minutes+':'+seconds;
    }





    useEffect(() => {
        setNow(startTime)
        setPausedSeconds(0)
        setPauseStartedTime(startTime)
    }, [startTime])

    useEffect(() => {
        // don't update clock if paused or on the starting block
        if (!props.paused && props.launches > 0) {
            setElapsedTime(generateElapsedTime())
        }
    },[now, props.paused, props.launches])

    useEffect(() => {
        if (props.launches != 1) {
            setStartTime(Date.now())
        }
    }, [props.launches])

    useEffect(() => {
        if (props.paused) {
            setPauseStartedTime(Date.now())
        } else {
            setPausedSeconds(pausedSeconds+Date.now()-pauseStartedTime)
        }
    }, [props.paused])

    useEffect(() => {
        const handleResetTimer = () => { 
            setStartTime(Date.now())
        }
        window.document.addEventListener("resetTimer", handleResetTimer);

        const updateNow = () => {
            setNow(Date.now())
        }
        const interval = setInterval(updateNow, 1000);
        return () => {
            window.document.removeEventListener("resetTimer", handleResetTimer);
            clearInterval(interval);
        }
    }, [])

    return (
        <Badge style={styles.button} count={elapsedTime}>
            <Tooltip placement="bottom" title="Elapsed Time">
                <Button style={styles.button} type="primary" icon="clock-circle" />
            </Tooltip>
        </Badge>
    )
}