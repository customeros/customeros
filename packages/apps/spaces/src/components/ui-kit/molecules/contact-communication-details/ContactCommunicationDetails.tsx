import React, { forwardRef } from 'react';
import styles from './contact-communication-details.module.scss';
import Image from 'next/image';

import { Button } from '../../atoms';
import { DetailItem } from './detail-item';
import { OverlayPanel } from '../../atoms/overlay-panel';
import { OverlayPanelEventType } from 'primereact';

const emailMock = [
  {
    id: '12',
    email: 'example1@example.com',
    label: 'Work',
    primary: true,
  },
  {
    id: '13',
    email: 'example2@example.com',
    label: 'Home',
    primary: false,
  },
  {
    id: '14',
    email: 'example3@example.com',
    label: 'Other',
    primary: false,
  },
];

const phones = [
  {
    id: '22',
    e164: '+1234567890',
    label: 'Mobile',
    primary: true,
  },
  {
    id: '122',
    e164: '+0987654321',
    label: 'Home',
    primary: false,
  },
  {
    id: '123',
    e164: '+1357908642',
    label: 'Work',
    primary: false,
  },
];
export const ContactCommunicationDetails = ({
  emails = emailMock,
  phoneNumbers = phones,
}: any) => {
  const addCommunicationChannelContainerRef = forwardRef<undefined | any>(
    undefined,
  );

  return (
    <div className={styles.contactDetails}>
      <ul className={styles.detailsList}>
        {emails.map((email) => (
          <DetailItem
            key={email.id}
            id={email.id}
            isPrimary={email.primary}
            label={email.label}
            data={email.email}
            onChange={() => null}
          />
        ))}
      </ul>
      <div className={styles.divider} />
      <ul className={styles.detailsList}>
        {phoneNumbers.map((phoneNr) => (
          <DetailItem
            key={phoneNr.id}
            id={phoneNr.id}
            isPrimary={phoneNr.primary}
            label={phoneNr.label}
            data={phoneNr.e164}
            onChange={() => null}
          />
        ))}
      </ul>
      <div className={styles.buttonWrapper}>
        <Button
          mode='secondary'
          onClick={(e: OverlayPanelEventType) =>
            //@ts-expect-error revisit later
            addCommunicationChannelContainerRef?.current?.toggle(e)
          }
          icon={<Image alt={''} src='/icons/plus.svg' width={15} height={15} />}
        >
          Add more details
        </Button>
        <OverlayPanel
          //@ts-expect-error revisit later
          ref={addCommunicationChannelContainerRef}
          model={[
            {
              label: 'Email',
              command: () => {
                // setEmails([
                //   ...emails,
                //   {
                //     id: undefined,
                //     email: '',
                //     label: 'asasa',
                //     primary: emails.length === 0,
                //     uiKey: uuidv4(), //TODO make sure the ID is unique in the array
                //     newItem: true, // this is used to remove the item from the emails array in case of cancel new item
                //   },
                // ]);
                //@ts-expect-error revisit later
                addCommunicationChannelContainerRef?.current?.hide();
              },
            },
            {
              label: 'Phone number',
              command: () => {
                // setPhoneNumbers([
                //   ...phoneNumbers,
                //   {
                //     id: undefined,
                //     e164: '',
                //     label: '',
                //     primary: phoneNumbers.length === 0,
                //     uiKey: uuidv4(), //TODO make sure the ID is unique in the array
                //     newItem: true, // this is used to remove the item from the phone numbers array in case of cancel new item
                //   },
                // ]);
                //@ts-expect-error revisit later
                addCommunicationChannelContainerRef?.current?.hide();
              },
            },
          ]}
        />
      </div>
    </div>
  );
};
