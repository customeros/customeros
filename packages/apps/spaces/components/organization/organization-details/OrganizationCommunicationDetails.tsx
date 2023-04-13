import React from 'react';
import { useRecoilValue } from 'recoil';
import { organizationDetailsEdit } from '../../../state';
import { CommunicationDetails } from '../../ui-kit/molecules';
import { useOrganizationCommunicationChannelsDetails } from '../../../hooks/useOrganization/useOrganizationCommunicationChannelsDetails';
import {
  useAddEmailToOrganizationEmail,
  useRemoveEmailFromOrganizationEmail,
  useUpdateOrganizationEmail,
} from '../../../hooks/useOrganizationEmail';
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

  const { onAddEmailToOrganization, loading: addingEmail } =
    useAddEmailToOrganizationEmail({
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

  const { onCreateOrganizationPhoneNumber, loading: addingPhoneNumber } =
    useCreateOrganizationPhoneNumber({
      organizationId: id,
    });
  const { onUpdateOrganizationPhoneNumber } = useUpdateOrganizationPhoneNumber({
    organizationId: id,
  });
  const { onRemovePhoneNumberFromOrganization, loading: removingPhoneNumber } =
    useRemovePhoneNumberFromOrganization({
      organizationId: id,
    });

  return (
    <div style={{ marginLeft: isEditMode ? 24 : 0, marginTop: 24 }}>
      <CommunicationDetails
        id={id}
        onAddEmail={(input) => onAddEmailToOrganization(input)}
        onAddPhoneNumber={(input) => onCreateOrganizationPhoneNumber(input)}
        onRemoveEmail={(id: string) => onRemoveEmailFromOrganization(id)}
        onRemovePhoneNumber={(id: string) =>
          onRemovePhoneNumberFromOrganization(id)
        }
        onUpdateEmail={(input) => onUpdateOrganizationEmail(input)}
        onUpdatePhoneNumber={(input) => onUpdateOrganizationPhoneNumber(input)}
        // @ts-expect-error fixme
        data={data}
        loading={
          loading || addingPhoneNumber || addingEmail || removingPhoneNumber
        }
        isEditMode={isEditMode}
      />
    </div>
  );
};
