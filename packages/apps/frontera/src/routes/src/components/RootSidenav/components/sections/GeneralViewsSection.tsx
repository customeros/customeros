import React from 'react';

import { observer } from 'mobx-react-lite';

import { cn } from '@ui/utils/cn.ts';
import { useStore } from '@shared/hooks/useStore';
import { Users01 } from '@ui/media/icons/Users01';
import { Invoice } from '@ui/media/icons/Invoice.tsx';
import { Building07 } from '@ui/media/icons/Building07';
import { TableIdType, TableViewType } from '@graphql/types';
import { Preferences } from '@shared/components/RootSidenav/hooks';

import { SidenavItem } from '../SidenavItem';
import { CollapsibleSection } from '../CollapsibleSection';

interface GeneralViewsSectionProps {
  preferences: Preferences;
  handleItemClick: (data: string) => void;
  togglePreference: (data: keyof Preferences) => void;
  checkIsActive: (
    path: string,
    options?: { preset: string | Array<string> },
  ) => boolean;
}

export const GeneralViewsSection = observer(
  ({
    preferences,
    togglePreference,
    handleItemClick,
    checkIsActive,
  }: GeneralViewsSectionProps) => {
    const store = useStore();
    const tableViewDefsList = store.tableViewDefs.toArray();
    const allOrganizationsView = tableViewDefsList.filter(
      (c) => c.value.tableId === TableIdType.Organizations && c.value.isPreset,
    );

    const allContactsView = tableViewDefsList.find(
      (e) => e.value.tableId === TableIdType.Contacts && e.value.isPreset,
    );
    const invoicesViews =
      tableViewDefsList.filter(
        (c) => c.value.tableType === TableViewType.Invoices && c.value.isPreset,
      ) ?? [];

    const upcomingInvoices = invoicesViews.find(
      (e) => e.value.tableId === TableIdType.UpcomingInvoices,
    );
    const allOrganizationsActivePreset = [allOrganizationsView?.[0]?.value?.id];
    const showInvoices = store.settings.tenant.value?.billingEnabled;

    return (
      <CollapsibleSection
        title='Views'
        isOpen={preferences.isViewsOpen}
        onToggle={() => togglePreference('isViewsOpen')}
      >
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

            {showInvoices && upcomingInvoices && (
              <SidenavItem
                label='Invoices'
                key={upcomingInvoices.value.id}
                dataTest={`side-nav-item-${upcomingInvoices.value.name}`}
                onClick={() =>
                  handleItemClick(`finder?preset=${upcomingInvoices.value.id}`)
                }
                isActive={checkIsActive('finder', {
                  preset: invoicesViews.map((e) => e.value.id),
                })}
                icon={(isActive) => {
                  return (
                    <Invoice
                      className={cn(
                        'w-5 h-5 text-gray-500',
                        isActive && 'text-gray-700',
                      )}
                    />
                  );
                }}
              />
            )}
          </>
        )}
      </CollapsibleSection>
    );
  },
);
