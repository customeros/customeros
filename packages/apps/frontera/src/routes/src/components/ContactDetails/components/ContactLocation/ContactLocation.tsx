import cityTimezone from 'city-timezones';
import { observer } from 'mobx-react-lite';

import { cn } from '@ui/utils/cn';
import { flags } from '@ui/media/flags';
import { getTimezone } from '@utils/getTimezone';
import { useStore } from '@shared/hooks/useStore';
import { Globe05 } from '@ui/media/icons/Globe05';

interface ContactLocationProps {
  contactId: string;
}

export const ContactLocation = observer(
  ({ contactId }: ContactLocationProps) => {
    const store = useStore();
    const contact = store.contacts.value.get(contactId);

    const countryA2 = contact?.value.locations?.[0]?.countryCodeA2;
    const countryA3 = contact?.value.locations?.[0]?.countryCodeA3;

    const flag = flags[countryA2 || ''];
    const city = contact?.value.locations?.[0]?.locality;

    const timezone = city
      ? cityTimezone.lookupViaCity(city).find((c) => {
          return c.iso2 === contact.value.locations?.[0].countryCodeA2;
        })?.timezone
      : null;

    return countryA2 ? (
      <div className='flex items-center cursor-not-allowed max-h-5'>
        <div className='mb-1'>{flag}</div>
        <div className={cn('flex items-center', countryA3 && 'gap-1')}>
          {countryA3 && <span className='ml-4 text-sm'>{countryA3}</span>}
          {countryA3 && city && timezone && <span>•</span>}
          {city && (
            <span className='overflow-hidden text-ellipsis whitespace-nowrap text-sm'>
              {city}
            </span>
          )}
          {city && timezone && <span>•</span>}
          {timezone && (
            <span className='w-[150px] text-sm'>
              {getTimezone(timezone || '')} local time
            </span>
          )}
        </div>
      </div>
    ) : (
      <div className='flex  items-center gap-4'>
        <Globe05 className='text-gray-500' />
        <p className='text-gray-400 text-sm'>Country</p>
      </div>
    );
  },
);
