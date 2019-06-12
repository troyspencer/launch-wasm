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
    
    const updateLaunches = (e) => setLaunchCount(e.launches)
    useEffect(() => {
        window.document.addEventListener("updateLaunches", updateLaunches);
        return () => {
            window.document.removeEventListener("updateLaunches", updateLaunches);
        } 
    },[])

    return (
        <Badge style={styles.button} count={launchCount}>
            <Tooltip placement="bottom" title="Launches">
                <Button style={styles.button} type="primary" icon="rise" />
            </Tooltip>
        </Badge>
    )
}   
