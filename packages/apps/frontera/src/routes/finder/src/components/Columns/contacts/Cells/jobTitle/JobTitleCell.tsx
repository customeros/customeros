import { useRef } from 'react';

import { observer } from 'mobx-react-lite';

import { Edit03 } from '@ui/media/icons/Edit03';
import { IconButton } from '@ui/form/IconButton';
import { useStore } from '@shared/hooks/useStore';

interface JobTitleCellProps {
  contactId: string;
}

export const JobTitleCell = observer(({ contactId }: JobTitleCellProps) => {
  const store = useStore();

  const ref = useRef(null);

  const contactStore = store.contacts.value.get(contactId);
  const jobRoles = store.contacts.getById(String(contactId))?.jobRoles;
  const findPrimaryJobRole = jobRoles?.find(
    (j) => j.primary && j.contact?.metadata.id === contactId,
  );
  const jobRolesStore = store.jobRoles.getById(findPrimaryJobRole?.id || '');

  const enrichingStatus = contactStore?.isEnriching;

  return (
    <div ref={ref} className='flex justify-between gap-2 group/jobTitle'>
      <div className='flex gap-2 truncate'>
        {!jobRolesStore?.value.jobTitle && (
          <p className='text-gray-400'>
            {enrichingStatus ? 'Enriching...' : 'Not set'}
          </p>
        )}
        {jobRolesStore?.value.jobTitle && (
          <p className='overflow-ellipsis overflow-hidden '>
            {jobRolesStore?.value.jobTitle}
          </p>
        )}
        <IconButton
          size='xxs'
          variant='ghost'
          icon={<Edit03 />}
          className='opacity-0 group-hover/jobTitle:opacity-100 mt-[3px]'
          aria-label={`navigate-to-${contactStore?.value.primaryOrganizationName}`}
          onClick={() => {
            store.ui.commandMenu.setType('EditJobTitle');
            store.ui.commandMenu.setOpen(true);
          }}
        />
      </div>
    </div>
  );
});
