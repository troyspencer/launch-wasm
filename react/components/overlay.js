import React, { useState, useEffect } from "react";
import Sidebar from "react-sidebar";
import SettingsButton from "./settingsButton"
import Stats from "./stats";
import Game from './game'
import SidebarContent from "./sidebarContent";
import FlexView from 'react-flexview';
import { Spin, Icon } from 'antd';

export default function Overlay() {
    const [showStats, setShowStats] = useState(true)
    const [sidebarOpen, setSidebarOpen] = useState(false)
    const [loading, setLoading] = useState(true)
    const [loaded, setLoaded] = useState(false)

    const styles = {
        spin: {
        marginTop: '20em',
        color: 'rgb(180,180,180',
        backgroundColor: 'black'
        }
    }

    useEffect(() => {
        if (sidebarOpen) {
            const event = new Event("pause")
            window.document.dispatchEvent(event)
        } else {
            const event = new Event("unpause")
            window.document.dispatchEvent(event)
        }
        
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
            <Game onLoadedChange={setLoaded} onLoadingChange={setLoading} />
            <FlexView hidden={!loaded || sidebarOpen} vAlignContent='top'>
            <SettingsButton onClick={() => {setSidebarOpen(!sidebarOpen)}} />
            <Stats paused={sidebarOpen} showStats={showStats} />
            </FlexView>
        </Spin>
    </Sidebar>
  );
}