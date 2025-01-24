import { useLocation } from 'react-router-dom';

import { observer } from 'mobx-react-lite';

import { cn } from '@ui/utils/cn';
import { Icon } from '@ui/media/Icon';
import { useStore } from '@shared/hooks/useStore';
import { Preferences } from '@shared/components/RootSidenav/hooks';
import { iconNameMap } from '@shared/components/RootSidenav/utils';
import { WelcomeSidenavItem } from '@shared/components/RootSidenav/components/WelcomeSideNavItem';
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
    const { pathname } = useLocation();

    const favoritesView =
      tableViewDefsList
        .filter((c) => !c.value.isPreset && !c.value.isShared)
        .sort((a, b) => a.value.order - b.value.order) ?? [];

    if (
      !favoritesView.length &&
      !store.globalCache.value?.user?.onboarding?.showOnboardingPage
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
              <Icon
                name='play'
                className={cn('text-gray-500', isActive && 'text-gray-700')}
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
              icon={(isActive) => (
                <Icon
                  name={iconNameMap?.[view.value.icon] ?? 'building-07'}
                  className={cn(
                    'min-w-4 text-gray-500',
                    isActive && 'text-gray-700',
                  )}
                />
              )}
            />
          ))}
      </CollapsibleSection>
    );
  },
);
