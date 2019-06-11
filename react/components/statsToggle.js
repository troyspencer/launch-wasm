import React from 'react'
import { Checkbox } from 'antd';

export default function StatsToggle(props) {
  const styles = {
    checkbox: {
      color: "rgb(180,180,180)",
      margin: "1em"
    }
  }

  return (
    <Checkbox style={styles.checkbox} checked={props.showStats} onChange={() => props.onShowStatsChange(!props.showStats)}>
      {props.showStats ? "Stats Visible" : "Stats Hidden"}
    </Checkbox>
  )
}
