import React from 'react'
import ReactDom from 'react-dom'
import App from './app'
import "bootstrap/dist/css/bootstrap.min.css";
import "shards-ui/dist/css/shards.min.css"
import 'antd/dist/antd.css';

ReactDom.render(
  <App />,
  document.getElementById('app')
)