import React, {createContext, ReactNode, useContext, useState} from 'react';
import {GCLIInputMode} from "../types";

interface GCLIContextProviderProps {
    label: string
    icon?: ReactNode
    inputPlaceholder?: string
    children: ReactNode
    loadSuggestions: (searchTerm: string) => void,
    loadingSuggestions: boolean,
    suggestionsLoaded: any[],
    onItemsChange: (items: any[]) => void
}

interface GCLIContextInterface {
    label: string
    icon?: ReactNode
    inputPlaceholder?: string
    loadSuggestions: (searchTerm: string) => void,
    loadingSuggestions: boolean,
    suggestionsLoaded: any[],
    onItemsChange: (items: any[]) => void
    selectedItems: Array<any>
    highlightedItemIndex: null | number
    onChangeMode: (mode: GCLIInputMode) => void
    onHighlightedItemChange: (newIndex: number | null) => void
    onSelectNextSuggestion: (
        {suggestions, currentIndex, onIndexSelect}:
            { suggestions: Array<any>, currentIndex: null | number, onIndexSelect: (index: number) => void }) => void // fixme
    onSelectPrevSuggestion: (
        {currentIndex, onIndexSelect}:
            { currentIndex: null | number, onIndexSelect: (index: number) => void }) => void // fixme
}

const GCLIContext = createContext<GCLIContextInterface>({
    label: '',
    icon: undefined,
    inputPlaceholder: undefined,
    loadSuggestions: (searchTerm: string) => {},
    loadingSuggestions: false,
    suggestionsLoaded: [],
    onItemsChange: (items: any[]) => {

    },
    selectedItems: [],
    highlightedItemIndex: null,
    onChangeMode: (newMode) => {
    },
    onHighlightedItemChange: (index) => {
    },
    onSelectNextSuggestion: (data: any) => {
    },
    onSelectPrevSuggestion: (data: any) => {
    }
});


const GCLIContextProvider = ({label, icon, inputPlaceholder, children, loadSuggestions, loadingSuggestions, suggestionsLoaded, onItemsChange}: GCLIContextProviderProps): JSX.Element => {
    const [mode, setMode] = useState<GCLIInputMode>('default');
    const [selectedItems, setSelectedItems] = useState([]);
    const [highlightedItemIndex, setHighlightedItemIndex] = useState<null | number>(null);

    const handleChangeMode = (newMode: GCLIInputMode) => {
        setMode(newMode);
    }

    const handleHighlightItem = (newItemIndex: number | null) => {
        console.log('HERE')
        setHighlightedItemIndex(newItemIndex);
    }

    //TODO fix types
    const handleSelectNextSuggestion = ({suggestions, currentIndex, onIndexSelect}: any) => {
        // select first item from the list -> if nothing is selected yet and there are available options
        if (currentIndex === null && suggestions?.length >= 0) {
            onIndexSelect(0)
        }
        // select next item if currently selected item is not last on the list
        if (suggestions.length - 1 >= currentIndex) {
            onIndexSelect(currentIndex + 1)
        } else {
            // go to top of the list ???
            onIndexSelect(0)
        }
    }

    const handleSelectPrevSuggestion = ({currentIndex, onIndexSelect}: any) => {
        // deselect list -> move focus back to input / previous context
        if (currentIndex === 0) {
            onIndexSelect(null)
        }
        // select prev
        if (currentIndex > 0) {
            onIndexSelect(currentIndex - 1)
        }
    }

    return (
        <GCLIContext.Provider value={{
            label,
            icon,
            inputPlaceholder,
            loadSuggestions,
            loadingSuggestions,
            suggestionsLoaded,
            onItemsChange,
            selectedItems,
            highlightedItemIndex,
            onChangeMode: handleChangeMode,
            onHighlightedItemChange: handleHighlightItem,
            onSelectNextSuggestion: handleSelectNextSuggestion,
            onSelectPrevSuggestion: handleSelectPrevSuggestion,
        }}>
            {children}
        </GCLIContext.Provider>
    );
}

const useGCLI = () => {
    const {
        label,
        icon,
        inputPlaceholder,
        loadSuggestions,
        loadingSuggestions,
        suggestionsLoaded,
        onItemsChange,
        selectedItems,
        highlightedItemIndex,
        onChangeMode,
        onHighlightedItemChange, onSelectNextSuggestion, onSelectPrevSuggestion
    } = useContext(GCLIContext)
    return {
        label,
        icon,
        inputPlaceholder,
        loadSuggestions,
        loadingSuggestions,
        suggestionsLoaded,
        onItemsChange,
        selectedItems,
        highlightedItemIndex,
        onChangeMode,
        onHighlightedItemChange,
        onSelectNextSuggestion,
        onSelectPrevSuggestion
    }
}

export {GCLIContext, GCLIContextProvider, useGCLI};