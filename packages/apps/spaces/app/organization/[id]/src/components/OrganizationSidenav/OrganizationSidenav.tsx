'use client';
import React from 'react';
import { useRouter, useParams, useSearchParams } from 'next/navigation';

import { useLocalStorage } from 'usehooks-ts';

import { Icons } from '@ui/media/Icon';
import { IconButton } from '@ui/form/IconButton';
import { Ticket02 } from '@ui/media/icons/Ticket02';
import { Trophy01 } from '@ui/media/icons/Trophy01';
import { Tooltip } from '@ui/overlay/Tooltip/Tooltip2';
import { getGraphQLClient } from '@shared/util/getGraphQLClient';
import { SidenavItem } from '@shared/components/RootSidenav/components/SidenavItem';
import { useOrganizationQuery } from '@organization/src/graphql/organization.generated';
import { NotificationCenter } from '@shared/components/Notifications/NotificationCenter';

export const OrganizationSidenav = () => {
  const router = useRouter();
  const params = useParams();
  const searchParams = useSearchParams();

  const [lastActivePosition, setLastActivePosition] = useLocalStorage(
    `customeros-player-last-position`,
    { [params?.id as string]: 'tab=about' },
  );

  const graphqlClient = getGraphQLClient();
  const { data } = useOrganizationQuery(graphqlClient, {
    id: params?.id as string,
  });
  const parentOrg = data?.organization?.subsidiaryOf?.[0]?.organization;

  const checkIsActive = (tab: string) => searchParams?.get('tab') === tab;

  const handleItemClick = (tab: string) => () => {
    const urlSearchParams = new URLSearchParams(searchParams?.toString());
    urlSearchParams.set('tab', tab);

    setLastActivePosition({
      ...lastActivePosition,
      [params?.id as string]: urlSearchParams.toString(),
    });
    router.push(`?${urlSearchParams}`);
  };

  return (
    <div className='px-2 py-4 h-full w-200 flex flex-col grid-area-sidebar bg-white relative border-r border-gray-200'>
      <div className='flex gap-2 items-center mb-4'>
        <IconButton
          size='xs'
          variant='ghost'
          aria-label='Go back'
          onClick={() => {
            router.push(`/${lastActivePosition?.root || 'organization'}`);
          }}
          icon={<Icons.ArrowNarrowLeft color='gray.700' boxSize='6' />}
        />

        <div className='flex flex-col'>
          {parentOrg && (
            <a
              className='text-xs truncate'
              href={`/organization/${parentOrg.id}?tab=about`}
            >
              {parentOrg.name}
            </a>
          )}
          <Tooltip
            // position={'bottom'}
            label={data?.organization?.name ?? ''}
            // className='bg-gray-700 text-white'
          >
            <span className='max-w-150px text-lg font-semibold text-gray-700 truncate whitespace-nowrap '>
              {data?.organization?.name || 'Organization'}
            </span>
          </Tooltip>
        </div>
      </div>
      <div className='space-y-1 w-full'>
        <SidenavItem
          label='About'
          isActive={checkIsActive('about') || !searchParams?.get('tab')}
          onClick={handleItemClick('about')}
          icon={<Icons.InfoSquare className='w-5 h-5' />}
        />
        <SidenavItem
          label='People'
          isActive={checkIsActive('people')}
          onClick={handleItemClick('people')}
          icon={<Icons.Users2 className='w-5 h-5' />}
        />
        <SidenavItem
          label='Account'
          isActive={checkIsActive('account') || checkIsActive('invoices')}
          onClick={handleItemClick('account')}
          icon={<Icons.ActivityHeart className='w-5 h-5' />}
        />
        <SidenavItem
          label='Success'
          isActive={checkIsActive('success')}
          onClick={handleItemClick('success')}
          icon={<Trophy01 className='w-5 h-5' />}
        />
        <SidenavItem
          label='Issues'
          isActive={checkIsActive('issues')}
          onClick={handleItemClick('issues')}
          icon={<Ticket02 className='w-5 h-5' />}
        />
      </div>
      <div className='flex flex-col flex-grow justify-end'>
        <NotificationCenter />
      </div>
      <div className='flex h-16' />
    </div>
  );
};
