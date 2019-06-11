import React from "react";
import StatsToggle from "./statsToggle";
import { Col } from "antd";

const styles = {
  header: {
    color: "rgb(180,180,180)",
    fontSize: "1.5em",
    textAlign: "center",
  }
}

export default function SidebarContent(props) {
    return (
      <Col>
        <div style={styles.header}>
          Settings 
        </div> 
        <StatsToggle showStats={props.showStats} onShowStatsChange={props.onShowStatsChange} />
      </Col>
    )
}
