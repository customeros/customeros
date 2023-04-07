import React, { useEffect, useState } from 'react';
import {
  DeleteConfirmationDialog,
  EditableContentInput,
  Inbox,
} from '../../ui-kit';
import styles from './contact-details.module.scss';
import {
  useArchiveContact,
  useContactPersonalDetails,
  useUpdateContactPersonalDetails,
} from '../../../hooks/useContact';
import { ContactDetailsSkeleton } from './skeletons';
import { ContactTags } from '../contact-tags';
import { ContactAvatar } from '../../ui-kit/molecules/organization-avatar';
import { useRecoilValue } from 'recoil';
import { contactDetailsEdit } from '../../../state';
import { JobRoleInput } from './edit';
import { IconButton } from '../../ui-kit/atoms';
import classNames from 'classnames';
import { useCreateContactJobRole } from '../../../hooks/useContactJobRole';

export const ContactPersonalDetails = ({ id }: { id: string }) => {
  const { data, loading, error } = useContactPersonalDetails({ id });
  const { isEditMode } = useRecoilValue(contactDetailsEdit);
  const { onCreateContactJobRole } = useCreateContactJobRole({ contactId: id });

  const { onUpdateContactPersonalDetails } = useUpdateContactPersonalDetails({
    contactId: id,
  });
  const [deleteConfirmationModalVisible, setDeleteConfirmationModalVisible] =
    useState(false);
  const { onArchiveContact } = useArchiveContact({ id });
  useEffect(() => {
    if (!loading && !data?.jobRoles?.length && isEditMode) {
      onCreateContactJobRole({ jobTitle: '' });
    }
  }, [loading, isEditMode]);

  if (loading) {
    return <ContactDetailsSkeleton />;
  }
  if (error) {
    return null;
  }
  return (
    <div className={styles.header}>
      <div className={styles.avatarWrapper}>
        <div className={styles.photo}>
          <ContactAvatar contactId={id} size={50} />
        </div>
        {isEditMode && (
          <>
            <IconButton
              className={styles.archiveContactButton}
              size='xxxs'
              mode='danger'
              onClick={() => setDeleteConfirmationModalVisible(true)}
              icon={<Inbox style={{ transform: 'scale(0.8)' }} />}
            />
            <DeleteConfirmationDialog
              deleteConfirmationModalVisible={deleteConfirmationModalVisible}
              setDeleteConfirmationModalVisible={
                setDeleteConfirmationModalVisible
              }
              deleteAction={() =>
                onArchiveContact().then(() =>
                  setDeleteConfirmationModalVisible(false),
                )
              }
              header='Confirm archive'
              confirmationButtonLabel='Archive contact'
              explanationText='Are you sure you want to archive this contact?'
            />
          </>
        )}
      </div>
      <div className={styles.name}>
        <div className={styles.nameAndEditButton}>
          <div className={styles.nameContainer}>
            <EditableContentInput
              isEditMode={isEditMode}
              value={data?.firstName || data?.name || ''}
              placeholder={isEditMode ? 'First name' : 'Unnamed'}
              onChange={(value: string) =>
                onUpdateContactPersonalDetails({
                  firstName: value,
                  lastName: data?.lastName || '',
                })
              }
            />
            <EditableContentInput
              isEditMode={isEditMode}
              value={data?.lastName || ''}
              placeholder={isEditMode ? 'Last name' : ''}
              onChange={(value: string) => {
                return onUpdateContactPersonalDetails({
                  lastName: value,
                  firstName: data?.firstName || '',
                });
              }}
            />
          </div>
        </div>

        {data?.jobRoles?.map((jobRole: any, index) => {
          return (
            <JobRoleInput
              key={jobRole.id}
              contactId={id}
              organization={jobRole.organization}
              jobRole={jobRole.jobTitle}
              roleId={jobRole.id}
              isEditMode={isEditMode}
              showAddButton={
                data?.jobRoles.length
                  ? data.jobRoles.length - 1 === index
                  : true
              }
            />
          );
        })}

        <ContactTags id={id} mode={isEditMode ? 'EDIT' : 'PREVIEW'} />
        <div
          className={classNames(styles.source, {
            [styles.sourceEditMode]: isEditMode,
          })}
        >
          Source:
          <div>{data?.source || 'OPENLINE'}</div>
        </div>
      </div>
    </div>
  );
};
