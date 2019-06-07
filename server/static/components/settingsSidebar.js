import React from "react";
import Sidebar from "react-sidebar";
import SettingsButton from "./settingsButton"

export default class SettingsSidebar extends React.Component {
  constructor(props) {
    super(props);
    this.state = {
      sidebarOpen: false
    };
    this.onSetSidebarOpen = this.onSetSidebarOpen.bind(this);
  }

  onSetSidebarOpen(open) {
    this.setState({ sidebarOpen: open });
  }

  onClickSettings = () => {
    console.log("ClickedJS")
    const event = new Event("increment", {"test": true})
    window.dispatchEvent(event)
    this.onSetSidebarOpen(true)
}

  render() {
    return (
      <Sidebar
        sidebar={<b>Sidebar content</b>}
        open={this.state.sidebarOpen}
        onSetOpen={this.onSetSidebarOpen}
        styles={{ sidebar: { background: "white" } }}
      >
        {this.props.children}
        <SettingsButton onClick={this.onClickSettings} />
      </Sidebar>
    );
  }
}
