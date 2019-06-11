import React from 'react'
import FlexView from 'react-flexview'
import { FormCheckbox } from "shards-react";

export default function StatsToggle(props) {
    const toggleShowStats = () => {
        props.onShowStatsChange(!props.showStats)
    }

    return (
        <FlexView vAlignContent="center" hAlignContent='center'>
            <FlexView marginLeft="10" marginRight="10" vAlignContent="center" hAlignContent='center'>
                <FormCheckbox checked={props.showStats} onChange={toggleShowStats}>
                    {props.showStats ? "Stats Visible" : "Stats Hidden"}
                </FormCheckbox>
            </FlexView>
        </FlexView>  
    )
}
