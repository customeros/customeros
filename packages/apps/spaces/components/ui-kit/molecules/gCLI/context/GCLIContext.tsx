import React, { createContext, ReactNode, useContext, useState } from 'react';

interface GCLIContextProviderProps {
  label: string;
  icon?: ReactNode;
  inputPlaceholder?: string;
  children: ReactNode;
  existingTerms?: Array<any>;
  loadSuggestions: (searchTerm: string) => void;
  loadingSuggestions: boolean;
  suggestionsLoaded: any[];
  selectedTermFormat?: (selectedTerm: any) => string;
  onItemsChange: (items: any[]) => void;
}

interface GCLIContextInterface {
  label: string;
  icon?: ReactNode;
  inputPlaceholder?: string;
  existingTerms?: Array<any>;
  loadSuggestions: (searchTerm: string) => void;
  loadingSuggestions: boolean;
  suggestionsLoaded: any[];
  selectedTermFormat?: (selectedTerm: any) => string;
  onItemsChange: (items: any[]) => void;
  selectedItems: Array<any>;
}

const GCLIContext = createContext<GCLIContextInterface>({
  label: '',
  icon: undefined,
  inputPlaceholder: undefined,
  existingTerms: [],
  loadSuggestions: (searchTerm: string) => {
    console.log(searchTerm);
  },
  loadingSuggestions: false,
  suggestionsLoaded: [],
  onItemsChange: (items: any[]) => {
    console.log(items);
  },
  selectedTermFormat: (selectedTerm: any) => {
    return selectedTerm.display;
  },
  selectedItems: [],
});

const GCLIContextProvider = ({
  label,
  icon,
  inputPlaceholder,
  children,
  existingTerms,
  loadSuggestions,
  loadingSuggestions,
  suggestionsLoaded,
  selectedTermFormat,
  onItemsChange,
}: GCLIContextProviderProps): JSX.Element => {
  const [selectedItems, setSelectedItems] = useState([]);

  return (
    <GCLIContext.Provider
      value={{
        label,
        icon,
        inputPlaceholder,
        existingTerms,
        loadSuggestions,
        loadingSuggestions,
        suggestionsLoaded,
        onItemsChange,
        selectedItems,
        selectedTermFormat,
      }}
    >
      {children}
    </GCLIContext.Provider>
  );
};

const useGCLI = () => {
  const {
    label,
    icon,
    inputPlaceholder,
    existingTerms,
    loadSuggestions,
    loadingSuggestions,
    suggestionsLoaded,
    onItemsChange,
    selectedItems,
    selectedTermFormat,
  } = useContext(GCLIContext);
  return {
    label,
    icon,
    inputPlaceholder,
    existingTerms,
    loadSuggestions,
    loadingSuggestions,
    suggestionsLoaded,
    onItemsChange,
    selectedItems,
    selectedTermFormat,
  };
};

export { GCLIContext, GCLIContextProvider, useGCLI };
