import React, { FC } from 'react';
import { useRecoilState } from 'recoil';
import { organizationDetailsEdit } from '../../../../state';
import { useRemoveOrganizationSubsidiary } from '@spaces/hooks/useOrganizationSubsidiaries';
import { DeleteIconButton } from '@spaces/atoms/icon-button';
import styles from './organization-subsidiaries.module.scss';
import { AnimatePresence, motion } from 'framer-motion';

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

  if (subsidiaries.length === 0 && !isEditMode) {
    return (
      <div className={styles.subsidiary}>This company has no branches</div>
    );
  }

  return (
    <ul>
      <AnimatePresence initial={false}>
        {subsidiaries.map((e) => (
          <motion.li
            key={e.organization.id}
            initial={{ opacity: 0 }}
            className={styles.subsidiary}
            animate={{ opacity: 1 }}
            exit={{ opacity: 0 }}
          >
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
          </motion.li>
        ))}
      </AnimatePresence>
    </ul>
  );
};
