import React, { useState, useEffect } from 'react'
import {Button, Tooltip, Badge} from 'antd'

const styles = {
    button: {
        marginTop: '15px',
        marginLeft: '10px'
    }
}

export default function Timer() {
    const initialNow = Date.now()
    const [startTime,setStartTime] = useState(initialNow)
    const [now,setNow] = useState(initialNow)
    const [elapsedTime,setElapsedTime] = useState(0)

    const updateNow = () => {
        setNow(Date.now())
    }

    const handleResetTimer = () => {
        const newNow = Date.now()
        setStartTime(newNow)
        setNow(newNow)
    }

    const generateElapsedTime = () => {
        const cleanSeconds = Math.round((now - startTime)/1000)
        
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
        setElapsedTime(generateElapsedTime())
    },[now])

    useEffect(() => {
        window.document.addEventListener("resetTimer", handleResetTimer);
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