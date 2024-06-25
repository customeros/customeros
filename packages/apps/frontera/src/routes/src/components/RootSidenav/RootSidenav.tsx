import React from 'react';
import { useNavigate, useLocation, useSearchParams } from 'react-router-dom';

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
import { HeartHand } from '@ui/media/icons/HeartHand';
import { Building07 } from '@ui/media/icons/Building07';
import { CheckHeart } from '@ui/media/icons/CheckHeart';
import { Settings01 } from '@ui/media/icons/Settings01';
import { Briefcase01 } from '@ui/media/icons/Briefcase01';
import { BrokenHeart } from '@ui/media/icons/BrokenHeart';
import { TableIdType, TableViewType } from '@graphql/types';
import { InvoiceCheck } from '@ui/media/icons/InvoiceCheck';
import { ArrowDropdown } from '@ui/media/icons/ArrowDropdown';
import { CoinsStacked01 } from '@ui/media/icons/CoinsStacked01';
import { InvoiceUpcoming } from '@ui/media/icons/InvoiceUpcoming';
import { ClockFastForward } from '@ui/media/icons/ClockFastForward';
import { SwitchHorizontal01 } from '@ui/media/icons/SwitchHorizontal01.tsx';
import { useOrganizationsMeta } from '@shared/state/OrganizationsMeta.atom';
import { NotificationCenter } from '@shared/components/Notifications/NotificationCenter';

import { SidenavItem } from './components/SidenavItem';
import logoCustomerOs from './assets/logo-customeros.png';
import { EmailExpiredSidebarNotification } from './components/EmailExpiredSidebarNotification';

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
  Churn: (props) => <BrokenHeart {...props} />,
};

export const RootSidenav = observer(() => {
  const navigate = useNavigate();
  const { pathname } = useLocation();
  const [searchParams] = useSearchParams();
  const [_, setOrganizationsMeta] = useOrganizationsMeta();
  const showMyViewsItems = useFeatureIsOn('my-views-nav-item');
  const showKanbanView = useFeatureIsOn('prospects');

  const store = useStore();

  const [lastActivePosition, setLastActivePosition] = useLocalStorage(
    `customeros-player-last-position`,
    { root: 'organization' },
  );
  const [preferences, setPreferences] = useLocalStorage(
    'customeros-preferences',
    {
      isInvoicesOpen: true,
      isGrowthOpen: true,
      isAcquisitionOpen: true,
      isMyViewsOpen: true,
    },
  );

  const tableViewDefsList = store.tableViewDefs.toArray();

  const myViews = store.demoMode
    ? []
    : tableViewDefsList.filter(
        (c) => c.value.tableType === TableViewType.Renewals,
      ) ?? [];

  const invoicesViews =
    tableViewDefsList.filter(
      (c) => c.value.tableType === TableViewType.Invoices,
    ) ?? [];
  const growthView =
    tableViewDefsList.filter((c) =>
      [TableIdType.Customers].includes(c.value.tableId),
    ) ?? [];

  const acquisitionView =
    tableViewDefsList
      .filter((c) =>
        [TableIdType.Leads, TableIdType.Nurture].includes(c.value.tableId),
      )
      .sort((a, b) => a.value.order - b.value.order) ?? [];

  const allOrganizationsView = tableViewDefsList.filter(
    (c) => c.value.tableId === TableIdType.Organizations,
  );

  const churnView = tableViewDefsList.filter(
    (c) => c.value.tableId === TableIdType.Churn,
  );

  const handleItemClick = (path: string) => {
    setLastActivePosition({ ...lastActivePosition, root: path });
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

  return (
    <div className='pb-4 h-full w-12.5 bg-white flex flex-col border-r border-gray-200 overflow-hidden'>
      <div className='px-2 pt-2.5 h-fit mb-2 ml-3 cursor-pointer flex justify-flex-start relative'>
        {!isLoading ? (
          <Image
            src={
              store.globalCache.value?.cdnLogoUrl ||
              store.settings.tenant.value?.logoRepositoryFileId ||
              logoCustomerOs
            }
            alt='CustomerOS'
            width={136}
            height={30}
            className='pointer-events-none transition-opacity-250 ease-in-out h-[30px]  w-auto'
          />
        ) : (
          <Skeleton className='w-full h-8 mr-2' />
        )}
      </div>

      <div className='px-2 pt-2.5 overflow-y-auto flex flex-col flex-1'>
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
          <div
            className='w-full pt-3 gap-1 flex justify-flex-start pl-3.5 cursor-pointer text-gray-500 hover:text-gray-700 transition-colors'
            onClick={() =>
              setPreferences((prev) => ({
                ...prev,
                isGrowthOpen: !prev.isGrowthOpen,
              }))
            }
          >
            <span className='text-sm text-gray-500'>Growth</span>
            <ArrowDropdown
              className={cn('w-5 h-5', {
                'transform -rotate-90': !preferences.isGrowthOpen,
              })}
            />
          </div>

          {preferences.isGrowthOpen && (
            <>
              {growthView
                .filter((o) => {
                  if (showKanbanView) {
                    return true;
                  }

                  return (
                    TableIdType.Leads !== o.value.tableId &&
                    TableIdType.Nurture !== o.value.tableId
                  );
                })
                .map((view) => (
                  <SidenavItem
                    key={view.value.id}
                    label={view.value.name}
                    isActive={checkIsActive('finder', {
                      preset: view.value.id,
                    })}
                    onClick={() =>
                      handleItemClick(`finder?preset=${view.value.id}`)
                    }
                    icon={(isActive) => {
                      const Icon = iconMap?.[view.value.icon];

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

          <div
            className='w-full pt-3 gap-1 flex justify-flex-start pl-3.5 cursor-pointer text-gray-500 hover:text-gray-700 transition-colors'
            onClick={() =>
              setPreferences((prev) => ({
                ...prev,
                isAcquisitionOpen: !prev.isAcquisitionOpen,
              }))
            }
          >
            <span className='text-sm text-gray-500'>Acquisition</span>

            <ArrowDropdown
              className={cn('w-5 h-5', {
                'transform -rotate-90': !preferences.isAcquisitionOpen,
              })}
            />
          </div>

          {preferences.isAcquisitionOpen && (
            <>
              {acquisitionView
                .filter((o) => {
                  if (showKanbanView) {
                    return true;
                  }

                  return (
                    TableIdType.Leads !== o.value.tableId &&
                    TableIdType.Nurture !== o.value.tableId
                  );
                })
                .reduce<JSX.Element[]>((acc, view, index) => {
                  const contractsPreset = tableViewDefsList.find(
                    (e) => e.value.tableType === TableViewType.Contacts,
                  )?.value.id;
                  const preset =
                    view.value.tableId === 'NURTURE' && contractsPreset
                      ? [view.value.id, contractsPreset]
                      : view.value.id;

                  acc.push(
                    <SidenavItem
                      key={view.value.id}
                      label={view.value.name}
                      isActive={checkIsActive('finder', {
                        preset,
                      })}
                      onClick={() =>
                        handleItemClick(`finder?preset=${view.value.id}`)
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
                  if (showKanbanView && index === 1) {
                    acc.push(
                      <SidenavItem
                        key={'kanban-experimental-view'}
                        label='Opportunities'
                        isActive={checkIsActive('prospects')}
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
                <span className='text-sm text-gray-500'>Invoices</span>
                <ArrowDropdown
                  className={cn('w-5 h-5', {
                    'transform -rotate-90': !preferences.isInvoicesOpen,
                  })}
                />
              </div>

              {preferences.isInvoicesOpen && (
                <>
                  {invoicesViews.map((view) => (
                    <SidenavItem
                      key={view.value.id}
                      label={view.value.name}
                      isActive={checkIsActive('invoices', {
                        preset: view.value.id,
                      })}
                      onClick={() =>
                        handleItemClick(`invoices?preset=${view.value.id}`)
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
            </>
          )}
        </div>

        <div className='space-y-1 w-full'>
          {(isOwner || showMyViewsItems) && (
            <div
              className='w-full gap-1 flex justify-flex-start pl-3.5 cursor-pointer text-gray-500 hover:text-gray-700 transition-colors'
              onClick={() =>
                setPreferences((prev) => ({
                  ...prev,
                  isMyViewsOpen: !prev.isMyViewsOpen,
                }))
              }
            >
              <span className='text-sm text-gray-500'>My views</span>
              <ArrowDropdown
                className={cn('w-5 h-5', {
                  'transform -rotate-90': !preferences.isMyViewsOpen,
                })}
              />
            </div>
          )}

          {preferences.isMyViewsOpen && (
            <>
              {isOwner &&
                acquisitionView
                  .filter((o) => o.value.name === 'My portfolio')
                  .map((view) => (
                    <SidenavItem
                      key={view.value.id}
                      label={view.value.name}
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
                  handleItemClick(`finder?preset=${churnView?.[0]?.value?.id}`)
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

        <div className='space-y-1 w-full mt-2'>
          <SidenavItem
            label='All orgs'
            isActive={checkIsActive('finder', {
              preset: allOrganizationsView?.[0]?.value?.id,
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
        </div>
      </div>

      <div className='space-y-1 flex flex-col flex-wrap-grow justify-end mt-auto sticky bottom-0 bg-white'>
        <EmailExpiredSidebarNotification />
        <NotificationCenter />

        <SidenavItem
          label='Settings'
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
    </div>
  );
});
