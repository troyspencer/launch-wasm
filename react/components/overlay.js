import React from "react";
import Sidebar from "react-sidebar";
import SettingsButton from "./settingsButton"
import Game from './game'
import SidebarContent from "./sidebarContent";

export default class Overlay extends React.Component {
  constructor(props) {
    super(props);
    this.state = {
      sidebarOpen: false,
      showStats: false,
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
        open={this.state.sidebarOpen}
        onSetOpen={this.onSidebarChange}
        styles={{ sidebar: { background: "#464646", color: "grey"} }}
      >
        <Game />
        <SettingsButton onClick={this.handleSidebarChange} />
      </Sidebar>
    );
  }
}
