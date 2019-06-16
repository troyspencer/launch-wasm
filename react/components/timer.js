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
    const [totalPausedTime,setTotalPausedTime] = useState(0)
    const [pauseStartedTime,setPauseStartedTime] = useState(initialNow)
    const [startedLevel, setStartedLevel] = useState(false)

    useEffect(() => {
        if (props.launches != 0) {
            setStartedLevel(true)
        }
    }, [props.launches])

    useEffect(() => {
        if (startedLevel) {
            setStartTime(Date.now())
        }
    }, [startedLevel])

    useEffect(() => {
        setNow(startTime)
        setPauseStartedTime(startTime)
        setTotalPausedTime(0)
    }, [startTime])

    useEffect(() => {
        if (!props.paused) {
            setPauseStartedTime(now)
        }
        if (!startedLevel) {
            setStartTime(now)
        }
    }, [now, props.paused, startedLevel])

    useEffect(() => {
        if (props.paused) {
            setPauseStartedTime(Date.now())
        } else {
            setTotalPausedTime(totalPausedTime + Date.now() - pauseStartedTime)
        }
    }, [props.paused])

    useEffect(() => {
        const handleResetTimer = () => { 
            setStartedLevel(false)
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

    const generateDisplayTime = (totalSeconds) => {
        if (totalSeconds <= 0) {
            return ''
        }

        var days   = Math.floor(totalSeconds / 86400);
        var hours   = Math.floor((totalSeconds - (days * 86400)) / 3600);
        var minutes = Math.floor((totalSeconds - (days * 86400) - (hours * 3600)) / 60);
        var seconds = totalSeconds - (days * 86400) - (hours * 3600) - (minutes * 60);
    
        if (minutes == 0 && hours == 0 && days == 0) {
            return seconds+''
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

    const generateElapsedTime = () => {
        const totalSeconds = Math.round((now - startTime - (totalPausedTime + now - pauseStartedTime))/1000)
        return generateDisplayTime(totalSeconds)
    }

    return (
        <Badge style={styles.button} count={generateElapsedTime()}>
            <Tooltip placement="bottom" title="Elapsed Time">
                <Button style={styles.button} type="primary" icon="clock-circle" />
            </Tooltip>
        </Badge>
    )
}