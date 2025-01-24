import { useLocation } from 'react-router-dom';

import { observer } from 'mobx-react-lite';
import { useFeatureIsOn } from '@growthbook/growthbook-react';

import { cn } from '@ui/utils/cn';
import { Atom01 } from '@ui/media/icons/Atom01';
import { useStore } from '@shared/hooks/useStore';
import { Bubbles } from '@ui/media/icons/Bubbles';
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

export const NavigationSections = observer(
  ({
    preferences,
    togglePreference,
    handleItemClick,
    checkIsActive,
  }: NavigationSectionsProps) => {
    const store = useStore();
    const { pathname } = useLocation();

    const showCustomerMap = useFeatureIsOn('show-customer-map');
    const flowSequencesView = store.tableViewDefs.getById(
      store.tableViewDefs.flowsPreset ?? '',
    );
    const isFlowEditorActive = pathname.includes('flow-editor');

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
                <Bubbles
                  className={cn(
                    'size-4 min-w-4 text-gray-500',
                    isActive && 'text-gray-700',
                  )}
                />
              )}
            />
          )}

          <SidenavItem
            label='Flows'
            dataTest={`side-nav-item-all-flows`}
            onClick={() =>
              handleItemClick(`finder?preset=${flowSequencesView?.value?.id}`)
            }
            isActive={
              checkIsActive('finder', {
                preset: flowSequencesView?.value?.id ?? '',
              }) || isFlowEditorActive
            }
            icon={(isActive) => (
              <Atom01
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
  },
);
