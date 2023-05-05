import React, {KeyboardEventHandler, useEffect, useRef} from "react";
import {AvailableAction} from './AvailableAction'
interface AvailableActionListProps {
    actionOptions:Array<string>,
    active: boolean,
    onKeyDown: KeyboardEventHandler
}
export const AvailableActionList = ({ actionOptions, active, onKeyDown}: AvailableActionListProps) => {
    return (
        <ul style={{padding: '0', margin: '0', width: '50%'}} role={'listbox'} >
            {actionOptions.map(action => (
                <AvailableAction
                    onKeyDown={onKeyDown}
                    action={action}
                    active={active}
                    key={action}/>
            ))}
        </ul>
    )
};