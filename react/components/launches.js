import React, { useEffect} from 'react'
import {Button, Tooltip, Badge} from 'antd'

const styles = {
    button: {
        marginTop: '15px',
        marginLeft: '10px'
    }
}

export default function Launches(props) { 
    
    useEffect(() => {
        const updateLaunches = (e) => props.setLaunches(e.launches)
        window.document.addEventListener("updateLaunches", updateLaunches);
        return () => {
            window.document.removeEventListener("updateLaunches", updateLaunches);
        } 
    },[])

    return (
        <Badge style={styles.button} count={props.launches}>
            <Tooltip placement="bottom" title="Launches">
                <Button style={styles.button} type="primary" icon="rise" />
            </Tooltip>
        </Badge>
    )
}   
