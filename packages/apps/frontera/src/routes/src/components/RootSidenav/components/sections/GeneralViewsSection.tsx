import React from 'react';

import { observer } from 'mobx-react-lite';
import { TableViewDefStore } from '@store/TableViewDefs/TableViewDef.store.ts';

import { cn } from '@ui/utils/cn.ts';
import { TableIdType } from '@graphql/types';
import { useStore } from '@shared/hooks/useStore';
import { Users01 } from '@ui/media/icons/Users01';
import { Invoice } from '@ui/media/icons/Invoice.tsx';
import { Building07 } from '@ui/media/icons/Building07';
import { Shuffle01 } from '@ui/media/icons/Shuffle01.tsx';
import { Signature } from '@ui/media/icons/Signature.tsx';
import { Preferences } from '@shared/components/RootSidenav/hooks';
import { CoinsStacked01 } from '@ui/media/icons/CoinsStacked01.tsx';
import { RootSidenavItem } from '@shared/components/RootSidenav/components/RootSidenavItem';

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

    const allContactsView = store.tableViewDefs.getById(
      store.tableViewDefs.contactsPreset ?? '',
    );
    const contractsView = store.tableViewDefs.getById(
      store.tableViewDefs.contractsPreset ?? '',
    );
    const opportunitiesView = store.tableViewDefs.getById(
      store.tableViewDefs.opportunitiesTablePreset ?? '',
    );
    const flowSequencesView = store.tableViewDefs.getById(
      store.tableViewDefs.flowsPreset ?? '',
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
              id={allOrganizationsView?.[0]?.value?.id}
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
                    'size-4 min-w-4 text-gray-500',
                    isActive && 'text-gray-700',
                  )}
                />
              )}
            />
            <RootSidenavItem
              label='Contacts'
              id={allContactsView?.value?.id}
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
                    'size-4 min-w-4 text-gray-500',
                    isActive && 'text-gray-700',
                  )}
                />
              )}
            />
            <RootSidenavItem
              label='Opportunities'
              id={opportunitiesView?.value?.id}
              dataTest={`side-nav-item-opportunities`}
              onClick={() =>
                handleItemClick(`finder?preset=${opportunitiesView?.value?.id}`)
              }
              isActive={checkIsActive('finder', {
                preset: opportunitiesView?.value?.id ?? '',
              })}
              icon={(isActive) => (
                <CoinsStacked01
                  className={cn(
                    'size-4 min-w-4 text-gray-500',
                    isActive && 'text-gray-700',
                  )}
                />
              )}
            />
            {showInvoices && upcomingInvoices && (
              <RootSidenavItem
                label='Invoices'
                id={upcomingInvoices.value.id}
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
                        'size-4 min-w-4 text-gray-500',
                        isActive && 'text-gray-700',
                      )}
                    />
                  );
                }}
              />
            )}
            <RootSidenavItem
              label='Contracts'
              id={contractsView?.value?.id}
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
                    'size-4 min-w-4 text-gray-500',
                    isActive && 'text-gray-700',
                  )}
                />
              )}
            />
            <SidenavItem
              label='Flows'
              dataTest={`side-nav-item-all-flows`}
              onClick={() =>
                handleItemClick(`finder?preset=${flowSequencesView?.value?.id}`)
              }
              isActive={checkIsActive('finder', {
                preset: flowSequencesView?.value?.id ?? '',
              })}
              icon={(isActive) => (
                <Shuffle01
                  className={cn(
                    'size-4 min-w-4 text-gray-500',
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
