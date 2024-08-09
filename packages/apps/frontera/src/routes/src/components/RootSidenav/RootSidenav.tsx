import React from 'react';

import { observer } from 'mobx-react-lite';

import {
  useNavigationManager,
  usePreferencesManager,
} from '@shared/components/RootSidenav/hooks';
import {
  LogoSection,
  UserActionSection,
  NavigationSections,
} from '@shared/components/RootSidenav/components/sections';

export const RootSidenav = observer(() => {
  const { preferences, togglePreference } = usePreferencesManager();
  const { handleItemClick, checkIsActive } = useNavigationManager();

  return (
    <div className='pb-4 h-full w-12.5 bg-white flex flex-col border-r border-gray-200 overflow-hidden'>
      <LogoSection />
      <NavigationSections
        preferences={preferences}
        checkIsActive={checkIsActive}
        handleItemClick={handleItemClick}
        togglePreference={togglePreference}
      />
      <UserActionSection />
    </div>
  );
});
