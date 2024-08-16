import { useParams, useNavigate, useSearchParams } from 'react-router-dom';

import { observer } from 'mobx-react-lite';
import { useLocalStorage } from 'usehooks-ts';

import { Users02 } from '@ui/media/icons/Users02';
import { useStore } from '@shared/hooks/useStore';
import { Ticket02 } from '@ui/media/icons/Ticket02';
import { Trophy01 } from '@ui/media/icons/Trophy01';
import { Tooltip } from '@ui/overlay/Tooltip/Tooltip';
import { InfoSquare } from '@ui/media/icons/InfoSquare';
import { IconButton } from '@ui/form/IconButton/IconButton';
import { ActivityHeart } from '@ui/media/icons/ActivityHeart';
import { ArrowNarrowLeft } from '@ui/media/icons/ArrowNarrowLeft';
import { SidenavItem } from '@shared/components/RootSidenav/components/SidenavItem';
import { NotificationCenter } from '@shared/components/Notifications/NotificationCenter';
import { useKeyboardNavigation } from '@shared/components/RootSidenav/hooks/useKeyboardNavigation';

export const OrganizationSidenav = observer(() => {
  const navigate = useNavigate();
  const params = useParams();
  const [searchParams, setSearchParams] = useSearchParams();
  const store = useStore();

  const organization = store.organizations.value.get(params?.id as string);

  const [lastActivePosition, setLastActivePosition] = useLocalStorage(
    `customeros-player-last-position`,
    { [params?.id as string]: 'tab=about' },
  );

  const checkIsActive = (tab: string) => searchParams?.get('tab') === tab;

  const handleItemClick = (tab: string) => () => {
    const urlSearchParams = new URLSearchParams(searchParams?.toString());

    urlSearchParams.set('tab', tab);

    setLastActivePosition({
      ...lastActivePosition,
      [params?.id as string]: urlSearchParams.toString(),
    });
    setSearchParams(urlSearchParams);
  };

  if (!organization) return null;
  const parentOrg = organization?.parentCompanies?.[0];
  const parentOrgName = parentOrg?.name;
  const parentOrgId = parentOrg?.metadata?.id;

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
      !store.ui.commandMenu.isOpen &&
      !store.ui.isEditingTableCell &&
      !store.ui.isFilteringTable,
  });

  return (
    <div className='px-2 pt-[6px] pb-4 h-full w-200 flex flex-col grid-area-sidebar bg-white relative border-r border-gray-200'>
      <div className='flex gap-2 items-center mb-4'>
        <IconButton
          size='xs'
          variant='ghost'
          aria-label='Go back'
          icon={<ArrowNarrowLeft className=' text-gray-700 ' />}
          onClick={() => {
            navigate(`/${lastActivePosition?.root || 'finder'}`);
          }}
        />

        <div className='flex flex-col line-clamp-1'>
          {parentOrg && (
            <a
              href={`/organization/${parentOrgId}?tab=about`}
              className='text-xs text-gray-600 truncate no-underline '
            >
              {parentOrgName}
            </a>
          )}
          <Tooltip label={organization?.value.name ?? ''}>
            <span className='max-w-150px  font-semibold text-gray-700 truncate whitespace-nowrap '>
              {organization?.value.name || 'Organization'}
            </span>
          </Tooltip>
        </div>
      </div>
      <div className='space-y-1 w-full'>
        <SidenavItem
          label='About'
          dataTest='org-side-nav-item-about'
          onClick={handleItemClick('about')}
          icon={<InfoSquare className='size-5' />}
          isActive={checkIsActive('about') || !searchParams?.get('tab')}
        />
        <SidenavItem
          label='People'
          isActive={checkIsActive('people')}
          dataTest='org-side-nav-item-people'
          onClick={handleItemClick('people')}
          icon={<Users02 className='size-5' />}
        />
        <SidenavItem
          label='Account'
          dataTest='org-side-nav-item-account'
          onClick={handleItemClick('account')}
          icon={<ActivityHeart className='size-5' />}
          isActive={checkIsActive('account') || checkIsActive('invoices')}
        />
        <SidenavItem
          label='Success'
          isActive={checkIsActive('success')}
          dataTest='org-side-nav-item-success'
          onClick={handleItemClick('success')}
          icon={<Trophy01 className='size-5' />}
        />
        <SidenavItem
          label='Issues'
          isActive={checkIsActive('issues')}
          dataTest='org-side-nav-item-issues'
          onClick={handleItemClick('issues')}
          icon={<Ticket02 className='size-5' />}
        />
      </div>
      <div className='flex flex-col flex-grow justify-end'>
        <NotificationCenter />
      </div>
    </div>
  );
});
