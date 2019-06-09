import React from 'react'
import Canvas from './canvas'
import Stats from './stats';

export default class Game extends React.Component {
  constructor(props) {
    super(props)
  }

  onShowStatsChange = (showStats) => {
    this.setState({showStats: showStats})
  }
  
  render() {
    return ( 
      <div>
        <div>
          <Canvas />
          <Stats showStats={this.props.showStats} />
        </div>
      </div>
    )
  }
}
