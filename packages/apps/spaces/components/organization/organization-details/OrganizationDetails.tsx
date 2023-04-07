import React from 'react';
import styles from './organization-details.module.scss';
import { useOrganizationDetails } from '../../../hooks/useOrganization';
import { Button, EditableContentInput, Link } from '../../ui-kit';
import { useUpdateOrganization } from '../../../hooks/useOrganization/useUpdateOrganization';
import { useRecoilState } from 'recoil';
import { organizationDetailsEdit } from '../../../state';
import { DebouncedTextArea } from '../../ui-kit/atoms/input/DebouncedTextArea';
export const OrganizationDetails = ({ id }: { id: string }) => {
  const { data } = useOrganizationDetails({ id });
  const [{ isEditMode }, setOrganizationDetailsEdit] = useRecoilState(
    organizationDetailsEdit,
  );
  const { onUpdateOrganization } = useUpdateOrganization({
    organizationId: id,
  });

  return (
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

        {(isEditMode || data?.industry?.length) && (
          <EditableContentInput
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
    </div>
  );
};
