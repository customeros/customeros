import React from 'react';
import { useSearchParams } from 'react-router-dom';

import { observer } from 'mobx-react-lite';
import { TableViewDefStore } from '@store/TableViewDefs/TableViewDef.store.ts';

import { cn } from '@ui/utils/cn.ts';
import { TableIdType } from '@graphql/types';
import { useStore } from '@shared/hooks/useStore';
import { Tag, TagLabel } from '@ui/presentation/Tag';
import { iconMap } from '@shared/components/RootSidenav/utils';
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

    const noOfOrganizationsMovedByICP = store.ui.movedIcpOrganization;

    const lifecycleStagesView: TableViewDefStore[] = [
      store.tableViewDefs.getById(store.tableViewDefs.targetsPreset ?? ''),
      store.tableViewDefs.getById(
        store.tableViewDefs.opportunitiesPreset ?? '',
      ),
      store.tableViewDefs.getById(store.tableViewDefs.defaultPreset ?? ''),
    ].filter((e): e is TableViewDefStore => e !== undefined);

    const contractsPreset = tableViewDefsList.find(
      (e) => e.value.tableId === TableIdType.ContactsForTargetOrganizations,
    )?.value.id;

    const currentPreset = searchParams?.get('preset');
    const activePreset = tableViewDefsList.find(
      (e) => e.value.id === currentPreset,
    )?.value?.id;
    const targetsPreset = tableViewDefsList.find(
      (e) => e.value.name === 'Targets',
    )?.value.id;

    if (activePreset === targetsPreset) {
      setTimeout(() => {
        store.ui.setMovedIcpOrganization(0);
      }, 2000);
    }

    return (
      <CollapsibleSection
        title='Lifecycle stages'
        isOpen={preferences.isLifecycleViewsOpen}
        onToggle={() => togglePreference('isLifecycleViewsOpen')}
      >
        {preferences.isLifecycleViewsOpen && (
          <>
            {lifecycleStagesView.map((view) => {
              const preset =
                view.value.tableId === TableIdType.Nurture && contractsPreset
                  ? [view.value.id, contractsPreset]
                  : view.value.id;

              return (
                <SidenavItem
                  key={view.value.id}
                  label={view.value.name}
                  dataTest={`side-nav-item-${view.value.name}`}
                  isActive={
                    view.value.tableId === TableIdType.Opportunities
                      ? checkIsActive('prospects')
                      : checkIsActive('finder', { preset })
                  }
                  onClick={() =>
                    view.value.tableId === TableIdType.Opportunities
                      ? handleItemClick(`prospects`)
                      : handleItemClick(`finder?preset=${view.value.id}`)
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

                    return Icon ? (
                      <Icon
                        className={cn(
                          'w-5 h-5 text-gray-500',
                          isActive && 'text-gray-700',
                        )}
                      />
                    ) : (
                      <div className='size-5' />
                    );
                  }}
                />
              );
            })}
          </>
        )}
      </CollapsibleSection>
    );
  },
);
