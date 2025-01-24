import { useSearchParams } from 'react-router-dom';

import { observer } from 'mobx-react-lite';
import { TableViewDef } from '@store/TableViewDefs/TableViewDef.dto';

import { cn } from '@ui/utils/cn';
import { Icon } from '@ui/media/Icon';
import { TableIdType } from '@graphql/types';
import { useStore } from '@shared/hooks/useStore';
import { Tag, TagLabel } from '@ui/presentation/Tag';
import { Preferences } from '@shared/components/RootSidenav/hooks';
import { iconNameMap } from '@shared/components/RootSidenav/utils';
import { SidenavItem } from '@shared/components/RootSidenav/components/SidenavItem';
import { RootSidenavItem } from '@shared/components/RootSidenav/components/RootSidenavItem';

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

    const lifecycleStagesView = [
      store.tableViewDefs.getById(store.tableViewDefs.targetsPreset ?? ''),
      store.tableViewDefs.getById(
        store.tableViewDefs.opportunitiesPreset ?? '',
      ),
      store.tableViewDefs.getById(store.tableViewDefs.defaultPreset ?? ''),
    ].filter(Boolean) as TableViewDef[];

    const contractsPreset = tableViewDefsList.find(
      (e) => e.value.tableId === TableIdType.ContactsForTargetOrganizations,
    )?.id;

    const currentPreset = searchParams?.get('preset');
    const activePreset = tableViewDefsList.find(
      (e) => e.value.id === currentPreset,
    )?.id;
    const targetsPreset = tableViewDefsList.find(
      (e) => e.value.name === 'Targets',
    )?.id;

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
                view.value.tableId === TableIdType.Targets && contractsPreset
                  ? [view.value.id, contractsPreset]
                  : view.value.id;

              if (view.value.tableId === TableIdType.Opportunities) {
                return (
                  <SidenavItem
                    key={view.value.id}
                    label={view.value.name}
                    isActive={checkIsActive('prospects')}
                    onClick={() => handleItemClick(`prospects`)}
                    dataTest={`side-nav-item-${view.value.name}`}
                    icon={(isActive) => (
                      <Icon
                        name={iconNameMap?.[view.value.icon] ?? 'building-07'}
                        className={cn(
                          'min-w-4 text-gray-500',
                          isActive && 'text-gray-700',
                        )}
                      />
                    )}
                  />
                );
              }

              return (
                <RootSidenavItem
                  id={view.value.id}
                  key={view.value.id}
                  label={view.value.name}
                  dataTest={`side-nav-item-${view.value.name}`}
                  isActive={checkIsActive('finder', { preset })}
                  onClick={() =>
                    handleItemClick(`finder?preset=${view.value.id}`)
                  }
                  icon={(isActive) => (
                    <Icon
                      name={iconNameMap?.[view.value.icon] ?? 'building-07'}
                      className={cn(
                        'min-w-4 text-gray-500',
                        isActive && 'text-gray-700',
                      )}
                    />
                  )}
                  rightElement={
                    noOfOrganizationsMovedByICP > 0 &&
                    view.value.tableId === TableIdType.Targets ? (
                      <Tag size='sm' variant='solid' colorScheme='gray'>
                        <TagLabel>{noOfOrganizationsMovedByICP}</TagLabel>
                      </Tag>
                    ) : null
                  }
                />
              );
            })}
          </>
        )}
      </CollapsibleSection>
    );
  },
);
