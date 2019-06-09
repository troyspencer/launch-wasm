import React from 'react'
import { Loader } from 'semantic-ui-react'

export default class LoaderCentered extends React.Component {
    constructor(props) {
        super(props)
    }
      
    render() {
      return ( 
        <Loader active={this.props.isLoading} size='massive' />
      )
    }
}