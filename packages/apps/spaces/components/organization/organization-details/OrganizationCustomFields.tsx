import React from 'react';
import { CustomFields } from '@spaces/molecules/custom-fields';
import { CustomField } from '@spaces/graphql';

export const OrganizationCustomFields = ({
  customFields,
}: {
  customFields?: Array<CustomField> | null;
}) => {
  if (!customFields) {
    return null;
  }

  return (
    <div style={{ marginLeft: 0, marginTop: customFields.length ? 24 : 0 }}>
      <CustomFields customFields={customFields} />
    </div>
  );
};
