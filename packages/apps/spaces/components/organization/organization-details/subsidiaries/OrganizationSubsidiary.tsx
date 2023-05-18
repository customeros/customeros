import React, { FC } from 'react';
import { useRecoilState } from 'recoil';
import { organizationDetailsEdit } from '../../../../state';
import { useRemoveOrganizationSubsidiary } from '@spaces/hooks/useOrganizationSubsidiaries';
import { DeleteIconButton } from '@spaces/atoms/icon-button';
import styles from './organization-subsidiaries.module.scss';

interface OrganizationSubsidiariesProps {
  subsidiaries: Array<any>;
  id: string;
}
export const OrganizationSubsidiary: FC<OrganizationSubsidiariesProps> = ({
  subsidiaries,
  id,
}) => {
  const [{ isEditMode }] = useRecoilState(organizationDetailsEdit);

  const { onRemoveOrganizationSubsidiary } = useRemoveOrganizationSubsidiary({
    organizationId: id,
  });

  return (
    <>
      {subsidiaries.map((e) => (
        <div key={e.organization.id} className={styles.subsidiary}>
          {isEditMode && (
            <DeleteIconButton
              onDelete={() =>
                onRemoveOrganizationSubsidiary({
                  subsidiaryId: e.organization.id,
                })
              }
            />
          )}

          <span style={{ marginLeft: isEditMode ? '8px' : '0' }}>
            {e.organization.name || 'Unnamed'}
          </span>
        </div>
      ))}
    </>
  );
};
