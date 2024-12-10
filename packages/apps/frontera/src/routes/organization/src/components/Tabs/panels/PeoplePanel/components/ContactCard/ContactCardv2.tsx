import { set } from 'lodash';
import { observer } from 'mobx-react-lite';

import { Input } from '@ui/form/Input';
import { Avatar } from '@ui/media/Avatar';
import { DateTimeUtils } from '@utils/date';
import { Phone } from '@ui/media/icons/Phone';
import { Tag01 } from '@ui/media/icons/Tag01';
import { User03 } from '@ui/media/icons/User03';
import { IconButton } from '@ui/form/IconButton';
import { useStore } from '@shared/hooks/useStore';
import { Linkedin } from '@ui/media/icons/Linkedin';
import { Tags } from '@organization/components/Tabs/shared/';
import { InputGroup, LeftElement } from '@ui/form/InputGroup';
import { ChevronCollapse } from '@ui/media/icons/ChevronCollapse';
import {
  Card,
  CardHeader,
  CardFooter,
  CardContent,
} from '@ui/presentation/Card/Card';
import {
  Tag,
  Social,
  DataSource,
  EntityType,
} from '@shared/types/__generated__/graphql.types';

import { EmailsSection } from './components/EmailsSection';
import { ContactCardMenu, ContactLocation } from './components';
import { ContactJobExperience } from './components/ContactJobExperience';

interface ContactCardProps {
  id: string;
}
export const ContactCardv2 = observer(({ id }: ContactCardProps) => {
  const store = useStore();
  const contactStore = store.contacts.value.get(id);

  const handleCreateOption = (value: string) => {
    store.tags?.create(
      { name: value },
      {
        onSucces: (id) => {
          contactStore?.value.tags?.push({
            name: value,
            metadata: {
              id,
              source: DataSource.Openline,
              sourceOfTruth: DataSource.Openline,
              appSource: 'organization',
              created: new Date().toISOString(),
              lastUpdated: new Date().toISOString(),
            },
            entityType: EntityType.Contact,
          } as Tag);
          contactStore?.commit();
        },
      },
    );
  };

  const updatedDaysAgo =
    DateTimeUtils.getDaysSinceDate(contactStore?.value.updatedAt) === 1
      ? `Last changed ${DateTimeUtils.getDaysSinceDate(
          contactStore?.value.updatedAt,
        )} day ago`
      : `Last changed ${DateTimeUtils.getDaysSinceDate(
          contactStore?.value.updatedAt,
        )} days ago`;

  if (!contactStore) return null;

  return (
    <Card className='bg-white mt-4 px-2 pb-2.5 pt-0.5'>
      <CardHeader className='pb-2'>
        <div className='flex items-center justify-between w-full'>
          <div className='flex items-center w-full'>
            <div>
              <Avatar
                size='sm'
                variant='outlineCircle'
                name={contactStore?.value.name ?? ''}
                icon={<User03 className='text-gray-700 size-6' />}
                src={
                  contactStore?.value?.profilePhotoUrl
                    ? contactStore.value.profilePhotoUrl
                    : undefined
                }
              />
            </div>

            <div className='flex flex-col w-full ml-2'>
              <div className='flex'>
                <Input
                  size='xxs'
                  variant='unstyled'
                  placeholder='First & last name'
                  value={contactStore?.name ?? ''}
                  dataTest='org-people-contact-name'
                  onBlur={() => contactStore.commit()}
                  className='placeholder:font-medium font-medium'
                  onChange={(e) => {
                    contactStore.value.name = e.target.value;
                  }}
                />
                <IconButton
                  size='xxs'
                  variant='ghost'
                  aria-label='collapse'
                  icon={<ChevronCollapse />}
                />
                <ContactCardMenu />
              </div>

              <Input
                size='xxs'
                variant='unstyled'
                placeholder='Job title'
                onBlur={() => contactStore.commit()}
                value={contactStore?.value?.jobRoles?.[0]?.jobTitle ?? ''}
                onChange={(e) => {
                  contactStore.value.jobRoles[0].jobTitle = e.target.value;
                }}
              />
            </div>
          </div>
        </div>
      </CardHeader>
      <CardContent className='pl-2 pr-0 pb-0 gap-1 flex flex-col'>
        <ContactLocation contactId={id} />
        <ContactJobExperience contactId={id} />
        <EmailsSection contactId={id} />
        <InputGroup>
          <LeftElement>
            <Phone className='text-gray-500 mr-1' />
          </LeftElement>
          <Input
            size='xxs'
            variant='unstyled'
            placeholder='Phone number'
            dataTest='org-people-contact-phone-number'
            value={contactStore?.value.phoneNumbers?.[0]?.rawPhoneNumber ?? ''}
            onBlur={() => {
              contactStore.draft();
              contactStore.commit();
            }}
            onChange={(e) => {
              set(
                contactStore.value,
                ['phoneNumbers', 0, 'rawPhoneNumber'],
                e.target.value,
              );
            }}
          />
        </InputGroup>

        <InputGroup>
          <LeftElement>
            <Linkedin className='text-gray-500 mr-1' />
          </LeftElement>
          <Input
            size='xs'
            variant='unstyled'
            placeholder='LinkedIn profile URL'
            disabled={
              !!contactStore?.value.socials?.find((s) =>
                s.url.includes('linkedin.com'),
              )
            }
            value={
              contactStore?.value.socials?.find((s) =>
                s.url.includes('linkedin.com'),
              )?.url
            }
            onBlur={(e) => {
              contactStore.draft();
              contactStore.value.socials.push({
                id: crypto.randomUUID(),
                url: e.target.value,
              } as Social);
              contactStore.commit();
            }}
          />
        </InputGroup>

        <Tags
          placeholder='Personas'
          onCreateOption={handleCreateOption}
          dataTest='org-people-contact-personas'
          icon={<Tag01 className='text-gray-500 w-[18px] h-4 mr-4 mt-[6px] ' />}
          value={
            contactStore?.value?.tags?.map((t) => ({
              label: t.name,
              value: t.metadata.id,
            })) ?? []
          }
          onChange={(e) => {
            contactStore.value.tags = e
              .map((tag) => store.tags?.value.get(tag.value)?.value as Tag)
              .filter(Boolean) as Tag[];

            contactStore.commit();
          }}
        />
        <CardFooter className='pt-0.5 pb-0.5 px-0'>
          <span className='text-[12px] text-grayModern-500'>
            {contactStore?.value?.updatedAt && updatedDaysAgo}
          </span>
        </CardFooter>
      </CardContent>
    </Card>
  );
});
