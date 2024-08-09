import React from 'react';
import { useLocation, useNavigate, useSearchParams } from 'react-router-dom';

import { produce } from 'immer';
import { observer } from 'mobx-react-lite';
import { useLocalStorage } from 'usehooks-ts';
import { useFeatureIsOn } from '@growthbook/growthbook-react';

import { cn } from '@ui/utils/cn';
import { Image } from '@ui/media/Image/Image';
import { Triage } from '@ui/media/icons/Triage';
import { Skeleton } from '@ui/feedback/Skeleton';
import { useStore } from '@shared/hooks/useStore';
import { Bubbles } from '@ui/media/icons/Bubbles';
import { LogOut01 } from '@ui/media/icons/LogOut01';
import { Tag, TagLabel } from '@ui/presentation/Tag';
import { Users01 } from '@ui/media/icons/Users01.tsx';
import { HeartHand } from '@ui/media/icons/HeartHand';
import { Tooltip } from '@ui/overlay/Tooltip/Tooltip';
import { Building07 } from '@ui/media/icons/Building07';
import { CheckHeart } from '@ui/media/icons/CheckHeart';
import { Settings01 } from '@ui/media/icons/Settings01';
import { Briefcase01 } from '@ui/media/icons/Briefcase01';
import { BrokenHeart } from '@ui/media/icons/BrokenHeart';
import { AlertSquare } from '@ui/media/icons/AlertSquare';
import { TableIdType, TableViewType } from '@graphql/types';
import { InvoiceCheck } from '@ui/media/icons/InvoiceCheck';
import { ArrowDropdown } from '@ui/media/icons/ArrowDropdown';
import { CoinsStacked01 } from '@ui/media/icons/CoinsStacked01';
import { InvoiceUpcoming } from '@ui/media/icons/InvoiceUpcoming';
import { ClockFastForward } from '@ui/media/icons/ClockFastForward';
import { SwitchHorizontal01 } from '@ui/media/icons/SwitchHorizontal01.tsx';
import { useOrganizationsMeta } from '@shared/state/OrganizationsMeta.atom';
import { NotificationCenter } from '@shared/components/Notifications/NotificationCenter';
import { EditableSideNavItem } from '@shared/components/RootSidenav/components/EditableSidenavItem.tsx';

import { SidenavItem } from './components/SidenavItem';
import logoCustomerOs from './assets/logo-customeros.png';

const iconMap: Record<
  string,
  (props: React.SVGAttributes<SVGElement>) => JSX.Element
> = {
  InvoiceUpcoming: (props) => <InvoiceUpcoming {...props} />,
  InvoiceCheck: (props) => <InvoiceCheck {...props} />,
  ClockFastForward: (props) => <ClockFastForward {...props} />,
  Briefcase01: (props) => <Briefcase01 {...props} />,
  Building07: (props) => <Building07 {...props} />,
  CheckHeart: (props) => <CheckHeart {...props} />,
  Seed: (props) => <HeartHand {...props} />,
  HeartHand: (props) => <HeartHand {...props} />,
  Triage: (props) => <Triage {...props} />,
  SwitchHorizontal01: (props) => <SwitchHorizontal01 {...props} />,
  BrokenHeart: (props) => <BrokenHeart {...props} />,
};

export const RootSidenav = observer(() => {
  const navigate = useNavigate();
  const { pathname } = useLocation();
  const [searchParams] = useSearchParams();
  const [_, setOrganizationsMeta] = useOrganizationsMeta();
  const showMyViewsItems = useFeatureIsOn('my-views-nav-item');
  const preset = searchParams.get('preset');
  const search = searchParams.get('search');
  const store = useStore();

  const [lastActivePosition, setLastActivePosition] = useLocalStorage(
    `customeros-player-last-position`,
    { root: 'organization' },
  );
  const [lastSearchForPreset, setLastSearchForPreset] = useLocalStorage<{
    [key: string]: string;
  }>(`customeros-last-search-for-preset`, { root: 'root' });

  const [preferences, setPreferences] = useLocalStorage(
    'customeros-preferences',
    {
      isLifecycleViewsOpen: true,
      isMyViewsOpen: true,
      isViewsOpen: true,
      isFavoritesOpen: true,
    },
  );

  const tableViewDefsList = store.tableViewDefs.toArray();

  const myViews = store.demoMode
    ? []
    : tableViewDefsList.filter(
        (c) => c.value.tableType === TableViewType.Renewals && c.value.isPreset,
      ) ?? [];

  const invoicesViews =
    tableViewDefsList.filter(
      (c) => c.value.tableType === TableViewType.Invoices && c.value.isPreset,
    ) ?? [];

  const lifecycleStagesView =
    tableViewDefsList
      .filter(
        (c) =>
          [
            TableIdType.Leads,
            TableIdType.Nurture,
            TableIdType.Customers,
          ].includes(c.value.tableId) && c.value.isPreset,
      )
      .sort((a, b) => a.value.order - b.value.order) ?? [];

  const favoritesView =
    tableViewDefsList
      .filter((c) => !c.value.isPreset)
      .sort((a, b) => a.value.order - b.value.order) ?? [];

  const allOrganizationsView = tableViewDefsList.filter(
    (c) => c.value.tableId === TableIdType.Organizations && c.value.isPreset,
  );

  const allContactsView = tableViewDefsList.find(
    (e) => e.value.tableId === TableIdType.Contacts && e.value.isPreset,
  );

  const churnView = tableViewDefsList.filter(
    (c) => c.value.tableId === TableIdType.Churn && c.value.isPreset,
  );

  const handleItemClick = (path: string) => {
    setLastActivePosition({ ...lastActivePosition, root: path });

    if (preset) {
      setLastSearchForPreset({
        ...lastSearchForPreset,
        [preset]: search ?? '',
      });
    }

    setOrganizationsMeta((prev) =>
      produce(prev, (draft) => {
        draft.getOrganization.pagination.page = 1;
      }),
    );

    navigate(`/${path}`);
  };

  const checkIsActive = (
    path: string,
    options?: { preset: string | Array<string> },
  ) => {
    const _pathName = path.split('?');

    const presetParam = new URLSearchParams(searchParams?.toString()).get(
      'preset',
    );

    if (options?.preset) {
      const isArr = Array.isArray(options.preset);

      if (isArr) {
        return (
          pathname?.startsWith(`/${_pathName}`) &&
          options.preset.includes(presetParam ?? '')
        );
      }

      return (
        pathname?.startsWith(`/${_pathName}`) && presetParam === options.preset
      );
    } else {
      return pathname?.startsWith(`/${_pathName}`) && !presetParam;
    }
  };

  const handleSignOutClick = () => {
    store.session.clearSession();

    if (store.demoMode) {
      window.location.reload();

      return;
    }
    navigate('/auth/signin');
  };
  const showInvoices = store.settings.tenant.value?.billingEnabled;
  const isLoading = store.globalCache?.isLoading;
  const isOwner = store?.globalCache?.value?.isOwner;

  const noOfOrganizationsMovedByICP = store.ui.movedIcpOrganization;

  const allOrganizationsActivePreset = [allOrganizationsView?.[0]?.value?.id];

  return (
    <div className='pb-4 h-full w-12.5 bg-white flex flex-col border-r border-gray-200 overflow-hidden'>
      <div className='px-2 pt-2.5 h-fit mb-2 ml-3 cursor-pointer flex justify-flex-start relative'>
        {!isLoading ? (
          <Image
            width={136}
            height={30}
            alt='CustomerOS'
            className='pointer-events-none transition-opacity-250 ease-in-out h-[30px]  w-auto'
            src={
              store.globalCache.value?.cdnLogoUrl ||
              store.settings.tenant.value?.logoRepositoryFileId ||
              logoCustomerOs
            }
          />
        ) : (
          <Skeleton className='w-full h-8 mr-2' />
        )}
      </div>

      <div className='px-2 pt-2.5 gap-4 overflow-y-auto flex flex-col flex-1'>
        <div className='w-full'>
          <SidenavItem
            label='Customer map'
            dataTest={`side-nav-item-customer-map`}
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

          <div
            className='w-full pt-3 gap-1 flex justify-flex-start pl-3.5 cursor-pointer text-gray-500 hover:text-gray-700 transition-colors'
            onClick={() =>
              setPreferences((prev) => ({
                ...prev,
                isLifecycleViewsOpen: !prev.isLifecycleViewsOpen,
              }))
            }
          >
            <span className='text-sm text-gray-500'>Lifecycle stages</span>

            <ArrowDropdown
              className={cn('w-5 h-5', {
                'transform -rotate-90': !preferences.isLifecycleViewsOpen,
              })}
            />
          </div>

          {preferences.isLifecycleViewsOpen && (
            <>
              {lifecycleStagesView.reduce<JSX.Element[]>((acc, view, index) => {
                const contractsPreset = tableViewDefsList.find(
                  (e) =>
                    e.value.tableId ===
                    TableIdType.ContactsForTargetOrganizations,
                )?.value.id;

                const preset =
                  view.value.tableId === TableIdType.Nurture && contractsPreset
                    ? [view.value.id, contractsPreset]
                    : view.value.id;

                const currentPreset = searchParams?.get('preset');

                const activePreset = store.tableViewDefs
                  ?.toArray()
                  .find((e) => e.value.id === currentPreset)?.value?.id;

                const targetsPreset = tableViewDefsList.find(
                  (e) => e.value.name === 'Targets',
                )?.value.id;

                if (activePreset === targetsPreset) {
                  setTimeout(() => {
                    store.ui.setMovedIcpOrganization(0);
                  }, 2000);
                }

                acc.push(
                  <SidenavItem
                    key={view.value.id}
                    label={view.value.name}
                    dataTest={`side-nav-item-${view.value.name}`}
                    isActive={checkIsActive('finder', {
                      preset,
                    })}
                    onClick={() =>
                      handleItemClick(`finder?preset=${view.value.id}`)
                    }
                    rightElement={
                      noOfOrganizationsMovedByICP > 0 &&
                      view.value.tableId === TableIdType.Nurture ? (
                        <Tag size='sm' variant='solid' colorScheme='gray'>
                          <TagLabel>{noOfOrganizationsMovedByICP}</TagLabel>
                        </Tag>
                      ) : null
                    }
                    icon={(isActive) => {
                      const Icon = iconMap?.[view.value.icon];

                      if (Icon) {
                        return (
                          <Icon
                            className={cn(
                              'w-5 h-5 text-gray-500',
                              isActive && 'text-gray-700',
                            )}
                          />
                        );
                      }

                      return <div className='size-5' />;
                    }}
                  />,
                );

                if (index === 1) {
                  acc.push(
                    <SidenavItem
                      label='Opportunities'
                      key={'kanban-experimental-view'}
                      isActive={checkIsActive('prospects')}
                      dataTest={`side-nav-item-opportunities`}
                      onClick={() => handleItemClick(`prospects`)}
                      icon={(isActive) => {
                        return (
                          <CoinsStacked01
                            className={cn(
                              'w-5 h-5 text-gray-500',
                              isActive && 'text-gray-700',
                            )}
                          />
                        );
                      }}
                    />,
                  );
                }

                return acc;
              }, [])}
            </>
          )}
        </div>

        {favoritesView.length > 0 && (
          <div className='w-full'>
            <div
              onClick={() =>
                setPreferences((prev) => ({
                  ...prev,
                  isFavoritesOpen: !prev.isFavoritesOpen,
                }))
              }
              className={cn(
                'w-full gap-1 flex justify-flex-start pl-3.5 cursor-pointer text-gray-500 hover:text-gray-700 transition-colors',
                {
                  'mb-1': preferences.isFavoritesOpen,
                },
              )}
            >
              <span className='text-sm text-gray-500'>Favorites</span>
              <ArrowDropdown
                className={cn('w-5 h-5', {
                  'transform -rotate-90': !preferences.isFavoritesOpen,
                })}
              />
            </div>

            {preferences.isFavoritesOpen &&
              favoritesView.map((view) => (
                <EditableSideNavItem
                  key={view.value.id}
                  label={view.value.name}
                  dataTest={`side-nav-item-${view.value.name}`}
                  isActive={checkIsActive('finder', {
                    preset: view.value.id,
                  })}
                  onClick={() =>
                    handleItemClick(`finder?preset=${view.value.id}`)
                  }
                  icon={(isActive) => {
                    const Icon = iconMap[view.value.icon];

                    return (
                      <Icon
                        className={cn(
                          'text-md text-gray-500',
                          isActive && 'text-gray-700',
                        )}
                      />
                    );
                  }}
                />
              ))}
          </div>
        )}

        {(isOwner || showMyViewsItems) && (
          <div className='w-full'>
            <div
              onClick={() =>
                setPreferences((prev) => ({
                  ...prev,
                  isMyViewsOpen: !prev.isMyViewsOpen,
                }))
              }
              className={cn(
                'w-full gap-1 flex justify-flex-start pl-3.5 cursor-pointer text-gray-500 hover:text-gray-700 transition-colors',
                {
                  'mb-1': preferences.isMyViewsOpen,
                },
              )}
            >
              <span className='text-sm text-gray-500'>My views</span>
              <ArrowDropdown
                className={cn('w-5 h-5', {
                  'transform -rotate-90': !preferences.isMyViewsOpen,
                })}
              />
            </div>
            {preferences.isMyViewsOpen && (
              <>
                {isOwner &&
                  lifecycleStagesView
                    .filter((o) => o.value.name === 'My portfolio')
                    .map((view) => (
                      <SidenavItem
                        key={view.value.id}
                        label={view.value.name}
                        dataTest={`side-nav-item-${view.value.name}`}
                        isActive={checkIsActive('finder', {
                          preset: view.value.id,
                        })}
                        onClick={() =>
                          handleItemClick(`finder?preset=${view.value.id}`)
                        }
                        icon={(isActive) => {
                          const Icon = iconMap[view.value.icon];

                          return (
                            <Icon
                              className={cn(
                                'w-5 h-5 text-gray-500',
                                isActive && 'text-gray-700',
                              )}
                            />
                          );
                        }}
                      />
                    ))}
                {showMyViewsItems &&
                  myViews.map((view) => (
                    <SidenavItem
                      key={view.value.id}
                      label={view.value.name}
                      dataTest={`side-nav-item-${view.value.name}`}
                      isActive={checkIsActive('renewals', {
                        preset: view.value.id,
                      })}
                      onClick={() =>
                        handleItemClick(`renewals?preset=${view.value.id}`)
                      }
                      icon={(isActive) => {
                        const Icon = iconMap[view.value.icon];

                        return (
                          <Icon
                            className={cn(
                              'w-5 h-5 text-gray-500',
                              isActive && 'text-gray-700',
                            )}
                          />
                        );
                      }}
                    />
                  ))}
                <SidenavItem
                  label={churnView.map((c) => c.value.name).join('')}
                  isActive={checkIsActive('finder', {
                    preset: churnView?.[0]?.value?.id,
                  })}
                  onClick={() =>
                    handleItemClick(
                      `finder?preset=${churnView?.[0]?.value?.id}`,
                    )
                  }
                  icon={(isActive) => (
                    <BrokenHeart
                      className={cn(
                        'w-5 h-5 text-gray-500',
                        isActive && 'text-gray-700',
                      )}
                    />
                  )}
                />
              </>
            )}
          </div>
        )}

        <div
          onClick={() =>
            setPreferences((prev) => ({
              ...prev,
              isViewsOpen: !prev.isViewsOpen,
            }))
          }
          className={cn(
            'w-full gap-1 flex justify-flex-start pl-3.5 cursor-pointer text-gray-500 hover:text-gray-700 transition-colors',
            {
              'mb-1': preferences.isViewsOpen,
            },
          )}
        >
          <span className='text-sm text-gray-500'>Views</span>

          <ArrowDropdown
            className={cn('w-5 h-5', {
              'transform -rotate-90': !preferences.isViewsOpen,
            })}
          />
        </div>

        <div className='w-full'>
          {preferences.isViewsOpen && (
            <>
              <SidenavItem
                label='Organizations'
                dataTest={`side-nav-item-all-orgs`}
                isActive={checkIsActive('finder', {
                  preset: allOrganizationsActivePreset,
                })}
                onClick={() =>
                  handleItemClick(
                    `finder?preset=${allOrganizationsView?.[0]?.value?.id}`,
                  )
                }
                icon={(isActive) => (
                  <Building07
                    className={cn(
                      'w-5 h-5 text-gray-500',
                      isActive && 'text-gray-700',
                    )}
                  />
                )}
              />
              <SidenavItem
                label='Contacts'
                dataTest={`side-nav-item-all-contacts`}
                onClick={() =>
                  handleItemClick(`finder?preset=${allContactsView?.value?.id}`)
                }
                isActive={checkIsActive('finder', {
                  preset: allContactsView?.value?.id ?? '',
                })}
                icon={(isActive) => (
                  <Users01
                    className={cn(
                      'w-5 h-5 text-gray-500',
                      isActive && 'text-gray-700',
                    )}
                  />
                )}
              />

              {showInvoices &&
                invoicesViews.map((view) => (
                  <SidenavItem
                    key={view.value.id}
                    label={`${view.value.name} Invoices`}
                    dataTest={`side-nav-item-${view.value.name}`}
                    isActive={checkIsActive('finder', {
                      preset: view.value.id,
                    })}
                    onClick={() =>
                      handleItemClick(`finder?preset=${view.value.id}`)
                    }
                    icon={(isActive) => {
                      const Icon = iconMap[view.value.icon];

                      return (
                        <Icon
                          className={cn(
                            'w-5 h-5 text-gray-500',
                            isActive && 'text-gray-700',
                          )}
                        />
                      );
                    }}
                  />
                ))}
            </>
          )}
        </div>
      </div>

      <div className='space-y-1 flex flex-col flex-wrap-grow justify-end mt-auto sticky bottom-0 bg-white'>
        {/* <EmailExpiredSidebarNotification /> */}
        <NotificationCenter />

        <SidenavItem
          label='Settings'
          dataTest={`side-nav-item-settings`}
          isActive={checkIsActive('settings')}
          onClick={() => navigate('/settings')}
          icon={(isActive) => (
            <Settings01
              className={cn(
                'w-5 h-5 text-gray-500',
                isActive && 'text-gray-700',
              )}
            />
          )}
          rightElement={
            store.globalCache.value?.inactiveEmailTokens &&
            store.globalCache.value?.inactiveEmailTokens.length > 0 ? (
              <Tooltip
                hasArrow
                className='max-w-[320px]'
                label={
                  'Your conversations and meetings are no longer syncing because access to some of your email accounts has expired'
                }
              >
                <span>
                  <AlertSquare className='text-warning-500' />
                </span>
              </Tooltip>
            ) : null
          }
        />
        <SidenavItem
          label='Sign out'
          isActive={false}
          onClick={handleSignOutClick}
          dataTest={`side-nav-item-sign-out`}
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
    </div>
  );
});
