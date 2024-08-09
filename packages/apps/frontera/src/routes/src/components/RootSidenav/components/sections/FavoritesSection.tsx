import React from 'react';

import { observer } from 'mobx-react-lite';

import { cn } from '@ui/utils/cn.ts';
import { useStore } from '@shared/hooks/useStore';
import { iconMap } from '@shared/components/RootSidenav/utils';
import { Preferences } from '@shared/components/RootSidenav/hooks';
import { EditableSideNavItem } from '@shared/components/RootSidenav/components/EditableSidenavItem';

import { CollapsibleSection } from '../CollapsibleSection';

interface FavoritesSectionProps {
  preferences: Preferences;
  handleItemClick: (data: string) => void;
  togglePreference: (data: keyof Preferences) => void;
  checkIsActive: (
    path: string,
    options?: { preset: string | Array<string> },
  ) => boolean;
}

export const FavoritesSection = observer(
  ({
    preferences,
    togglePreference,
    handleItemClick,
    checkIsActive,
  }: FavoritesSectionProps) => {
    const store = useStore();
    const tableViewDefsList = store.tableViewDefs.toArray();

    const favoritesView =
      tableViewDefsList
        .filter((c) => !c.value.isPreset)
        .sort((a, b) => a.value.order - b.value.order) ?? [];

    if (!favoritesView.length) return null;

    return (
      <CollapsibleSection
        title='Favorites'
        isOpen={preferences.isFavoritesOpen}
        onToggle={() => togglePreference('isFavoritesOpen')}
      >
        {preferences.isFavoritesOpen &&
          favoritesView.map((view) => (
            <EditableSideNavItem
              id={view.value.id}
              key={view.value.id}
              label={view.value.name}
              dataTest={`side-nav-item-${view.value.name}`}
              onClick={() => handleItemClick(`finder?preset=${view.value.id}`)}
              isActive={checkIsActive('finder', {
                preset: view.value.id,
              })}
              icon={(isActive) => {
                const Icon = iconMap[view.value.icon];

                return (
                  <Icon
                    className={cn(
                      'text-md text-gray-500',
                      isActive && 'text-gray-700',
                    )}
                  />
                );
              }}
            />
          ))}
      </CollapsibleSection>
    );
  },
);
