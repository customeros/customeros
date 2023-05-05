import React, {KeyboardEventHandler, useEffect, useRef} from "react";

interface ActionItemSuggestionProps {
    action:string,
    active: boolean,
    onKeyDown: KeyboardEventHandler
}
export const AvailableAction = ({ action, active, onKeyDown,}: ActionItemSuggestionProps) => {
    const ref = useRef(null)
    useEffect(() => {
        if(ref.current) {
            if(active) {
                //@ts-ignore
                ref.current?.focus()
            }
        }
    },[active])

    return (
        <li
            tabIndex={0}
            ref={ref}
            className={`list_item`}
            onClick={() => console.log('EXECUTING ACTION')}
            role="listitem"
            onKeyDown={onKeyDown}
        >
            {action}
        </li>
    )
};
