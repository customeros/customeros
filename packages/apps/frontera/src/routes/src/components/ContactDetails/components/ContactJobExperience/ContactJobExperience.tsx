import { observer } from 'mobx-react-lite';
import { formatDistanceToNow, differenceInCalendarMonths } from 'date-fns';

import { useStore } from '@shared/hooks/useStore';
import { GitTimeline } from '@ui/media/icons/GitTimeline';

interface ContactJobExperienceProps {
  contactId: string;
}

export const ContactJobExperience = observer(
  ({ contactId }: ContactJobExperienceProps) => {
    const store = useStore();
    const contactStore = store.contacts.value.get(contactId);

    const orgName = contactStore?.value?.primaryOrganizationName;

    const startedAt = contactStore?.value?.primaryOrganizationJobRoleStartDate;

    const timeAtOrg = startedAt ? timeAt(startedAt, orgName ?? '') : null;

    return timeAtOrg ? (
      <div className='flex items-center cursor-not-allowed text-sm'>
        <GitTimeline className='text-gray-500' />
        <p className='ml-4 first-letter:capitalize'>{timeAtOrg}</p>
      </div>
    ) : (
      <div className='flex items-center gap-4'>
        <GitTimeline className='text-gray-500' />
        <p className='text-gray-400 text-sm'>Tenure at Catalog</p>
      </div>
    );
  },
);

const timeAt = (startedAt: string, organizationName: string) => {
  if (!organizationName) return;

  const months = Math.abs(
    differenceInCalendarMonths(new Date(startedAt), new Date()),
  );

  if (months < 0) return `Less than a month at ${organizationName}`;
  if (months === 1) return `${months} month at ${organizationName}`;
  if (months > 1 && months < 12)
    return `${months} months at ${organizationName}`;
  if (months === 12) return `1 year at ${organizationName}`;
  if (months > 12)
    return `${formatDistanceToNow(new Date(startedAt))} at ${organizationName}`;
};
