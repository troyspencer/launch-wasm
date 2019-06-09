import React from "react";
import Sidebar from "react-sidebar";
import SettingsButton from "./settingsButton"
import Game from './game'
import StatsToggle from "./statsToggle";
import FlexView from "react-flexview/lib";
import { Header } from "semantic-ui-react";

const styles = {
  header: {
    color: "grey",
    
  }
}

export default class SettingsSidebar extends React.Component {
  constructor(props) {
    super(props);
    this.state = {
      sidebarOpen: false,
      showStats: false,
    };
    this.handleSidebarChange = this.handleSidebarChange.bind(this);
    this.onSetSidebarOpen = this.onSetSidebarOpen.bind(this);
    this.handleStatsChange = this.handleStatsChange.bind(this);
    this.onShowStatsChange = this.onShowStatsChange.bind(this);
  }

  handleSidebarChange() {
    this.onSetSidebarOpen(!this.state.sidebarOpen);
  }

  onSetSidebarOpen(open) {
    this.setState({ sidebarOpen: open });
  }

  handleStatsChange(show) {
    this.onShowStatsChange(show)
  }

  onShowStatsChange = (showStats) => {
    this.setState({showStats: showStats})
  }

  render() {
    const sidebarContent = (
      <FlexView column marginTop="10" vAlignContent='center'>
        <Header style={styles.header} textAlign="center" content="Settings" />
        <StatsToggle showStats={this.state.showStats} onShowStatsChange={this.handleStatsChange} />
      </FlexView>
    )
    
    return (
      <Sidebar
        sidebar={sidebarContent}
        open={this.state.sidebarOpen}
        onSetOpen={this.onSetSidebarOpen}
        styles={{ sidebar: { background: "#464646", color: "grey"} }}
      >
        <Game showStats={this.state.showStats} />
        <SettingsButton onClick={this.handleSidebarChange} />
      </Sidebar>
    );
  }
}
