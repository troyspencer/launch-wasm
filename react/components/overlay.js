import React, { useState, useEffect } from "react";
import Sidebar from "react-sidebar";
import SettingsButton from "./settingsButton"
import Stats from "./stats";
import GameView from './gameView'
import SidebarContent from "./sidebarContent";
import FlexView from 'react-flexview';
import { Spin, Icon } from 'antd';

export default function Overlay() {
    const [showStats, setShowStats] = useState(true)
    const [sidebarOpen, setSidebarOpen] = useState(false)
    const [loading, setLoading] = useState(true)
    const [loaded, setLoaded] = useState(false)
    const [paused, setPaused] = useState(false)

    const styles = {
        spin: {
        marginTop: '20em',
        color: 'rgb(180,180,180',
        backgroundColor: 'black'
        }
    }

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

    useEffect(() => {
        setPaused(sidebarOpen)
    }, [sidebarOpen])

    return ( 
        <Sidebar
            sidebar={<SidebarContent showStats={showStats} onShowStatsChange={setShowStats} />}
            open={sidebarOpen}
            onSetOpen={setSidebarOpen}
            styles={{ sidebar: { background: "#464646", color: "rgb(180,180,180)", width: "10em"} }}>
            <Spin 
                tip="Loading WebAssembly..." 
                size="large" 
                spinning={loading}
                style={styles.spin}
                indicator={<Icon type="loading" spin />}>
                <GameView onLoadedChange={setLoaded} onLoadingChange={setLoading} />
                <FlexView hidden={!loaded || sidebarOpen} vAlignContent='top'>
                    <SettingsButton onClick={() => {setSidebarOpen(!sidebarOpen)}} />
                    <Stats paused={paused} showStats={showStats} />
                </FlexView>
            </Spin>
        </Sidebar>
    );
}