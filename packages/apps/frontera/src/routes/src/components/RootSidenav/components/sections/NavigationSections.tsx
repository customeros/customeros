import { useFeatureIsOn } from '@growthbook/growthbook-react';

import { cn } from '@ui/utils/cn';
import { Icon } from '@ui/media/Icon';
import { Preferences } from '@shared/components/RootSidenav/hooks';
import { SidenavItem } from '@shared/components/RootSidenav/components/SidenavItem';
import { TeamViewsSectionSection } from '@shared/components/RootSidenav/components/sections/TeamViewsSection';

import { FavoritesSection } from './FavoritesSection';
import { GeneralViewsSection } from './GeneralViewsSection';

interface NavigationSectionsProps {
  preferences: Preferences;
  handleItemClick: (data: string) => void;
  togglePreference: (data: keyof Preferences) => void;
  checkIsActive: (
    path: string,
    options?: { preset: string | Array<string> },
  ) => boolean;
}

export const NavigationSections = ({
  preferences,
  togglePreference,
  handleItemClick,
  checkIsActive,
}: NavigationSectionsProps) => {
  const showCustomerMap = useFeatureIsOn('show-customer-map');

  return (
    <div className='px-2 pt-2.5 gap-4 overflow-y-auto overflow-hidden flex flex-col flex-1'>
      <div className='flex flex-col'>
        {showCustomerMap && (
          <SidenavItem
            label='Customer map'
            dataTest={`side-nav-item-customer-map`}
            isActive={checkIsActive('customer-map')}
            onClick={() => handleItemClick('customer-map')}
            icon={(isActive) => (
              <Icon
                name='bubbles'
                className={cn('text-gray-500', isActive && 'text-gray-700')}
              />
            )}
          />
        )}

        <SidenavItem
          label='Agents'
          isActive={checkIsActive('agents')}
          dataTest={`side-nav-item-all-agents`}
          onClick={() => handleItemClick('agents')}
          icon={(isActive) => (
            <Icon
              name='atom-01'
              className={cn(
                'size-4 min-w-4 text-gray-500',
                isActive && 'text-gray-700',
              )}
            />
          )}
        />
      </div>

      <TeamViewsSectionSection
        preferences={preferences}
        checkIsActive={checkIsActive}
        handleItemClick={handleItemClick}
        togglePreference={togglePreference}
      />
      <FavoritesSection
        preferences={preferences}
        checkIsActive={checkIsActive}
        handleItemClick={handleItemClick}
        togglePreference={togglePreference}
      />
      <GeneralViewsSection
        preferences={preferences}
        checkIsActive={checkIsActive}
        handleItemClick={handleItemClick}
        togglePreference={togglePreference}
      />
    </div>
  );
};
