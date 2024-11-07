import { useNavigate, useSearchParams } from 'react-router-dom';

import { useLocalStorage } from 'usehooks-ts';

import { cn } from '@ui/utils/cn';
import { Link01 } from '@ui/media/icons/Link01';
import { Receipt } from '@ui/media/icons/Receipt';
import { useStore } from '@shared/hooks/useStore';
import { Dataflow03 } from '@ui/media/icons/Dataflow03';
import { IconButton } from '@ui/form/IconButton/IconButton';
import { ArrowNarrowLeft } from '@ui/media/icons/ArrowNarrowLeft';
import { SidenavItem } from '@shared/components/RootSidenav/components/SidenavItem';
import { useKeyboardNavigation } from '@shared/components/RootSidenav/hooks/useKeyboardNavigation';

import { WorkspaceSection } from './components';
// import { FieldsSection } from './components/FieldsSection';

export const SettingsSidenav = () => {
  const navigate = useNavigate();
  const store = useStore();
  const [searchParams, setSearchParams] = useSearchParams();

  const [lastActivePosition, setLastActivePosition] = useLocalStorage(
    `customeros-player-last-position`,
    { ['settings']: 'oauth', root: 'organization' },
  );

  const checkIsActive = (tab: string) => searchParams?.get('tab') === tab;

  const handleItemClick = (tab: string) => () => {
    const params = new URLSearchParams(searchParams?.toString() ?? '');

    params.set('tab', tab);
    setLastActivePosition({ ...lastActivePosition, settings: tab });
    setSearchParams(params.toString());
  };

  const presets = {
    targetsPreset: store.tableViewDefs.targetsPreset,
    customersPreset: store.tableViewDefs.defaultPreset,
    organizationsPreset: store.tableViewDefs.organizationsPreset,
    upcomingInvoicesPreset: store.tableViewDefs.upcomingInvoicesPreset,
    contractsPreset: store.tableViewDefs.contractsPreset,
    flowSequencesPreset: store.tableViewDefs.flowsPreset,
  };

  useKeyboardNavigation(presets, {
    when:
      !store.ui.commandMenu.isOpen &&
      !store.ui.isEditingTableCell &&
      !store.ui.isFilteringTable,
  });

  return (
    <div className='px-2 pt-[6px] h-full w-[200px] bg-white flex flex-col relative border-r border-gray-200'>
      <div className='flex gap-2 items-center mb-4'>
        <IconButton
          size='xs'
          variant='ghost'
          aria-label='Go back'
          dataTest='settings-go-back'
          icon={<ArrowNarrowLeft className='text-gray-700' />}
          onClick={() => navigate(`/${lastActivePosition.root}`)}
        />

        <p className='font-semibold text-gray-700 break-keep line-clamp-1'>
          Settings
        </p>
      </div>

      <div className='flex flex-col space-y-2 w-full'>
        <WorkspaceSection
          checkIsActive={checkIsActive}
          handleItemClick={handleItemClick}
        />
        <SidenavItem
          label='Accounts'
          dataTest='settings-accounts'
          onClick={handleItemClick('oauth')}
          isActive={checkIsActive('oauth') || !searchParams?.get('tab')}
          icon={
            <Link01
              className={cn(
                checkIsActive('oauth') ? 'text-gray-700' : 'text-gray-500',
                'size-5',
              )}
            />
          }
        />
        {/* <FieldsSection
          checkIsActive={checkIsActive}
          handleItemClick={handleItemClick}
        /> */}
        <SidenavItem
          label='Customer billing'
          isActive={checkIsActive('billing')}
          onClick={handleItemClick('billing')}
          icon={
            <Receipt
              className={cn(
                checkIsActive('billing') ? 'text-gray-700' : 'text-gray-500',
                'size-5',
              )}
            />
          }
        />
        <SidenavItem
          label='Integrations'
          isActive={checkIsActive('integrations')}
          onClick={handleItemClick('integrations')}
          icon={
            <Dataflow03
              className={cn(
                checkIsActive('integrations')
                  ? 'text-gray-700'
                  : 'text-gray-500',
              )}
            />
          }
        />
      </div>
      <div className='flex flex-col space-y-1 flex-grow justify-end'>
        {/* <NotificationCenter /> */}
      </div>
    </div>
  );
};
