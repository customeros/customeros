import React, { useEffect, useRef } from 'react';
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
import {
  Button,
  DeleteIconButton,
  EditableContentInput,
  Stop,
} from '../../ui-kit';
import { OverlayPanel } from '../../ui-kit/atoms/overlay-panel';
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
          label: EmailLabel.Work,
          primary: true,
          email: '',
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
                    key={`detail-item-email-label-${rest.id}`}
                    className={classNames(styles.communicationItem)}
                  >
                    <th className={classNames(styles.tableHeader)} colSpan={1}>
                      {isEditMode && (
                        <>
                          <DeleteIconButton
                            onDelete={() => onRemoveEmailFromContact(rest.id)}
                            style={{
                              position: 'absolute',
                              left: -20,
                            }}
                          />

                          <Button
                            mode='link'
                            style={{
                              display: 'inline-flex',
                              padding: 0,
                            }}
                            onClick={(e: OverlayPanelEventType) =>
                              //@ts-expect-error revisit later
                              addEmailContainerRef?.current?.toggle(e)
                            }
                          >
                            <div className={styles.editLabelIcon}>
                              {label?.toLowerCase()}
                              <Image
                                src='/icons/code.svg'
                                alt={'Change label'}
                                height={12}
                                width={12}
                              />
                            </div>
                          </Button>
                          <OverlayPanel
                            ref={addEmailContainerRef}
                            model={getLabelOptions(
                              EmailLabel,
                              (newLabel: EmailLabel) => {
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

                      {!isEditMode && label?.toLowerCase()}
                    </th>
                  </tr>
                ))}

              {(!!data?.emails?.length || !!data?.phoneNumbers?.length) && (
                <tr className={styles.divider} />
              )}
              {!hidePhoneNumberInReadOnlyIfNoData &&
                data?.phoneNumbers.map(({ label, ...rest }) => (
                  <tr
                    key={`detail-item-phone-number-label-${rest.id}`}
                    className={classNames(styles.communicationItem)}
                  >
                    <th className={classNames(styles.tableHeader)} colSpan={1}>
                      {isEditMode && (
                        <>
                          <DeleteIconButton
                            onDelete={() =>
                              onRemovePhoneNumberFromContact(rest.id)
                            }
                            style={{
                              position: 'absolute',
                              left: -20,
                            }}
                          />
                          <Button
                            mode='link'
                            style={{
                              display: 'inline-flex',
                              padding: 0,
                            }}
                            onClick={(e: OverlayPanelEventType) =>
                              //@ts-expect-error revisit later
                              addPhoneNumberContainerRef?.current?.toggle(e)
                            }
                          >
                            <div className={styles.editLabelIcon}>
                              {label?.toLowerCase()}
                              <Image
                                src='/icons/code.svg'
                                alt={'Change label'}
                                height={12}
                                width={12}
                              />
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

                      {!isEditMode && label?.toLowerCase()}
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
                        key={`detail-item-email-content-${emailId}`}
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
                            inputSize='xxxxs'
                            value={email || ''}
                            placeholder='email'
                            isEditMode={isEditMode}
                          />
                        </td>

                        {isEditMode && (
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
                        )}

                        {index === data?.emails.length - 1 && isEditMode && (
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
                                onAddEmailToContact({
                                  label: EmailLabel.Work,
                                  primary: false,
                                  email: '',
                                })
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
                        key={`detail-item-phone-number-content-${phoneNumberId}`}
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
                        {isEditMode && (
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
                        )}
                        {index === data?.phoneNumbers.length - 1 &&
                          isEditMode && (
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
                                  onCreateContactPhoneNumber({
                                    phoneNumber: '',
                                    label: PhoneNumberLabel.Work,
                                    primary: false,
                                  })
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
