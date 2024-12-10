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

    const orgName =
      contactStore?.value?.latestOrganizationWithJobRole?.organization?.name;

    const startedAt =
      contactStore?.value?.latestOrganizationWithJobRole?.jobRole.startedAt;

    const timeAtOrg = startedAt ? timeAt(startedAt, orgName ?? '') : null;

    return (
      timeAtOrg && (
        <div className='flex items-center mb-1 cursor-not-allowed text-sm'>
          <GitTimeline className='text-gray-500' />
          <p className='ml-4 capitalize'>{timeAtOrg}</p>
        </div>
      )
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
