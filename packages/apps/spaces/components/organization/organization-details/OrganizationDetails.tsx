import React, { useState } from 'react';
import styles from './organization-details.module.scss';
import {
  useDeleteOrganization,
  useOrganizationDetails,
} from '@spaces/hooks/useOrganization';
import Trash from '@spaces/atoms/icons/Trash';
import { useRecoilState } from 'recoil';
import { organizationDetailsEdit } from '../../../state';
import { DebouncedTextArea } from '@spaces/atoms/input/DebouncedTextArea';
import { OrganizationCommunicationDetails } from './OrganizationCommunicationDetails';
import { IconButton } from '@spaces/atoms/icon-button/IconButton';
import { Button } from '@spaces/atoms/button';
import { DeleteConfirmationDialog } from '@spaces/atoms/delete-confirmation-dialog';
import { EditableContentInput } from '@spaces/atoms/input/EditableContentInput';
import {
  useUpdateOrganizationDescription,
  useUpdateOrganizationIndustry,
  useUpdateOrganizationName,
  useUpdateOrganizationWebsite,
} from '@spaces/hooks/useOrganizationDetails';
import { OrganizationDetailsSkeleton } from './skeletons';
import { OrganizationCustomFields } from '@spaces/organization/organization-details/OrganizationCustomFields';
import { OrganizationSubsidiaries } from '@spaces/organization/organization-details/subsidiaries/OrganizationSubsidiaries';
import { OrganizationOwner } from '@spaces/organization/organization-details/owner';
import { ExternalLink } from '@spaces/atoms/external-link/ExternalLink';

export const OrganizationDetails = ({ id }: { id: string }) => {
  const { data, loading } = useOrganizationDetails({ id });
  const [{ isEditMode }, setOrganizationDetailsEdit] = useRecoilState(
    organizationDetailsEdit,
  );
  const [deleteConfirmationModalVisible, setDeleteConfirmationModalVisible] =
    useState(false);
  const { onUpdateOrganizationName } = useUpdateOrganizationName({
    organizationId: id,
  });
  const { onUpdateOrganizationIndustry } = useUpdateOrganizationIndustry({
    organizationId: id,
  });
  const { onUpdateOrganizationDescription } = useUpdateOrganizationDescription({
    organizationId: id,
  });
  const { onUpdateOrganizationWebsite } = useUpdateOrganizationWebsite({
    organizationId: id,
  });

  const { onDeleteOrganization } = useDeleteOrganization({
    id,
  });

  if (!data || loading) {
    return <OrganizationDetailsSkeleton />;
  }

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
            {!!data.subsidiaryOf.length && (
              <span className={styles.parent_company_name}>
                {data.subsidiaryOf[0].organization.name}
              </span>
            )}

            <EditableContentInput
              id={`organization-details-name-${id}`}
              label='Name'
              isEditMode={isEditMode}
              value={data?.name}
              placeholder={isEditMode ? 'Organization' : 'Unnamed'}
              onChange={(value: string) =>
                onUpdateOrganizationName({
                  name: value,
                  industry: data?.industry,
                  description: data?.description,
                  website: data?.website,
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
                onUpdateOrganizationIndustry({
                  industry: value,
                  name: data?.name || '',
                  description: data?.description,
                  website: data?.website,
                })
              }
            />
          )}
        </div>

        <DebouncedTextArea
          id={`organization-details-description-${id}`}
          aria-label='Description'
          isEditMode={isEditMode}
          value={data?.description || ''}
          placeholder={isEditMode ? 'Description' : ''}
          onChange={(value: string) =>
            onUpdateOrganizationDescription({
              description: value,
              industry: data?.industry,
              name: data?.name || '',
              website: data?.website,
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
                onUpdateOrganizationWebsite({
                  name: data?.name || '',
                  description: data?.description,
                  industry: data?.industry,
                  website: value,
                })
              }
            />
          )}

          {data?.website && !isEditMode && <ExternalLink url={data.website} />}
        </div>
        {isEditMode && (
          <div className={styles.deleteButton}>
            <IconButton
              label='Delete'
              size='sm'
              mode='danger'
              onClick={() => setDeleteConfirmationModalVisible(true)}
              icon={<Trash height={16} />}
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
      <OrganizationOwner id={id} />
      <OrganizationCustomFields id={id} />
      <OrganizationSubsidiaries id={id} />
    </div>
  );
};
