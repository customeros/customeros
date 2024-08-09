import React from 'react';

import { observer } from 'mobx-react-lite';
import { useFeatureIsOn } from '@growthbook/growthbook-react';

import { cn } from '@ui/utils/cn.ts';
import { useStore } from '@shared/hooks/useStore';
import { BrokenHeart } from '@ui/media/icons/BrokenHeart';
import { TableIdType, TableViewType } from '@graphql/types';
import { iconMap } from '@shared/components/RootSidenav/utils';
import { Preferences } from '@shared/components/RootSidenav/hooks';

import { SidenavItem } from '../SidenavItem';
import { CollapsibleSection } from '../CollapsibleSection';

interface MyViewsSectionProps {
  preferences: Preferences;
  handleItemClick: (data: string) => void;
  togglePreference: (data: keyof Preferences) => void;
  checkIsActive: (
    path: string,
    options?: { preset: string | Array<string> },
  ) => boolean;
}

export const MyViewsSection = observer(
  ({
    preferences,
    togglePreference,
    handleItemClick,
    checkIsActive,
  }: MyViewsSectionProps) => {
    const store = useStore();
    const showMyViewsItems = useFeatureIsOn('my-views-nav-item');
    const isOwner = store?.globalCache?.value?.isOwner;

    if (!showMyViewsItems) return null;

    const tableViewDefsList = store.tableViewDefs.toArray();

    const myViews = store.demoMode
      ? []
      : tableViewDefsList.filter(
          (c) =>
            c.value.tableType === TableViewType.Renewals && c.value.isPreset,
        ) ?? [];
    const churnView = tableViewDefsList.filter(
      (c) => c.value.tableId === TableIdType.Churn && c.value.isPreset,
    );

    return (
      <CollapsibleSection
        title='My views'
        isOpen={preferences.isMyViewsOpen}
        onToggle={() => togglePreference('isMyViewsOpen')}
      >
        <div className='w-full'>
          {preferences.isMyViewsOpen && (
            <>
              {isOwner &&
                tableViewDefsList
                  .filter(
                    (o) => o.value.name === 'My portfolio' && o.value.isPreset,
                  )
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
      </CollapsibleSection>
    );
  },
);
