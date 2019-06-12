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
      marginTop: "50%",
    }
  }

  return (
    <Sidebar
      sidebar={<SidebarContent showStats={showStats} onShowStatsChange={setShowStats} />}
      open={sidebarOpen}
      onSetOpen={setSidebarOpen}
      styles={{ sidebar: { background: "#464646", color: "rgb(180,180,180)", width: "10em"} }}
    >
      <Spin 
      tip="Loading WebAssembly..." 
      size="large" 
      spinning={loading}
      style={styles.spin}
      indicator={<Icon type="loading" spin />}>
        <Game onLoadingChange={setLoading} />
        <FlexView hidden={loading} vAlignContent='top'>
          <SettingsButton onClick={() => {setSidebarOpen(!sidebarOpen)}} />
          <Stats showStats={showStats} />
        </FlexView>
      </Spin>
    </Sidebar>
  );
}