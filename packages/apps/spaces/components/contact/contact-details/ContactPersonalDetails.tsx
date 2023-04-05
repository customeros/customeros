import React, { useState } from 'react';
import {
  Avatar,
  Button,
  DebouncedInput,
  EditableContentInput,
} from '../../ui-kit';
import styles from './contact-details.module.scss';
import {
  useContactPersonalDetails,
  useCreateContact,
  useUpdateContactPersonalDetails,
} from '../../../hooks/useContact';
import { ContactDetailsSkeleton } from './skeletons';
import { ContactTags } from '../contact-tags';
import { ContactAvatar } from '../../ui-kit/molecules/organization-avatar';
import { useRecoilState } from 'recoil';
import { contactDetailsEdit, editorMode } from '../../../state';
import { JobRoleInput } from './edit/JobRoleInput';
import { User } from '../../ui-kit/atoms';

export const ContactPersonalDetails = ({ id }: { id: string }) => {
  const { data, loading, error } = useContactPersonalDetails({ id });
  const [{ isEditMode }, setContactDetailsEdit] =
    useRecoilState(contactDetailsEdit);

  const { onUpdateContactPersonalDetails } = useUpdateContactPersonalDetails({
    contactId: id,
  });

  if (loading) {
    return <ContactDetailsSkeleton />;
  }
  if (error) {
    return <>ERROR</>;
  }
  return (
    <div className={styles.header}>
      <div className={styles.photo}>
        <ContactAvatar contactId={id} size={50} />
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

        {(
          data?.jobRoles || [
            { organization: { id: '', name: '' }, jobTitle: '' },
          ]
        )?.map((jobRole: any) => {
          return (
            <JobRoleInput
              key={jobRole.id}
              contactId={id}
              organization={jobRole.organization}
              jobRole={jobRole.jobTitle}
              roleId={jobRole.id}
            />
          );
        })}

        <ContactTags id={id} mode={isEditMode ? 'EDIT' : 'PREVIEW'} />
        <div className={styles.source}>
          <span>Source:</span>
          {data?.source || 'OPENLINE'}
        </div>
      </div>
    </div>
  );
};
