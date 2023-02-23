import React, { useRef, useState } from 'react';
import styles from './contact-communication-details.module.scss';
import Image from 'next/image';

import { Button } from '../../atoms';
import { DetailItem, DetailItemEditMode } from './detail-item';
import { OverlayPanel } from '../../atoms/overlay-panel';
import { OverlayPanelEventType } from 'primereact';
import {
  useContactCommunicationChannelsDetails,
  useRemoveEmailFromContactEmail,
  useAddEmailToContactEmail,
} from '../../../../hooks/useContact';
import { useAddPhoneToContactMutation } from '../../../../graphQL/generated';
import { useCreateContactPhoneNumber } from '../../../../hooks/useContactPhoneNumber';

export const ContactCommunicationDetails = ({ id }: { id: string }) => {
  const addCommunicationChannelContainerRef = useRef(null);
  const [newEmail, setNewEmail] = useState<any>(false);
  const [newPhoneNumber, setNewPhoneNumber] = useState<any>(false);
  const { data, loading, error } = useContactCommunicationChannelsDetails({
    id,
  });

  const { onAddEmailToContact } = useAddEmailToContactEmail({ contactId: id });
  const { onCreateContactPhoneNumber } = useCreateContactPhoneNumber({
    contactId: id,
  });
  const { onRemoveEmailFromContact } = useRemoveEmailFromContactEmail({
    contactId: id,
  });

  if (loading) {
    return <>LOADING</>;
  }
  if (error) {
    return <>ERROR</>;
  }

  return (
    <div className={styles.contactDetails}>
      <div className={styles.buttonWrapper}>
        <Button
          mode='secondary'
          onClick={(e: OverlayPanelEventType) =>
            //@ts-expect-error revisit later
            addCommunicationChannelContainerRef?.current?.toggle(e)
          }
          icon={<Image alt={''} src='/icons/plus.svg' width={15} height={15} />}
        >
          Add more
        </Button>
        <OverlayPanel
          ref={addCommunicationChannelContainerRef}
          model={[
            {
              label: 'Email',
              command: () => {
                setNewEmail({ primary: true, label: 'WORK', email: '' });
                //@ts-expect-error revisit later
                addCommunicationChannelContainerRef?.current?.hide();
              },
            },
            {
              label: 'Phone number',
              command: () => {
                setNewPhoneNumber({
                  primary: false,
                  label: 'MAIN',
                  phoneNumber: '',
                });
                //@ts-expect-error revisit later
                addCommunicationChannelContainerRef?.current?.hide();
              },
            },
          ]}
        />
      </div>

      <ul className={styles.detailsList}>
        {newEmail && (
          <DetailItemEditMode
            id={'new-email'}
            isPrimary={false}
            label={'WORK'}
            data={newEmail.email}
            onChange={(e) =>
              setNewEmail({ ...newEmail, email: e.target.value })
            }
            onExitEditMode={() => {
              onAddEmailToContact(newEmail).then((e) => setNewEmail(false));
            }}
          />
        )}
        {data?.emails.map((email) => (
          <DetailItem
            key={email.id}
            id={email.id}
            isPrimary={email.primary}
            label={email.label}
            data={email.email}
            onDelete={() => onRemoveEmailFromContact(email.id)}
          />
        ))}
      </ul>
      <div className={styles.divider} />
      <ul className={styles.detailsList}>
        {newPhoneNumber && (
          <DetailItemEditMode
            id={'new-phoneNumber'}
            isPrimary={false}
            label={'WORK'}
            data={newPhoneNumber.phoneNumber}
            onChange={(e) =>
              setNewPhoneNumber({
                ...newPhoneNumber,
                phoneNumber: e.target.value,
              })
            }
            onExitEditMode={() => {
              onCreateContactPhoneNumber(newPhoneNumber).then((e) =>
                setNewPhoneNumber(false),
              );
            }}
          />
        )}
        {data?.phoneNumbers.map((phoneNr) => (
          <DetailItem
            key={phoneNr.id}
            id={phoneNr.id}
            isPrimary={phoneNr.primary}
            label={phoneNr.label}
            data={phoneNr.rawPhoneNumber}
            onChange={() => null}
          />
        ))}
      </ul>
    </div>
  );
};
