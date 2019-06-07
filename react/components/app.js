import React from 'react'
import SettingsSidebar from './settingsSidebar';
import Canvas from './canvas'

export default class App extends React.Component {
  render() {
    return ( 
      <div className="App">
        <SettingsSidebar>
          <Canvas />
        </SettingsSidebar>
      </div>
    )
  }
}
