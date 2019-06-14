import React, {useEffect, useState} from 'react'
import Overlay from './overlay';

export default function Game() {
    const [loading, setLoading] = useState(true)
    const [loaded, setLoaded] = useState(false)
    const [paused, setPaused] = useState(false)

    const pauseEvent = new Event("pause")
    const unpauseEvent = new Event("unpause")

    useEffect(() => {
        if (paused) {
            window.document.dispatchEvent(pauseEvent)
        } else {
            window.document.dispatchEvent(unpauseEvent)
        }
    }, [paused])

    const handleKey = (e) => {
        if (e.which == 32) {
            setPaused(!paused)
        }
    }

    useEffect(() => {
        window.addEventListener("keyup", handleKey);
        return () => {
            window.removeEventListener("keyup", handleKey);
        }
    }, [paused])

    return (
        <Overlay 
        paused={paused} onPausedChange={setPaused}
        loading={loading} onLoadingChange={setLoading}
        loaded={loaded} onLoadedChange={setLoaded} />
    );
}