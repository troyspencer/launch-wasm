import React, { useState } from 'react'
import Launches from './launches'
import Timer from './timer'
import { Row, Col } from 'antd'

export default function Stats(props) {
    const [launches,setLaunches] = useState(0)

    return (
        <div hidden={!props.showStats}>
            <Row>
                <Col span={12}>
                    <Launches launches={launches} onLaunchesChanged={setLaunches} />
                </Col>
                <Col span={12}>
                    <Timer launches={launches} paused={props.paused} />
                </Col>
            </Row>
        </div>
    )
}
