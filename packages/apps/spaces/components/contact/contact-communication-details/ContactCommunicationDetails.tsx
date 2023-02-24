import React, { useRef, useState } from 'react';
import styles from './contact-communication-details.module.scss';
import Image from 'next/image';

import { DetailItem, DetailItemEditMode } from './detail-item';
import { OverlayPanelEventType } from 'primereact';
import {
  useContactCommunicationChannelsDetails,
  useRemoveEmailFromContactEmail,
  useAddEmailToContactEmail,
} from '../../../hooks/useContact';
import {
  useCreateContactPhoneNumber,
  useRemovePhoneNumberFromContact,
  useUpdateContactPhoneNumber,
} from '../../../hooks/useContactPhoneNumber';
import { useUpdateContactEmail } from '../../../hooks/useContactEmail';
import {
  Email,
  PhoneNumber,
  PhoneNumberUpdateInput,
} from '../../../graphQL/generated';
import { Button } from '../../ui-kit';
import { OverlayPanel } from '../../ui-kit/atoms/overlay-panel';
import { ListSkeleton } from './skeletons/ListSkeleton';
import { ContactCommunicationDetailsSkeleton } from './skeletons';

export const ContactCommunicationDetails = ({ id }: { id: string }) => {
  const addCommunicationChannelContainerRef = useRef(null);
  const [newEmail, setNewEmail] = useState<any>(false);
  const [newPhoneNumber, setNewPhoneNumber] = useState<any>(false);
  const { data, loading, error } = useContactCommunicationChannelsDetails({
    id,
  });

  const { onAddEmailToContact } = useAddEmailToContactEmail({ contactId: id });

  const { onRemoveEmailFromContact } = useRemoveEmailFromContactEmail({
    contactId: id,
  });
  const { onUpdateContactEmail } = useUpdateContactEmail({
    contactId: id,
  });

  const { onCreateContactPhoneNumber } = useCreateContactPhoneNumber({
    contactId: id,
  });
  const { onUpdateContactPhoneNumber } = useUpdateContactPhoneNumber({
    contactId: id,
  });
  const { onRemovePhoneNumberFromContact } = useRemovePhoneNumberFromContact({
    contactId: id,
  });

  if (loading) {
    return <ContactCommunicationDetailsSkeleton />;
  }
  if (error) {
    return null;
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
            mode='ADD'
            id={'new-email'}
            isPrimary={false}
            label={'WORK'}
            value={newEmail.email}
            onChange={(e) =>
              setNewEmail({ ...newEmail, email: e.target.value })
            }
            onChangeLabelAndPrimary={(newValue) =>
              setNewEmail({ ...newEmail, ...newValue })
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
            onChangeLabelAndPrimary={(newValue: Email) =>
              onUpdateContactEmail(newValue, email as Email)
            }
            onDelete={() => onRemoveEmailFromContact(email.id)}
          />
        ))}
      </ul>
      <div className={styles.divider} />
      <ul className={styles.detailsList}>
        {newPhoneNumber && (
          <DetailItemEditMode
            mode='ADD'
            id={'new-phoneNumber'}
            isPrimary={false}
            label={'WORK'}
            value={newPhoneNumber.phoneNumber}
            onChange={(e) =>
              setNewPhoneNumber({
                ...newPhoneNumber,
                phoneNumber: e.target.value,
              })
            }
            onChangeLabelAndPrimary={(newValue) =>
              setNewPhoneNumber({ ...newPhoneNumber, ...newValue })
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
            // @ts-expect-error this should be revisited on phoneNumber schema change
            data={phoneNr?.rawPhoneNumber || phoneNr?.e164}
            onChange={() => null}
            onChangeLabel={(newValue: PhoneNumberUpdateInput) =>
              onUpdateContactPhoneNumber(newValue, phoneNr as PhoneNumber)
            }
            onDelete={() => onRemovePhoneNumberFromContact(phoneNr.id)}
          />
        ))}
      </ul>
    </div>
  );
};
