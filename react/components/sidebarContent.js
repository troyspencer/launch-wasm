import React from "react";
import StatsToggle from "./statsToggle";
import FlexView from "react-flexview";

const styles = {
  header: {
    color: "grey",
    fontSize: "1.5em",
    textAlign: "center",
  }
}

export default function SidebarContent(props) {
    return (
      <FlexView column marginTop="10" vAlignContent='center'>
        <div style={styles.header} >
          Settings 
        </div> 
        <StatsToggle showStats={props.showStats} onShowStatsChange={props.onShowStatsChange} />
      </FlexView>
    )
}
