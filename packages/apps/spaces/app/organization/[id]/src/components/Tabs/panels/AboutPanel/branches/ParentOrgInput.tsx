import React from 'react';
import { ArrowCircleBrokenUpLeft } from '@ui/media/icons/ArrowCircleBrokenUpLeft';
import { FormSelect } from '@ui/form/SyncSelect';
import { useAddSubsidiaryToOrganizationMutation } from '@organization/src/graphql/addSubsidiaryToOrganization.generated';
import { getGraphQLClient } from '@shared/util/getGraphQLClient';

interface ParentOrgInputProps {
  parentOrg: any;
  id: string;
}

export const ParentOrgInput: React.FC<ParentOrgInputProps> = ({
  parentOrg,
  id,
}) => {
  const client = getGraphQLClient();
  const addSubsidiaryToOrganizationMutation =
    useAddSubsidiaryToOrganizationMutation(client, {
      onSuccess: (data, variables, context) => {
        console.log('🏷️ ----- data: ', data);
        console.log('🏷️ ----- variables: ', variables);
      },
    });

  console.log('🏷️ ----- parentOrg: ', parentOrg);
  return (
    <FormSelect
      isClearable
      name='subsidiaryOf'
      value={parentOrg}
      onChange={(e) => {
        console.log('🏷️ ----- e: ', e);
        addSubsidiaryToOrganizationMutation.mutate({
          input: {
            organizationId: e.value,
            subOrganizationId: id,
          },
        });
      }}
      options={[
        {
          value: '09176502-a2cb-4f77-8291-e42c0708a985',
          label: 'Steyn org',
        },
        {
          value: 'b1dcb956-df32-466f-a985-9273cb972506',
          label: 'Silviu org',
        },
      ]}
      formId='organization-parent'
      placeholder='Parent organization'
      leftElement={<ArrowCircleBrokenUpLeft color='gray.500' mr='3' />}
    />
  );
};
