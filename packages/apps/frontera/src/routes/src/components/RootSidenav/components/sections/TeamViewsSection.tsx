import React from 'react';

import { observer } from 'mobx-react-lite';

import { cn } from '@ui/utils/cn.ts';
import { useStore } from '@shared/hooks/useStore';
import { iconMap } from '@shared/components/RootSidenav/utils';
import { Preferences } from '@shared/components/RootSidenav/hooks';
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

    const teamViewsSectionView =
      tableViewDefsList
        .filter((c) => !c.value.isPreset && c.value.isShared)
        .sort((a, b) => a.value.order - b.value.order) ?? [];

    if (!teamViewsSectionView.length) return null;

    return (
      <CollapsibleSection
        title='Team views'
        isOpen={preferences.isTeamViewsOpen}
        onToggle={() => togglePreference('isTeamViewsOpen')}
      >
        {preferences.isTeamViewsOpen &&
          teamViewsSectionView.map((view) => (
            <EditableSideNavItem
              id={view.value.id}
              key={view.value.id}
              label={view.value.name}
              dataTest={`side-nav-item-${view.value.name}`}
              onClick={() => handleItemClick(`finder?preset=${view.value.id}`)}
              isActive={checkIsActive('finder', {
                preset: view.value.id,
              })}
              icon={(isActive) => {
                const Icon = iconMap?.[view.value.icon];

                if (!Icon) return <div />;

                return (
                  <Icon
                    className={cn(
                      'size-4 min-w-4 text-gray-500',
                      isActive && 'text-gray-700',
                    )}
                  />
                );
              }}
            />
          ))}
      </CollapsibleSection>
    );
  },
);
