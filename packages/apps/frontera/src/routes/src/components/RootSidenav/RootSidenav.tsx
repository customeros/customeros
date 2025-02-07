import { useSearchParams } from 'react-router-dom';

import { observer } from 'mobx-react-lite';

import { useStore } from '@shared/hooks/useStore';
import { useKeyboardNavigation } from '@shared/components/RootSidenav/hooks/useKeyboardNavigation';
import {
  useNavigationManager,
  usePreferencesManager,
} from '@shared/components/RootSidenav/hooks';
import { SettingsSection } from '@shared/components/RootSidenav/components/sections/SettingsSection.tsx';
import {
  LogoSection,
  NavigationSections,
} from '@shared/components/RootSidenav/components/sections';

import { SystemUpdateNotification } from './components/SystemUpdateNotification';

export const RootSidenav = observer(() => {
  const { preferences, togglePreference } = usePreferencesManager();
  const { handleItemClick, checkIsActive } = useNavigationManager();
  const store = useStore();
  const [searchParams] = useSearchParams();

  const preset = searchParams?.get('preset');

  const tableViewDef = store.tableViewDefs.getById(preset ?? '1');
  const isPlatformOwner = store?.globalCache?.value?.isPlatformOwner ?? false;

  const presets = {
    targetsPreset: store.tableViewDefs.targetsPreset,
    customersPreset: store.tableViewDefs.customersPreset,
    organizationsPreset: store.tableViewDefs.organizationsPreset,
    upcomingInvoicesPreset: store.tableViewDefs.upcomingInvoicesPreset,
    contactsPreset: store.tableViewDefs.contactsPreset,
    contractsPreset: store.tableViewDefs.contractsPreset,
    flowSequencesPreset: store.tableViewDefs.flowsPreset,
  };

  useKeyboardNavigation(
    presets,
    {
      when:
        store.ui.isSearching !==
          tableViewDef?.value?.tableType?.toLowerCase() &&
        !store.ui.commandMenu.isOpen &&
        !store.ui.isEditingTableCell &&
        !store.ui.isFilteringTable,
    },
    isPlatformOwner,
  );

  return (
    <div className='pb-4 h-full w-12.5 bg-white flex flex-col border-r border-gray-200 overflow-hidden'>
      <LogoSection />

      <NavigationSections
        preferences={preferences}
        checkIsActive={checkIsActive}
        handleItemClick={handleItemClick}
        togglePreference={togglePreference}
      />
      {/* <UserActionSection /> */}
      <SystemUpdateNotification />
      <SettingsSection />
    </div>
  );
});
