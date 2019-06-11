import React from 'react'
import Launches from './launches'
import Timer from './timer'
import { Row, Col } from 'antd'

export default function Stats(props) {
    return (
        <div hidden={!props.showStats}>
            <Row>
                <Col span={12}>
                    <Launches />
                </Col>
                <Col span={12}>
                    <Timer />
                </Col>
            </Row>
        </div>
    )
}
