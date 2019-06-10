import React from "react";
import Sidebar from "react-sidebar";
import SettingsButton from "./settingsButton"
import Stats from "./stats";
import Game from './game'
import SidebarContent from "./sidebarContent";
import FlexView from 'react-flexview';

export default class Overlay extends React.Component {
  constructor(props) {
    super(props);
    this.state = {
      sidebarOpen: false,
      showStats: true,
    };
    this.handleSidebarChange = this.handleSidebarChange.bind(this);
    this.onSidebarChange = this.onSidebarChange.bind(this);
    this.handleStatsChange = this.handleStatsChange.bind(this);
    this.onShowStatsChange = this.onShowStatsChange.bind(this);
  }

  handleSidebarChange() {
    this.onSidebarChange(!this.state.sidebarOpen);
  }

  onSidebarChange(open) {
    this.setState({ sidebarOpen: open });
  }

  handleStatsChange(show) {
    this.onShowStatsChange(show)
  }

  onShowStatsChange(showStats) {
    this.setState({showStats: showStats})
  }

  render() {
    
    return (
      <Sidebar
        sidebar={<SidebarContent showStats={this.state.showStats} onShowStatsChange={this.handleStatsChange} />}
        docked={this.state.sidebarOpen}
        open={this.state.sidebarOpen}
        onSetOpen={this.onSidebarChange}
        styles={{ sidebar: { background: "#464646", color: "grey"} }}
      >
        <Game />
        <FlexView vAlignContent='top'>
          <SettingsButton onClick={this.handleSidebarChange} />
          <Stats showStats={this.state.showStats} />
        </FlexView>
        
      </Sidebar>
    );
  }
}
