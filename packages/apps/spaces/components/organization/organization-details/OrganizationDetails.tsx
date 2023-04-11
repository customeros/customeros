import React, { useState } from 'react';
import styles from './organization-details.module.scss';
import {
  useDeleteOrganization,
  useOrganizationDetails,
} from '../../../hooks/useOrganization';
import {
  Button,
  DeleteConfirmationDialog,
  EditableContentInput,
  Link,
  Pencil,
  Trash,
} from '../../ui-kit';
import { useUpdateOrganization } from '../../../hooks/useOrganization/useUpdateOrganization';
import { useRecoilState } from 'recoil';
import { organizationDetailsEdit } from '../../../state';
import { DebouncedTextArea } from '../../ui-kit/atoms/input/DebouncedTextArea';
import { OrganizationCommunicationDetails } from './OrganizationCommunicationDetails';
import { Check, IconButton } from '../../ui-kit/atoms';
export const OrganizationDetails = ({ id }: { id: string }) => {
  const { data } = useOrganizationDetails({ id });
  const [{ isEditMode }, setOrganizationDetailsEdit] = useRecoilState(
    organizationDetailsEdit,
  );
  const [deleteConfirmationModalVisible, setDeleteConfirmationModalVisible] =
    useState(false);
  const { onUpdateOrganization } = useUpdateOrganization({
    organizationId: id,
  });
  const { onDeleteOrganization } = useDeleteOrganization({
    id,
  });

  return (
    <div className={styles.detailsAndCommunicationChannel}>
      <div className={styles.organizationDetails}>
        <div>
          <div className={styles.editButton}>
            <div style={{ marginLeft: '4px' }}>
              <Button
                mode='secondary'
                onClick={() =>
                  setOrganizationDetailsEdit({ isEditMode: !isEditMode })
                }
              >
                {isEditMode ? 'Done' : 'Edit'}
              </Button>
            </div>
          </div>
          <h1 className={styles.name}>
            <EditableContentInput
              id={`organization-details-name-${id}`}
              label='Name'
              isEditMode={isEditMode}
              value={data?.name || ''}
              placeholder={isEditMode ? 'Organization' : 'Unnamed'}
              onChange={(value: string) =>
                onUpdateOrganization({
                  name: value,
                })
              }
            />
          </h1>

          {(isEditMode || !!data?.industry?.length) && (
            <EditableContentInput
              id={`organization-details-industry-${id}`}
              label='Industry'
              isEditMode={isEditMode}
              value={data?.industry || ''}
              placeholder={isEditMode ? 'Industry' : ''}
              onChange={(value: string) =>
                onUpdateOrganization({
                  name: data?.name || '',
                  description: data?.description,
                  industry: value,
                })
              }
            />
          )}
        </div>

        <DebouncedTextArea
          id={`organization-details-description-${id}`}
          label='Description'
          isEditMode={isEditMode}
          value={data?.description || ''}
          placeholder={isEditMode ? 'Description' : ''}
          onChange={(value: string) =>
            onUpdateOrganization({
              name: data?.name || '',
              description: value,
              industry: data?.industry || '',
            })
          }
        />
        <div>
          {isEditMode && (
            <EditableContentInput
              id={`organization-details-website-${id}`}
              label='Website'
              isEditMode={isEditMode}
              value={data?.website || ''}
              placeholder={isEditMode ? 'Website' : ''}
              onChange={(value: string) =>
                onUpdateOrganization({
                  name: data?.name || '',
                  description: data?.description,
                  industry: data?.industry,
                  website: value,
                })
              }
            />
          )}

          {data?.website && !isEditMode && (
            <Link href={data.website}> {data.website} </Link>
          )}
        </div>
        {isEditMode && (
          <div className={styles.deleteButton}>
            <IconButton
              size='sm'
              mode='danger'
              onClick={() => setDeleteConfirmationModalVisible(true)}
              icon={<Trash style={{ transform: 'scale(0.6)' }} />}
            />
            <DeleteConfirmationDialog
              deleteConfirmationModalVisible={deleteConfirmationModalVisible}
              setDeleteConfirmationModalVisible={
                setDeleteConfirmationModalVisible
              }
              deleteAction={() => {
                setDeleteConfirmationModalVisible(false);
                onDeleteOrganization();
              }}
              header='Confirm delete'
              confirmationButtonLabel='Delete organization'
              explanationText='Are you sure you want to delete this organization?'
            />
          </div>
        )}
      </div>
      <OrganizationCommunicationDetails id={id} />
    </div>
  );
};
