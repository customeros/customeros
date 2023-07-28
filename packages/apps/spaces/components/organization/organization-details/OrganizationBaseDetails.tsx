import React, { useState } from 'react';
import styles from './organization-details.module.scss';
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
import { OrganizationHealthIndicator } from '@spaces/organization/organization-details/health-indicator';
import { OrganizationDetailsProps } from '@spaces/organization/organization-details/type';

export const OrganizationBaseDetails = ({
  id,
  loading,
  organization,
}: OrganizationDetailsProps) => {
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



  if (loading) {
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
            {!!organization?.subsidiaryOf.length && (
              <span className={styles.parent_company_name}>
                {organization?.subsidiaryOf[0].organization.name}
              </span>
            )}

            <EditableContentInput
              id={`organization-details-name-${id}`}
              label='Name'
              isEditMode={isEditMode}
              value={organization?.name}
              placeholder={isEditMode ? 'Organization' : 'Unnamed'}
              onBlur={(value: string) =>
                onUpdateOrganizationName({
                  name: value,
                  industry: organization?.industry,
                  description: organization?.description,
                  website: organization?.website,
                })
              }
            />
          </h1>

          {(isEditMode || !!organization?.industry?.length) && (
            <EditableContentInput
              id={`organization-details-industry-${id}`}
              label='Industry'
              isEditMode={isEditMode}
              value={organization?.industry || ''}
              placeholder={isEditMode ? 'Industry' : ''}
              onBlur={(value: string) =>
                onUpdateOrganizationIndustry({
                  industry: value,
                  name: organization?.name || '',
                  description: organization?.description,
                  website: organization?.website,
                })
              }
            />
          )}
        </div>

        <DebouncedTextArea
          id={`organization-details-description-${id}`}
          aria-label='Description'
          isEditMode={isEditMode}
          value={organization?.description || ''}
          placeholder={isEditMode ? 'Description' : ''}
          onChange={(value: string) =>
            onUpdateOrganizationDescription({
              description: value,
              industry: organization?.industry,
              name: organization?.name || '',
              website: organization?.website,
            })
          }
        />
        <div>
          {isEditMode && (
            <EditableContentInput
              id={`organization-details-website-${id}`}
              label='Website'
              isEditMode={isEditMode}
              value={organization?.website || ''}
              placeholder={isEditMode ? 'Website' : ''}
              onBlur={(value: string) =>
                onUpdateOrganizationWebsite({
                  name: organization?.name || '',
                  description: organization?.description,
                  industry: organization?.industry,
                  website: value,
                })
              }
            />
          )}

          {organization?.website && !isEditMode && (
            <ExternalLink url={organization.website} />
          )}
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
              }}
              header='Confirm delete'
              confirmationButtonLabel='Delete organization'
              explanationText='Are you sure you want to delete this organization?'
            />
          </div>
        )}
      </div>
      <OrganizationCommunicationDetails
        id={id}
        organization={organization}
        loading={loading}
      />
      <OrganizationOwner id={id} owner={organization?.owner} />
      <OrganizationHealthIndicator
        id={id}
        healthIndicator={organization?.healthIndicator}
      />
      <OrganizationCustomFields
        //@ts-expect-error fixme
        customFields={organization?.customFields}
      />
      <OrganizationSubsidiaries
        id={id}
        subsidiaries={organization?.subsidiaries}
      />
    </div>
  );
};
