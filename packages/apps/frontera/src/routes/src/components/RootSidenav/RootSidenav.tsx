import React from 'react';
import { useSearchParams } from 'react-router-dom';

import { observer } from 'mobx-react-lite';

import { useStore } from '@shared/hooks/useStore';
import { useKeyboardNavigation } from '@shared/components/RootSidenav/hooks/useKeyboardNavigation';
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
  const store = useStore();
  const [searchParams] = useSearchParams();

  const preset = searchParams?.get('preset');

  const tableViewDef = store.tableViewDefs.getById(preset ?? '1');

  const presets = {
    targetsPreset: store.tableViewDefs.targetsPreset,
    churnedPreset: store.tableViewDefs.churnedPreset,
    customersPreset: store.tableViewDefs.defaultPreset,
    addressBookPreset: store.tableViewDefs.addressBookPreset,
    upcomingInvoicesPreset: store.tableViewDefs.upcomingInvoicesPreset,
    myPortfolioPreset: store.tableViewDefs.myPortfolioPreset,
  };

  useKeyboardNavigation(presets, {
    when:
      store.ui.isSearching !== tableViewDef?.value?.tableType?.toLowerCase() &&
      !store.ui.commandMenu.isOpen &&
      !store.ui.isEditingTableCell &&
      !store.ui.isFilteringTable,
  });

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
