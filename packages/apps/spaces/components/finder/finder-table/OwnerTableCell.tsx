import React, {useEffect, useRef} from 'react';
import { FinderCell } from '@spaces/finder/finder-table/FinderTableCell';
import { OrganizationOwnerAutocomplete } from '@spaces/organization/organization-details/owner/OrganizationOwnerAutocomplete';

export const OwnerTableCell = ({
  organizationId,
  owner,
}: {
  organizationId: string;
  owner: any;
}) => {
  const op = useRef(null);

  const [editMode, setEditMode] = React.useState(!owner);

  return (
    <FinderCell
      label={
        <OrganizationOwnerAutocomplete
          id={organizationId}
          editMode={editMode}
          switchEditMode={() => setEditMode(!editMode)}
        />
      }
    />
  );
};
