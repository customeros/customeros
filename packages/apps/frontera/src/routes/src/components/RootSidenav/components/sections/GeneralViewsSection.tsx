import React from 'react';

import { observer } from 'mobx-react-lite';
import { TableViewDefStore } from '@store/TableViewDefs/TableViewDef.store.ts';

import { cn } from '@ui/utils/cn.ts';
import { TableIdType } from '@graphql/types';
import { useStore } from '@shared/hooks/useStore';
import { Users01 } from '@ui/media/icons/Users01';
import { Invoice } from '@ui/media/icons/Invoice.tsx';
import { Building07 } from '@ui/media/icons/Building07';
import { Signature } from '@ui/media/icons/Signature.tsx';
import { Preferences } from '@shared/components/RootSidenav/hooks';
import { RootSidenavItem } from '@shared/components/RootSidenav/components/RootSidenavItem';

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

    const allContactsView = store.tableViewDefs.getById(
      store.tableViewDefs.contactsPreset ?? '',
    );
    const contractsView = store.tableViewDefs.getById(
      store.tableViewDefs.contractsPreset ?? '',
    );

    const invoicesViews = [
      store.tableViewDefs.getById(
        store.tableViewDefs.upcomingInvoicesPreset ?? '',
      ),
      store.tableViewDefs.getById(store.tableViewDefs.pastInvoicesPreset ?? ''),
    ].filter((e): e is TableViewDefStore => e !== undefined);

    const upcomingInvoices = invoicesViews[0];
    const allOrganizationsActivePreset = [allOrganizationsView?.[0]?.value?.id];
    const showInvoices = store.settings.tenant.value?.billingEnabled;

    return (
      <CollapsibleSection
        title='Records'
        isOpen={preferences.isViewsOpen}
        onToggle={() => togglePreference('isViewsOpen')}
      >
        {preferences.isViewsOpen && (
          <>
            <RootSidenavItem
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
            <RootSidenavItem
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
              <RootSidenavItem
                label='Invoices'
                key={upcomingInvoices.value.id}
                dataTest={`side-nav-item-${upcomingInvoices.value.name}`}
                onClick={() =>
                  handleItemClick(`finder?preset=${upcomingInvoices.value.id}`)
                }
                isActive={checkIsActive('finder', {
                  preset: invoicesViews.map((e) => e?.value?.id),
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
            <RootSidenavItem
              label='Contracts'
              dataTest={`side-nav-item-all-contracts`}
              onClick={() =>
                handleItemClick(`finder?preset=${contractsView?.value?.id}`)
              }
              isActive={checkIsActive('finder', {
                preset: contractsView?.value?.id ?? '',
              })}
              icon={(isActive) => (
                <Signature
                  className={cn(
                    'w-5 h-5 text-gray-500',
                    isActive && 'text-gray-700',
                  )}
                />
              )}
            />
          </>
        )}
      </CollapsibleSection>
    );
  },
);
