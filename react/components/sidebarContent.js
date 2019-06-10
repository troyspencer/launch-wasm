import React from "react";
import StatsToggle from "./statsToggle";
import Stats from "./stats";
import FlexView from "react-flexview";

const styles = {
  header: {
    color: "grey",
    fontSize: "1.5em",
    textAlign: "center",
  }
}

export default class SidebarContent extends React.Component {
  constructor(props) {
    super(props);
    this.handleStatsChange = this.handleStatsChange.bind(this);
    this.onShowStatsChange = this.onShowStatsChange.bind(this);
  }

  handleStatsChange(show) {
    this.onShowStatsChange(show)
  }

  onShowStatsChange = (showStats) => {
    this.props.onShowStatsChange(showStats)
  }

  render() {
    return (
      <FlexView column marginTop="10" vAlignContent='center'>
        <div style={styles.header} >
          Settings 
        </div> 
        <StatsToggle showStats={this.props.showStats} onShowStatsChange={this.handleStatsChange} />
        <Stats showStats={this.props.showStats} />
      </FlexView>
    )
  }
}
