import { cn } from '@ui/utils/cn.ts';
import { Bubbles } from '@ui/media/icons/Bubbles';
import { Preferences } from '@shared/components/RootSidenav/hooks';
import { SidenavItem } from '@shared/components/RootSidenav/components/SidenavItem';
import { TeamViewsSectionSection } from '@shared/components/RootSidenav/components/sections/TeamViewsSection.tsx';

import { FavoritesSection } from './FavoritesSection';
import { GeneralViewsSection } from './GeneralViewsSection';
import { LifecycleStagesSection } from './LifecycleStagesSection';

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
  return (
    <div className='px-2 pt-2.5 gap-4 overflow-y-auto overflow-hidden flex flex-col flex-1'>
      <SidenavItem
        label='Customer map'
        dataTest={`side-nav-item-customer-map`}
        isActive={checkIsActive('customer-map')}
        onClick={() => handleItemClick('customer-map')}
        icon={(isActive) => (
          <Bubbles
            className={cn(
              'size-4 min-w-4 text-gray-500',
              isActive && 'text-gray-700',
            )}
          />
        )}
      />
      <LifecycleStagesSection
        preferences={preferences}
        checkIsActive={checkIsActive}
        handleItemClick={handleItemClick}
        togglePreference={togglePreference}
      />
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
