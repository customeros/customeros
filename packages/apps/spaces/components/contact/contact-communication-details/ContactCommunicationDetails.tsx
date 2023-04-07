import React, { MutableRefObject, useEffect, useRef, useState } from 'react';
import styles from './contact-communication-details.module.scss';
import Image from 'next/image';
import { OverlayPanelEventType } from 'primereact';
import {
  useAddEmailToContactEmail,
  useContactCommunicationChannelsDetails,
  useRemoveEmailFromContactEmail,
} from '../../../hooks/useContact';
import {
  useCreateContactPhoneNumber,
  useRemovePhoneNumberFromContact,
  useUpdateContactPhoneNumber,
} from '../../../hooks/useContactPhoneNumber';
import { useUpdateContactEmail } from '../../../hooks/useContactEmail';
import {
  EmailLabel,
  PhoneNumberLabel,
} from '../../../graphQL/__generated__/generated';
import { Button, EditableContentInput, Stop, Trash } from '../../ui-kit';
import { OverlayPanel } from '../../ui-kit/atoms/overlay-panel';
import { ContactCommunicationDetailsSkeleton } from './skeletons';
import classNames from 'classnames';
import { CheckSquare, IconButton, Plus } from '../../ui-kit/atoms';
import { useRecoilValue } from 'recoil';
import { contactDetailsEdit } from '../../../state';
import { MenuItemCommandParams } from 'primereact/menuitem';

export const ContactCommunicationDetails = ({ id }: { id: string }) => {
  const { isEditMode } = useRecoilValue(contactDetailsEdit);
  const addEmailContainerRef = useRef(null);
  const addPhoneNumberContainerRef = useRef(null);

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

  useEffect(() => {
    if (!loading && isEditMode) {
      if (!data?.emails?.length) {
        onAddEmailToContact({
          email: '',
          label: EmailLabel.Work,
          primary: true,
        });
      }
      if (!data?.phoneNumbers?.length) {
        onCreateContactPhoneNumber({
          phoneNumber: '',
          label: PhoneNumberLabel.Main,
          primary: true,
        });
      }
    }
  }, [loading, isEditMode]);

  const getLabelOptions = (label: any, onChange: (d: any) => void, ref: any) =>
    Object.values(label).map((labelOption) => ({
      // @ts-expect-error fixme
      label: labelOption.toLowerCase(),
      command: (event: MenuItemCommandParams) => {
        event.originalEvent.stopPropagation();
        event.originalEvent.preventDefault();
        onChange(labelOption);
        ref?.current?.toggle(event);
      },
    }));

  if (!data && !error) {
    return <ContactCommunicationDetailsSkeleton />;
  }
  if (error) {
    return null;
  }

  const hideEmailInReadOnlyIfNoData =
    (!data?.emails.length ||
      (data?.emails.length === 1 && !data.emails[0].email)) &&
    !isEditMode;
  const hidePhoneNumberInReadOnlyIfNoData =
    (!data?.phoneNumbers.length ||
      (data?.phoneNumbers.length === 1 &&
        (!data.phoneNumbers[0]?.rawPhoneNumber ||
          !data.phoneNumbers[0]?.e164))) &&
    !isEditMode;

  return (
    <div className={styles.contactDetails}>
      <div className={styles.detailsList}>
        <>
          <table className={styles.table}>
            <thead>
              {!hideEmailInReadOnlyIfNoData &&
                data?.emails.map(({ label, ...rest }, index) => (
                  <tr
                    key={`detail-item-${label}-${index}`}
                    className={classNames(styles.communicationItem)}
                  >
                    <th
                      className={classNames(
                        styles.listContent,
                        styles.tableHeader,
                      )}
                      colSpan={1}
                    >
                      {isEditMode && (
                        <>
                          <IconButton
                            size={'xxxxs'}
                            mode='dangerLink'
                            style={{
                              width: '24px',
                              height: '16px',
                              position: 'absolute',
                              left: -20,
                            }}
                            onClick={() => onRemoveEmailFromContact(rest.id)}
                            icon={
                              <Trash
                                style={{
                                  transform: 'scale(0.6)',
                                  position: 'absolute',
                                }}
                              />
                            }
                          />
                          <Button
                            mode='link'
                            style={{
                              display: 'inline-flex',
                              paddingTop: 0,
                              paddingBottom: 0,
                              paddingRight: 0,
                            }}
                            onClick={(e: OverlayPanelEventType) =>
                              //@ts-expect-error revisit later
                              addEmailContainerRef?.current?.toggle(e)
                            }
                          >
                            <div className={styles.editLabelIcon}>
                              <Image
                                src='/icons/code.svg'
                                alt={'Change label'}
                                height={12}
                                width={12}
                              />
                              {label}
                            </div>
                          </Button>
                          <OverlayPanel
                            ref={addEmailContainerRef}
                            model={getLabelOptions(
                              EmailLabel,
                              (newLabel: EmailLabel) => {
                                console.log('🏷️ ----- newLabel : ', newLabel);
                                onUpdateContactEmail({
                                  label: newLabel,
                                  id: rest.id,
                                  email: rest.email,
                                  primary: rest.primary,
                                });
                              },
                              addEmailContainerRef,
                            )}
                          />
                        </>
                      )}

                      {!isEditMode && label}
                    </th>
                  </tr>
                ))}

              {(!!data?.emails?.length || !!data?.phoneNumbers?.length) && (
                <tr className={styles.divider} />
              )}
              {!hidePhoneNumberInReadOnlyIfNoData &&
                data?.phoneNumbers.map(({ label, ...rest }, index) => (
                  <tr
                    key={`detail-item-${label}-${index}`}
                    className={classNames(styles.communicationItem)}
                  >
                    <th
                      className={classNames(
                        styles.listContent,
                        styles.tableHeader,
                      )}
                      colSpan={1}
                    >
                      {isEditMode && (
                        <>
                          <IconButton
                            size={'xxxxs'}
                            mode='dangerLink'
                            style={{
                              width: '24px',
                              height: '16px',
                              position: 'absolute',
                              left: -20,
                            }}
                            onClick={() => onRemoveEmailFromContact(rest.id)}
                            icon={
                              <Trash
                                style={{
                                  transform: 'scale(0.6)',
                                }}
                              />
                            }
                          />
                          <Button
                            mode='link'
                            style={{
                              display: 'inline-flex',
                              paddingTop: 0,
                              paddingBottom: 0,
                              paddingRight: 0,
                            }}
                            onClick={(e: OverlayPanelEventType) =>
                              //@ts-expect-error revisit later
                              addPhoneNumberContainerRef?.current?.toggle(e)
                            }
                          >
                            <div className={styles.editLabelIcon}>
                              <Image
                                src='/icons/code.svg'
                                alt={'Change label'}
                                height={12}
                                width={12}
                              />
                              {label}
                            </div>
                          </Button>
                          <OverlayPanel
                            ref={addPhoneNumberContainerRef}
                            model={getLabelOptions(
                              PhoneNumberLabel,
                              (newLabel: PhoneNumberLabel) => {
                                onUpdateContactPhoneNumber({
                                  label: newLabel,
                                  id: rest.id,
                                  primary: rest?.primary || true,
                                  phoneNumber: rest.rawPhoneNumber || rest.e164,
                                });
                              },
                              addPhoneNumberContainerRef,
                            )}
                          />
                        </>
                      )}

                      {!isEditMode && label}
                    </th>
                  </tr>
                ))}
            </thead>
            <tbody>
              {!hideEmailInReadOnlyIfNoData &&
                data?.emails.map(
                  ({ label, email, primary, id: emailId }, index) => {
                    return (
                      <tr
                        key={`detail-item-${label}-${emailId}`}
                        className={classNames(styles.communicationItem, {})}
                      >
                        <td
                          className={classNames(styles.communicationItem, {})}
                        >
                          <EditableContentInput
                            onChange={(value: string) =>
                              onUpdateContactEmail({
                                id: emailId,
                                label,
                                primary: primary,
                                email: value,
                              })
                            }
                            inlineMode
                            inputSize='xxxxs'
                            value={email || ''}
                            placeholder='email'
                            isEditMode={isEditMode}
                          />
                        </td>
                        <td>
                          <Button
                            mode='text'
                            className={styles.primaryButton}
                            style={{
                              display: 'inline-flex',
                              padding: 0,
                              fontWeight: 'normal',
                            }}
                            onClick={() =>
                              onUpdateContactEmail({
                                id: emailId,
                                label,
                                email,
                                primary: !primary,
                              })
                            }
                          >
                            <div className={styles.editLabelIcon}>
                              {primary ? (
                                <CheckSquare
                                  style={{ transform: 'scale(0.6)' }}
                                />
                              ) : (
                                <Stop style={{ transform: 'scale(0.6)' }} />
                              )}

                              <span>Primary</span>
                            </div>
                          </Button>
                        </td>
                        {index === data?.emails.length - 1 && (
                          <td>
                            <IconButton
                              size={'xxxxs'}
                              mode='dangerLink'
                              style={{
                                width: '24px',
                                height: '16px',
                                position: 'relative',
                              }}
                              onClick={() => onAddEmailToContact({ email: '' })}
                              icon={
                                <Plus style={{ transform: 'scale(0.6)' }} />
                              }
                            />
                          </td>
                        )}
                      </tr>
                    );
                  },
                )}
              {(hideEmailInReadOnlyIfNoData ||
                !!data?.phoneNumbers?.length) && (
                <tr className={styles.divider} />
              )}
              {!hidePhoneNumberInReadOnlyIfNoData &&
                data?.phoneNumbers.map(
                  (
                    { label, rawPhoneNumber, e164, primary, id: phoneNumberId },
                    index,
                  ) => {
                    return (
                      <tr
                        key={`detail-item-${label}-${rawPhoneNumber || e164}`}
                        className={classNames(styles.communicationItem, {})}
                      >
                        <td
                          className={classNames(styles.communicationItem, {})}
                        >
                          <EditableContentInput
                            isEditMode={isEditMode}
                            onChange={(value: string) =>
                              onUpdateContactPhoneNumber({
                                id: phoneNumberId,
                                label,
                                phoneNumber: value,
                              })
                            }
                            inputSize='xxxxs'
                            value={rawPhoneNumber || e164 || ''}
                            placeholder='phone'
                          />
                        </td>
                        <td>
                          <Button
                            mode='text'
                            className={styles.primaryButton}
                            style={{
                              display: 'inline-flex',
                              padding: 0,
                              fontWeight: 'normal',
                            }}
                            onClick={() =>
                              onUpdateContactPhoneNumber({
                                id: phoneNumberId,
                                label,
                                phoneNumber: rawPhoneNumber || e164 || '',
                                primary: !primary,
                              })
                            }
                          >
                            <div className={styles.editLabelIcon}>
                              {primary ? (
                                <CheckSquare
                                  style={{ transform: 'scale(0.6)' }}
                                />
                              ) : (
                                <Stop style={{ transform: 'scale(0.6)' }} />
                              )}

                              <span>Primary</span>
                            </div>
                          </Button>
                        </td>
                        {index === data?.phoneNumbers.length - 1 && (
                          <td>
                            <IconButton
                              size={'xxxxs'}
                              mode='dangerLink'
                              style={{
                                width: '24px',
                                height: '16px',
                                position: 'relative',
                              }}
                              onClick={() =>
                                onCreateContactPhoneNumber({ phoneNumber: '' })
                              }
                              icon={
                                <Plus style={{ transform: 'scale(0.6)' }} />
                              }
                            />
                          </td>
                        )}
                      </tr>
                    );
                  },
                )}
            </tbody>
          </table>
        </>
      </div>
    </div>
  );
};
