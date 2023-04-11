import React from 'react';
import { useRecoilValue } from 'recoil';
import { organizationDetailsEdit } from '../../../state';
import { CommunicationDetails } from '../../ui-kit/molecules';
import { useOrganizationCommunicationChannelsDetails } from '../../../hooks/useOrganization/useOrganizationCommunicationChannelsDetails';
import { useAddEmailToOrganizationEmail } from '../../../hooks/useOrganizationEmail/useAddOrganizationEmail';
import { useRemoveEmailFromOrganizationEmail } from '../../../hooks/useOrganizationEmail/useDeleteOrganizationEmail';
import { useUpdateOrganizationEmail } from '../../../hooks/useOrganizationEmail/useUpdateOrganizationEmail';
import {
  useCreateOrganizationPhoneNumber,
  useRemovePhoneNumberFromOrganization,
  useUpdateOrganizationPhoneNumber,
} from '../../../hooks/useOrganizationPhoneNumber';

export const OrganizationCommunicationDetails = ({ id }: { id: string }) => {
  const { isEditMode } = useRecoilValue(organizationDetailsEdit);

  const { data, loading, error } = useOrganizationCommunicationChannelsDetails({
    id,
  });

  const { onAddEmailToOrganization } = useAddEmailToOrganizationEmail({
    organizationId: id,
  });

  const { onRemoveEmailFromOrganization } = useRemoveEmailFromOrganizationEmail(
    {
      organizationId: id,
    },
  );
  const { onUpdateOrganizationEmail } = useUpdateOrganizationEmail({
    organizationId: id,
  });

  const { onCreateOrganizationPhoneNumber } = useCreateOrganizationPhoneNumber({
    organizationId: id,
  });
  const { onUpdateOrganizationPhoneNumber } = useUpdateOrganizationPhoneNumber({
    organizationId: id,
  });
  const { onRemovePhoneNumberFromOrganization } =
    useRemovePhoneNumberFromOrganization({
      organizationId: id,
    });

  return (
    <div style={{ marginLeft: isEditMode ? 24 : 0, marginTop: 24 }}>
      <CommunicationDetails
        id={id}
        onAddEmail={(input: any) => onAddEmailToOrganization(input)}
        onAddPhoneNumber={(input: any) =>
          onCreateOrganizationPhoneNumber(input)
        }
        onRemoveEmail={(input: any) => onRemoveEmailFromOrganization(input)}
        onRemovePhoneNumber={(input: any) =>
          onRemovePhoneNumberFromOrganization(input)
        }
        onUpdateEmail={(input: any) => onUpdateOrganizationEmail(input)}
        onUpdatePhoneNumber={(input: any) =>
          onUpdateOrganizationPhoneNumber(input)
        }
        data={data}
        loading={loading}
        isEditMode={isEditMode}
      />
    </div>
  );
};
