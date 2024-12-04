import React from 'react';
import { useLocation } from 'react-router-dom';

import { observer } from 'mobx-react-lite';

import { cn } from '@ui/utils/cn.ts';
import { Play } from '@ui/media/icons/Play.tsx';
import { useStore } from '@shared/hooks/useStore';
import { iconMap } from '@shared/components/RootSidenav/utils';
import { Preferences } from '@shared/components/RootSidenav/hooks';
import { EditableSideNavItem } from '@shared/components/RootSidenav/components/EditableSidenavItem';
import { WelcomeSidenavItem } from '@shared/components/RootSidenav/components/WelcomeSideNavItem.tsx';

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
    const { pathname } = useLocation();

    const favoritesView =
      tableViewDefsList
        .filter((c) => !c.value.isPreset && !c.value.isShared)
        .sort((a, b) => a.value.order - b.value.order) ?? [];

    if (
      !favoritesView.length &&
      store.globalCache.value?.user?.onboarding?.showOnboardingPage
    )
      return null;

    return (
      <CollapsibleSection
        title='My views'
        isOpen={preferences.isFavoritesOpen}
        onToggle={() => togglePreference('isFavoritesOpen')}
      >
        {store.globalCache.value?.user?.onboarding?.showOnboardingPage && (
          <WelcomeSidenavItem
            label='Welcome'
            isActive={pathname.includes('welcome')}
            onClick={() => handleItemClick(`welcome`)}
            icon={(isActive) => (
              <Play
                className={cn(
                  'size-4 min-w-4 text-gray-500',
                  isActive && 'text-gray-700',
                )}
              />
            )}
          />
        )}

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
                const Icon = iconMap?.[view.value.icon];

                if (!Icon) return <div />;

                return (
                  <Icon
                    className={cn(
                      'size-4 min-w-4 text-gray-500',
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
