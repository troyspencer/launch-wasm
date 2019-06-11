import React, { useState } from "react";
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

  const styles = {
    spin: {
      margin: 0,
      position: "absolute",
      top: "50%",
      left: "50%",
      transform: "translate(-50%, -50%)"
    }
  }

  return (
    <Sidebar
      sidebar={<SidebarContent showStats={showStats} onShowStatsChange={setShowStats} />}
      docked={sidebarOpen}
      open={sidebarOpen}
      onSetOpen={setSidebarOpen}
      styles={{ sidebar: { background: "#464646", color: "rgb(180,180,180)", width: "10em"} }}
    >
      <Spin 
      tip="Loading..." 
      size="large" 
      spinning={loading}
      style={styles.spin}
      indicator={<Icon type="loading" spin />}>
        <Game onLoadingChange={setLoading} />
        <FlexView vAlignContent='top'>
          <SettingsButton onClick={() => {setSidebarOpen(!sidebarOpen)}} />
          <Stats showStats={showStats} />
        </FlexView>
      </Spin>
    </Sidebar>
  );
}