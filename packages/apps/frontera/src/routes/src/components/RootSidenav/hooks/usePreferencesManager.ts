import { useLocalStorage } from 'usehooks-ts';

interface Preferences {
  isViewsOpen: boolean;
  isMyViewsOpen: boolean;
  isFavoritesOpen: boolean;
  isLifecycleViewsOpen: boolean;
}

export const usePreferencesManager = () => {
  const [preferences, setPreferences] = useLocalStorage(
    'customeros-preferences',
    {
      isLifecycleViewsOpen: true,
      isMyViewsOpen: true,
      isViewsOpen: true,
      isFavoritesOpen: true,
    } as Preferences,
  );

  const togglePreference = (key: keyof Preferences) => {
    setPreferences((prev) => ({
      ...prev,
      [key]: !prev?.[key],
    }));
  };

  return { preferences, togglePreference };
};

export type { Preferences };
