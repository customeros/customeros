import { useRef, useEffect } from 'react';
import { useSearchParams } from 'react-router-dom';

import { useDebounce } from 'rooks';
import { useLocalStorage } from 'usehooks-ts';

export const useSearchPersistence = () => {
  const [searchParams] = useSearchParams();
  const [lastSearchForPreset, setLastSearchForPreset] = useLocalStorage<{
    [key: string]: string;
  }>('customeros-last-search-for-preset', { root: 'root' });

  const preset = searchParams.get('preset');
  const searchValue = searchParams.get('search') ?? '';

  const previousPresetRef = useRef<string | null>(null);

  const saveToStorage = (currentPreset: string, currentSearchValue: string) => {
    if (currentPreset) {
      setLastSearchForPreset((prevState) => ({
        ...prevState,
        [currentPreset]: currentSearchValue,
      }));
    }
  };

  const debouncedSave = useDebounce(saveToStorage, 500);

  useEffect(() => {
    if (!preset) return;

    if (preset !== previousPresetRef.current) {
      debouncedSave.flush();
      previousPresetRef.current = preset;
    }

    debouncedSave(preset, searchValue);
  }, [preset, searchValue, debouncedSave]);

  return { lastSearchForPreset };
};
