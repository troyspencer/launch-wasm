import React, { useState, useEffect } from "react";
import Sidebar from "react-sidebar";
import SettingsButton from "./settingsButton"
import Stats from "./stats";
import GameView from './gameView'
import SidebarContent from "./sidebarContent";
import FlexView from 'react-flexview';
import { Spin, Icon } from 'antd';

export default function Overlay(props) {
    const [showStats, setShowStats] = useState(true)
    const [sidebarOpen, setSidebarOpen] = useState(false)

    const styles = {
        spin: {
        marginTop: '20em',
        color: 'rgb(180,180,180',
        backgroundColor: 'black'
        }
    }

    useEffect(() => {
        if (props.setPaused) {
            props.setPaused(sidebarOpen)
        }
    }, [sidebarOpen, props.setPaused])

    return ( 
        <Sidebar
            sidebar={<SidebarContent showStats={showStats} setShowStats={setShowStats} />}
            open={sidebarOpen}
            onSetOpen={setSidebarOpen}
            styles={{ sidebar: { background: "#464646", color: "rgb(180,180,180)", width: "10em"} }}>
            <Spin 
                tip="Loading WebAssembly..." 
                size="large" 
                spinning={props.loading}
                style={styles.spin}
                indicator={<Icon type="loading" spin />}>
                <GameView 
                    setLoading={props.setLoading}
                    setLoaded={props.setLoaded}/>
                <FlexView hidden={!props.loaded || sidebarOpen} vAlignContent='top'>
                    <SettingsButton onClick={() => {setSidebarOpen(!sidebarOpen)}} />
                    <Stats paused={props.paused} showStats={showStats} setPaused={props.setPaused} />
                </FlexView>
            </Spin>
        </Sidebar>
    );
}