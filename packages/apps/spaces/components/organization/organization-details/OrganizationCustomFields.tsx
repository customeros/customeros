import React from 'react';

import { useOrganizationCustomFields } from '@spaces/hooks/useOrganization/useOrganizationCustomFields';
import { CustomFields } from '@spaces/molecules/custom-fields';

export const OrganizationCustomFields = ({ id }: { id: string }) => {
  const { data, loading, error } = useOrganizationCustomFields({
    id,
  });

  return (
    <div
      style={{ marginLeft: 0, marginTop: data?.customFields.length ? 24 : 0 }}
    >
      <CustomFields
        id={id}
        // @ts-expect-error fixme
        data={data}
        loading={loading}
      />
    </div>
  );
};
