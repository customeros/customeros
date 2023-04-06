import React, { useEffect, useState } from 'react';
import {
  Button,
  DeleteConfirmationDialog,
  EditableContentInput,
  Folder,
  Inbox,
  Trash,
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
import { useRecoilState } from 'recoil';
import { contactDetailsEdit } from '../../../state';
import { JobRoleInput } from './edit/JobRoleInput';
import { IconButton } from '../../ui-kit/atoms';
import classNames from 'classnames';

export const ContactPersonalDetails = ({ id }: { id: string }) => {
  const { data, loading, error } = useContactPersonalDetails({ id });
  const [{ isEditMode }, setContactDetailsEdit] =
    useRecoilState(contactDetailsEdit);

  const { onUpdateContactPersonalDetails } = useUpdateContactPersonalDetails({
    contactId: id,
  });
  const [deleteConfirmationModalVisible, setDeleteConfirmationModalVisible] =
    useState(false);
  const { onArchiveContact } = useArchiveContact({ id });

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
              size='xxxxs'
              mode='text'
              onClick={() => setDeleteConfirmationModalVisible(true)}
              icon={<Inbox />}
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
              placeholder='First name'
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
              placeholder='Last name'
              onChange={(value: string) => {
                return onUpdateContactPersonalDetails({
                  lastName: value,
                  firstName: data?.firstName || '',
                });
              }}
            />
          </div>

          <div style={{ marginLeft: '4px' }}>
            <Button
              className={styles.editButton}
              mode='secondary'
              onClick={() => setContactDetailsEdit({ isEditMode: !isEditMode })}
            >
              {isEditMode ? 'Done' : 'Edit'}
            </Button>
          </div>
        </div>

        {(!data?.jobRoles.length
          ? [{ organization: { id: '', name: '' }, jobTitle: '' }]
          : data?.jobRoles
        )?.map((jobRole: any, index) => {
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
          <span>Source:</span>
          {data?.source || 'OPENLINE'}
        </div>
      </div>
    </div>
  );
};
