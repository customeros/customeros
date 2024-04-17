'use client';
import Image from 'next/image';
import React, { useMemo, useEffect } from 'react';
import { useRouter, usePathname, useSearchParams } from 'next/navigation';

import { produce } from 'immer';
import { signOut } from 'next-auth/react';
import { useLocalStorage } from 'usehooks-ts';
import { useFeatureIsOn } from '@growthbook/growthbook-react';
import { useTenantSettingsQuery } from '@settings/graphql/getTenantSettings.generated';

import { cn } from '@ui/utils/cn';
import { Icons } from '@ui/media/Icon';
import { Skeleton } from '@ui/feedback/Skeleton';
import { Bubbles } from '@ui/media/icons/Bubbles';
import { LogOut01 } from '@ui/media/icons/LogOut01';
import { InvoiceCheck } from '@ui/media/icons/InvoiceCheck';
import { ArrowDropdown } from '@ui/media/icons/ArrowDropdown';
import { mockedTableDefs } from '@shared/util/tableDefs.mock';
import { getGraphQLClient } from '@shared/util/getGraphQLClient';
import { InvoiceUpcoming } from '@ui/media/icons/InvoiceUpcoming';
import { ClockFastForward } from '@ui/media/icons/ClockFastForward';
import { useOrganizationsMeta } from '@shared/state/OrganizationsMeta.atom';
import { useGlobalCacheQuery } from '@shared/graphql/global_Cache.generated';
import { NotificationCenter } from '@shared/components/Notifications/NotificationCenter';
import { useGetAllInvoicesCountQuery } from '@shared/graphql/getAllInvoicesCount.generated';

import { SidenavItem } from './components/SidenavItem';
import logoCustomerOs from './assets/logo-customeros.png';
import { GoogleSidebarNotification } from './components/GoogleSidebarNotification';

export const RootSidenav = () => {
  const client = getGraphQLClient();
  const router = useRouter();
  const pathname = usePathname();
  const searchParams = useSearchParams();
  const [_, setOrganizationsMeta] = useOrganizationsMeta();
  const showMyViewsItems = useFeatureIsOn('my-views-nav-item');

  const [lastActivePosition, setLastActivePosition] = useLocalStorage(
    `customeros-player-last-position`,
    { root: 'organization' },
  );
  const [preferences, setPreferences] = useLocalStorage(
    'customeros-preferences',
    {
      isInvoicesOpen: true,
      isMyViewsOpen: true,
    },
  );

  const { data: tenantSettingsData } = useTenantSettingsQuery(client);
  const { data: totalInvoices } = useGetAllInvoicesCountQuery(client);
  const tableViewDefsData = mockedTableDefs;

  const { data, isLoading } = useGlobalCacheQuery(client);
  const globalCache = data?.global_Cache;
  const myViews =
    tableViewDefsData.filter((c) => ['1', '2', '3'].includes(c.id)) ?? [];

  const handleItemClick = (path: string) => {
    setLastActivePosition({ ...lastActivePosition, root: path });
    setOrganizationsMeta((prev) =>
      produce(prev, (draft) => {
        draft.getOrganization.pagination.page = 1;
      }),
    );

    router.push(`/${path}`);
  };
  const checkIsActive = (path: string, options?: { preset: string }) => {
    const [_pathName, _searchParams] = path.split('?');
    const presetParam = new URLSearchParams(searchParams?.toString()).get(
      'preset',
    );

    if (options?.preset) {
      return (
        pathname?.startsWith(`/${_pathName}`) && presetParam === options.preset
      );
    } else {
      return pathname?.startsWith(`/${_pathName}`) && !presetParam;
    }
  };

  const handleSignOutClick = () => {
    signOut();
  };

  const showInvoices = useMemo(() => {
    return (
      tenantSettingsData?.tenantSettings?.billingEnabled ||
      totalInvoices?.invoices?.totalElements > 0
    );
  }, [tenantSettingsData?.tenantSettings?.billingEnabled]);

  useEffect(() => {
    [
      '/organizations',
      '/organizations?preset=customer',
      '/organizations?preset=portfolio',
      '/renewals?preset=1',
      '/renewals?preset=2',
      '/renewals?preset=3',
      '/invoices?preset=4',
      '/invoices?preset=5',
    ].forEach((path) => {
      router.prefetch(path);
    });
  }, []);

  const cdnLogoUrl = data?.global_Cache?.cdnLogoUrl;

  return (
    <div className='px-2 pt-2.5 pb-4 h-full w-12.5 bg-white flex flex-col border-r border-gray-200'>
      <div className='mb-2 ml-3 cursor-pointer flex justify-flex-start overflow-hidden relative'>
        {!isLoading ? (
          <Image
            src={cdnLogoUrl ?? logoCustomerOs}
            alt='CustomerOS'
            width={136}
            height={30}
            className='pointer-events-none transition-opacity-250 ease-in-out h-[30px] w-auto'
          />
        ) : (
          <Skeleton className='w-full h-8 mr-2' />
        )}
      </div>

      <div className='space-y-1 w-full mb-4'>
        <SidenavItem
          label='Customer map'
          isActive={checkIsActive('customer-map')}
          onClick={() => handleItemClick('customer-map')}
          icon={(isActive) => (
            <Bubbles
              className={cn(
                'w-5 h-5 text-gray-500',
                isActive && 'text-gray-700',
              )}
            />
          )}
        />
        <SidenavItem
          label='Organizations'
          isActive={checkIsActive('organizations')}
          onClick={() => handleItemClick('organizations')}
          icon={(isActive) => (
            <Icons.Building7
              className={cn(
                'w-5 h-5 text-gray-500',
                isActive && 'text-gray-700',
              )}
            />
          )}
        />
        <SidenavItem
          label='Customers'
          isActive={checkIsActive('organizations', { preset: 'customer' })}
          onClick={() => handleItemClick('organizations?preset=customer')}
          icon={(isActive) => (
            <Icons.CheckHeart
              className={cn(
                'w-5 h-5 text-gray-500',
                isActive && 'text-gray-700',
              )}
            />
          )}
        />

        {showInvoices && (
          <>
            <div
              className='w-full pt-3 gap-1 flex justify-flex-start pl-3.5 cursor-pointer text-gray-500 hover:text-gray-700 transition-colors'
              onClick={() =>
                setPreferences((prev) => ({
                  ...prev,
                  isInvoicesOpen: !prev.isInvoicesOpen,
                }))
              }
            >
              <span className='text-sm'>Invoices</span>
              <ArrowDropdown className='w-5 h-5' />
            </div>

            {preferences.isInvoicesOpen && (
              <>
                <SidenavItem
                  label='Upcoming'
                  isActive={checkIsActive('invoices', { preset: '4' })}
                  onClick={() => handleItemClick('invoices?preset=4')}
                  icon={(isActive) => (
                    <InvoiceUpcoming
                      className={cn(
                        'w-5 h-5 text-gray-500',
                        isActive && 'text-gray-700',
                      )}
                    />
                  )}
                />
                <SidenavItem
                  label='Past'
                  isActive={checkIsActive('invoices', { preset: '5' })}
                  onClick={() => handleItemClick('invoices?preset=5')}
                  icon={(isActive) => (
                    <InvoiceCheck
                      className={cn(
                        'w-5 h-5 text-gray-500',
                        isActive && 'text-gray-700',
                      )}
                    />
                  )}
                />
              </>
            )}
          </>
        )}
      </div>

      <div className='space-y-1 w-full'>
        {(globalCache?.isOwner || showMyViewsItems) && (
          <div
            className='w-full gap-1 flex justify-flex-start pl-3.5 cursor-pointer text-gray-500 hover:text-gray-700 transition-colors'
            onClick={() =>
              setPreferences((prev) => ({
                ...prev,
                isMyViewsOpen: !prev.isMyViewsOpen,
              }))
            }
          >
            <span className='text-sm'>My views</span>
            <ArrowDropdown className='w-5 h-5' />
          </div>
        )}

        {preferences.isMyViewsOpen && (
          <>
            {globalCache?.isOwner && (
              <SidenavItem
                label='My portfolio'
                isActive={checkIsActive('organizations', {
                  preset: 'portfolio',
                })}
                onClick={() =>
                  handleItemClick('organizations?preset=portfolio')
                }
                icon={(isActive) => (
                  <Icons.Briefcase1
                    className={cn(
                      'w-5 h-5 text-gray-500',
                      isActive && 'text-gray-700',
                    )}
                  />
                )}
              />
            )}
            {showMyViewsItems &&
              myViews.map((view) => (
                <SidenavItem
                  key={view.id}
                  label={view.name}
                  isActive={checkIsActive('renewals', { preset: view.id })}
                  onClick={() => handleItemClick(`renewals?preset=${view.id}`)}
                  icon={(isActive) => (
                    <ClockFastForward
                      className={cn(
                        'w-5 h-5 text-gray-500',
                        isActive && 'text-gray-700',
                      )}
                    />
                  )}
                />
              ))}
          </>
        )}
      </div>

      <div className='space-y-1 flex flex-col flex-wrap-grow justify-end mt-auto'>
        <GoogleSidebarNotification />
        <NotificationCenter />

        <SidenavItem
          label='Settings'
          isActive={checkIsActive('settings')}
          onClick={() => router.push('/settings')}
          icon={(isActive) => (
            <Icons.Settings
              className={cn(
                'w-5 h-5 text-gray-500',
                isActive && 'text-gray-700',
              )}
            />
          )}
        />
        <SidenavItem
          label='Sign out'
          isActive={false}
          onClick={handleSignOutClick}
          icon={(isActive) => (
            <LogOut01
              className={cn(
                'w-5 h-5 text-gray-500',
                isActive && 'text-gray-700',
              )}
            />
          )}
        />
      </div>
      <div className='flex h-16' />
    </div>
  );
};
