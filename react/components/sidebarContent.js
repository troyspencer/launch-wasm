import React from "react";
import StatsToggle from "./statsToggle";
import FlexView from "react-flexview/lib";
import { Header } from "semantic-ui-react";

const styles = {
  header: {
    color: "grey",
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
        <Header style={styles.header} textAlign="center" content="Settings" />
        <StatsToggle showStats={this.props.showStats} onShowStatsChange={this.handleStatsChange} />
      </FlexView>
    )
  }
}
