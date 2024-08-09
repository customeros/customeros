import React from 'react';
import { useSearchParams } from 'react-router-dom';

import { observer } from 'mobx-react-lite';

import { cn } from '@ui/utils/cn.ts';
import { TableIdType } from '@graphql/types';
import { useStore } from '@shared/hooks/useStore';
import { Tag, TagLabel } from '@ui/presentation/Tag';
import { iconMap } from '@shared/components/RootSidenav/utils';
import { CoinsStacked01 } from '@ui/media/icons/CoinsStacked01';
import { Preferences } from '@shared/components/RootSidenav/hooks';

import { SidenavItem } from '../SidenavItem';
import { CollapsibleSection } from '../CollapsibleSection';

interface LifecycleStagesSectionProps {
  preferences: Preferences;
  handleItemClick: (data: string) => void;
  togglePreference: (data: keyof Preferences) => void;
  checkIsActive: (
    path: string,
    options?: { preset: string | Array<string> },
  ) => boolean;
}

export const LifecycleStagesSection = observer(
  ({
    preferences,
    togglePreference,
    handleItemClick,
    checkIsActive,
  }: LifecycleStagesSectionProps) => {
    const store = useStore();
    const tableViewDefsList = store.tableViewDefs.toArray();
    const [searchParams] = useSearchParams();

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

    const noOfOrganizationsMovedByICP = store.ui.movedIcpOrganization;

    return (
      <CollapsibleSection
        title='Lifecycle stages'
        isOpen={preferences.isLifecycleViewsOpen}
        onToggle={() => togglePreference('isLifecycleViewsOpen')}
      >
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
      </CollapsibleSection>
    );
  },
);
