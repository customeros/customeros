import React from 'react';
import { useMergeOrganizations } from '@spaces/hooks/useOrganization';
import { ActionColumn } from '@spaces/finder/finder-table';
import { useSetRecoilState } from 'recoil';
import { tableMode } from '@spaces/finder/state';
import { useRouter } from 'next/router';

export const OrganizationActionColumn: React.FC = () => {
  const { onMergeOrganizations } = useMergeOrganizations();
  const setTableMode = useSetRecoilState(tableMode);
  const { push } = useRouter();

  return (
    <ActionColumn
      onMerge={({ primaryId, mergeIds }) =>
        onMergeOrganizations({
          primaryOrganizationId: primaryId,
          mergedOrganizationIds: mergeIds,
        })
      }
      actions={[
        {
          label: 'Add organization',
          command() {
            push('/organization/new');
          },
        },
        {
          label: 'Merge organizations',
          command() {
            return setTableMode('MERGE');
          },
        },
      ]}
    />
  );
};
