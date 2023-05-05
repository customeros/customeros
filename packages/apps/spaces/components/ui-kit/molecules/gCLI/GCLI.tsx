import React, {ReactNode} from 'react';
import {GCLIContextProvider} from "./context/GCLIContext";
import {GCLIInput} from "./GCLIInput";

interface GCLIProps {
    label: string,
    icon?: ReactNode,
    loadSuggestions: (searchTerm: string) => void,
    loadingSuggestions: boolean,
    suggestionsLoaded: any[],
    onItemsChange: (items: any[]) => void
}

export const GCLI = ({label, icon, suggestionsLoaded, loadingSuggestions, loadSuggestions, onItemsChange}: GCLIProps) => {
    return (
        <GCLIContextProvider label={label} icon={icon} suggestionsLoaded={suggestionsLoaded} loadingSuggestions={loadingSuggestions} loadSuggestions={loadSuggestions} onItemsChange={onItemsChange}>
            <GCLIInput/>
        </GCLIContextProvider>
    );
}

