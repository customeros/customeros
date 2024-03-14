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
import { Receipt } from '@ui/media/icons/Receipt';
import { Bubbles } from '@ui/media/icons/Bubbles';
import { LogOut01 } from '@ui/media/icons/LogOut01';
import { mockedTableDefs } from '@shared/util/tableDefs.mock';
import { getGraphQLClient } from '@shared/util/getGraphQLClient';
import { ClockFastForward } from '@ui/media/icons/ClockFastForward';
import { useOrganizationsMeta } from '@shared/state/OrganizationsMeta.atom';
import { useGlobalCacheQuery } from '@shared/graphql/global_Cache.generated';
import { useTableViewDefsQuery } from '@shared/graphql/tableViewDefs.generated';
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

  const { data: tenantSettingsData } = useTenantSettingsQuery(client);
  const { data: totalInvoices } = useGetAllInvoicesCountQuery(client);

  const { data: tableViewDefsData } = useTableViewDefsQuery(
    client,
    {
      pagination: { limit: 100, page: 1 },
    },
    {
      enabled: false,
      placeholderData: { tableViewDefs: { content: mockedTableDefs } },
    },
  );
  const { data, isLoading } = useGlobalCacheQuery(client);
  const globalCache = data?.global_Cache;
  const myViews = tableViewDefsData?.tableViewDefs?.content ?? [];

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
    ].forEach((path) => {
      router.prefetch(path);
    });
  }, []);

  const cdnLogoUrl = data?.global_Cache?.cdnLogoUrl;

  return (
    <div className='px-2 pt-2.5 pb-4 h-full w-12.5 bg-white flex flex-col border-r border-gray-200'>
      <div
        className='mb-2 ml-3 cursor-pointer flex justify-flex-start overflow-hidden relative'
        tabIndex={0}
        role='button'
      >
        {!isLoading ? (
          <Image
            src={cdnLogoUrl ?? logoCustomerOs}
            alt='CustomerOS'
            width={136}
            height={30}
            className='pointer-events-none transition-opacity-250 ease-in-out'
          />
        ) : (
          <Skeleton className='w-full h-8 mr-2' />
        )}
      </div>

      <div className='space-y-2 w-full mb-4'>
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
          <SidenavItem
            label='Invoices'
            isActive={checkIsActive('invoices')}
            onClick={() => handleItemClick('invoices')}
            icon={(isActive) => (
              <Receipt
                className={cn(
                  'w-5 h-5 text-gray-500',
                  isActive && 'text-gray-700',
                )}
              />
            )}
          />
        )}
      </div>

      <div className='space-y-2 w-full'>
        {(globalCache?.isOwner || showMyViewsItems) && (
          <div className='w-full flex justify-flex-start pl-3.5'>
            <span className='text-gray-500 text-sm'>My views</span>
          </div>
        )}

        {globalCache?.isOwner && (
          <SidenavItem
            label='My portfolio'
            isActive={checkIsActive('organizations', { preset: 'portfolio' })}
            onClick={() => handleItemClick('organizations?preset=portfolio')}
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
      </div>

      <div className='space-y-1 flex flex-col flex-wrap-grow justify-end mt-auto'>
        <NotificationCenter />
        <GoogleSidebarNotification />

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
