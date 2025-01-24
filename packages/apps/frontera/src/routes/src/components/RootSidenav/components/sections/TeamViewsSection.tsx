import { observer } from 'mobx-react-lite';
import { TableViewDef } from '@store/TableViewDefs/TableViewDef.dto.ts';

import { cn } from '@ui/utils/cn';
import { Icon } from '@ui/media/Icon';
import { useStore } from '@shared/hooks/useStore';
import { Preferences } from '@shared/components/RootSidenav/hooks';
import { iconNameMap } from '@shared/components/RootSidenav/utils';
import { EditableSideNavItem } from '@shared/components/RootSidenav/components/EditableSidenavItem';

import { CollapsibleSection } from '../CollapsibleSection';

interface TeamViewsSectionSectionProps {
  preferences: Preferences;
  handleItemClick: (data: string) => void;
  togglePreference: (data: keyof Preferences) => void;
  checkIsActive: (
    path: string,
    options?: { preset: string | Array<string> },
  ) => boolean;
}

export const TeamViewsSectionSection = observer(
  ({
    preferences,
    togglePreference,
    handleItemClick,
    checkIsActive,
  }: TeamViewsSectionSectionProps) => {
    const store = useStore();
    const tableViewDefsList = store.tableViewDefs.toArray();
    const lifecycleStagesView = [
      store.tableViewDefs.getById(store.tableViewDefs.targetsPreset ?? ''),
      store.tableViewDefs.getById(store.tableViewDefs.customersPreset ?? ''),
    ].filter(Boolean) as TableViewDef[];

    const teamViewsSectionView =
      tableViewDefsList
        .filter((c) => !c.value.isPreset && c.value.isShared)
        .sort((a, b) => a.value.order - b.value.order) ?? [];

    if (!teamViewsSectionView.length && !lifecycleStagesView.length)
      return null;

    return (
      <CollapsibleSection
        title='Team views'
        isOpen={preferences.isTeamViewsOpen}
        onToggle={() => togglePreference('isTeamViewsOpen')}
      >
        {preferences.isTeamViewsOpen &&
          [...(lifecycleStagesView || []), ...(teamViewsSectionView || [])].map(
            (view) => (
              <EditableSideNavItem
                id={view.value.id}
                key={view.value.id}
                label={view.value.name}
                dataTest={`side-nav-item-${view.value.name}`}
                isActive={checkIsActive('finder', {
                  preset: view.value.id,
                })}
                onClick={() =>
                  handleItemClick(`finder?preset=${view.value.id}`)
                }
                icon={(isActive) => (
                  <Icon
                    name={iconNameMap?.[view.value.icon] ?? 'building-07'}
                    className={cn(
                      'size-4 min-w-4 text-gray-500',
                      isActive && 'text-gray-700',
                    )}
                  />
                )}
              />
            ),
          )}
      </CollapsibleSection>
    );
  },
);
