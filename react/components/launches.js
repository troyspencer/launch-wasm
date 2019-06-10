import React, {useState, useEffect} from 'react'
import {Button, Tooltip, Badge} from 'antd'


const styles = {
    button: {
        marginTop: '15px',
        marginLeft: '10px'
    }
}

export default function Launches() {
    const [launchCount,setLaunchCount] = useState(0)
    
    useEffect(() => {
        window.document.addEventListener("updateLaunches", (e) => setLaunchCount(e.launches));
        return () => {
            window.document.removeEventListener("updateLaunches", (e) => setLaunchCount(e.launches));
        } 
    },[])

    return (
        <Badge style={styles.button} count={launchCount}>
            <Tooltip placement="right" title="Launches">
                <Button style={styles.button} icon="rise" />
            </Tooltip>
        </Badge>
    )
}   
