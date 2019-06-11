import React, { useState } from "react";
import Sidebar from "react-sidebar";
import SettingsButton from "./settingsButton"
import Stats from "./stats";
import Game from './game'
import SidebarContent from "./sidebarContent";
import FlexView from 'react-flexview';

export default function Overlay() {
  const [showStats, setShowStats] = useState(true)
  const [sidebarOpen, setSidebarOpen] = useState(false)

  const toggleSidebarOpen = () => {
    setSidebarOpen(!sidebarOpen)
  }

  return (
    <Sidebar
      sidebar={<SidebarContent showStats={showStats} onShowStatsChange={setShowStats} />}
      docked={sidebarOpen}
      open={sidebarOpen}
      onSetOpen={setSidebarOpen}
      styles={{ sidebar: { background: "#464646", color: "grey"} }}
    >
      <Game />
      <FlexView vAlignContent='top'>
        <SettingsButton onClick={toggleSidebarOpen} />
        <Stats showStats={showStats} />
      </FlexView>
    </Sidebar>
  );
}