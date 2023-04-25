import React, { useEffect, useRef, useState } from 'react';
import styles from './communication-details.module.scss';
import Image from 'next/image';
import { OverlayPanelEventType } from 'primereact';
import { OverlayPanel } from '../../atoms/overlay-panel';
import classNames from 'classnames';
import { MenuItemCommandParams } from 'primereact/menuitem';
import {
  EmailInput,
  PhoneNumber,
  Email,
  EmailLabel,
  PhoneNumberInput,
  PhoneNumberLabel,
  EmailUpdateInput,
  PhoneNumberUpdateInput,
} from '../../../../graphQL/__generated__/generated';
import {
  Button,
  DeleteIconButton,
  EditableContentInput,
  AddIconButton,
  Checkbox,
} from '../../atoms';
import { useAutoAnimate } from '@formkit/auto-animate/react';

interface Props {
  onAddEmail: (input: EmailInput) => void;
  onAddPhoneNumber: (input: PhoneNumberInput) => void;
  onRemoveEmail: (id: string) => Promise<any>;
  onRemovePhoneNumber: (id: string) => Promise<any>;
  onUpdateEmail: (input: EmailUpdateInput) => Promise<any>;
  onUpdatePhoneNumber: (input: PhoneNumberUpdateInput) => Promise<any>;
  data: {
    emails: Array<Email>;
    phoneNumbers: Array<PhoneNumber>;
  };
  loading: boolean;
  isEditMode: boolean;
}

export const CommunicationDetails = ({
  onAddEmail,
  onAddPhoneNumber,
  onRemoveEmail,
  onRemovePhoneNumber,
  onUpdateEmail,
  onUpdatePhoneNumber,
  data,
  loading,
  isEditMode,
}: Props) => {
  const addEmailContainerRef = useRef([]);
  const addPhoneNumberContainerRef = useRef([]);
  const [canAddEmail, setAddEmail] = useState(true);
  const [canAddPhoneNumber, setAddPhoneNumber] = useState(true);
  const [animatedLabelItemsParent] = useAutoAnimate(/* optional config */);

  const handleAddEmptyEmail = () =>
    onAddEmail({
      label: EmailLabel.Work,
      primary: true,
      email: '',
    });

  const handleAddEmptyPhoneNumber = () =>
    onAddPhoneNumber({
      phoneNumber: '',
      label: PhoneNumberLabel.Main,
      primary: true,
    });

  useEffect(() => {
    if (!loading && isEditMode) {
      setTimeout(() => {
        if (!data?.emails?.length && canAddEmail) {
          handleAddEmptyEmail();
          setAddEmail(false);
        }
        if (data?.phoneNumbers?.length === 0 && canAddPhoneNumber) {
          handleAddEmptyPhoneNumber();
          setAddPhoneNumber(false);
        }
      }, 300);
    }
  }, [data]);

  const getLabelOptions = (
    label: any,
    onChange: (d: any) => void,
    type: 'phone' | 'email',
    index: number,
  ) =>
    Object.values(label).map((labelOption) => ({
      // @ts-expect-error fixme
      label: labelOption.toLowerCase(),
      command: (event: MenuItemCommandParams) => {
        event.originalEvent.stopPropagation();
        event.originalEvent.preventDefault();
        onChange(labelOption);

        if (type === 'phone') {
          //@ts-expect-error revisit later
          addPhoneNumberContainerRef.current?.[index]?.toggle(event);
        }
        if (type === 'email') {
          //@ts-expect-error revisit later
          addEmailContainerRef.current?.[index]?.toggle(event);
        }
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
                data?.emails
                  .filter((email) => (isEditMode ? true : email.email?.length))
                  .map(({ label, ...rest }, index) => (
                    <tr
                      key={`detail-item-email-label-${index}`}
                      className={classNames(styles.communicationItem, {
                        [styles.primary]: rest.primary,
                      })}
                    >
                      <th
                        className={classNames(styles.tableHeader, {
                          [styles.primary]: rest.primary,
                        })}
                        colSpan={1}
                      >
                        {isEditMode && (
                          <>
                            {index === 0 &&
                            data?.emails?.length === 1 &&
                            !rest.email ? null : (
                              <DeleteIconButton
                                onDelete={() => onRemoveEmail(rest.id)}
                                style={{
                                  position: 'absolute',
                                  left: -20,
                                }}
                              />
                            )}

                            <Button
                              mode='link'
                              style={{
                                display: 'inline-flex',
                                padding: 0,
                              }}
                              onClick={(e: OverlayPanelEventType) =>
                                //@ts-expect-error revisit later
                                addEmailContainerRef?.current?.[index]?.toggle(
                                  e,
                                )
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
                              ref={(el) =>
                                // @ts-expect-error revisit types
                                (addEmailContainerRef.current[index] = el)
                              }
                              model={getLabelOptions(
                                EmailLabel,
                                (newLabel: EmailLabel) => {
                                  onUpdateEmail({
                                    label: newLabel,
                                    id: rest.id,
                                    email: rest.email,
                                    primary: rest.primary,
                                  });
                                },
                                'email',
                                index,
                              )}
                            />
                          </>
                        )}

                        {!isEditMode && index === 0 && (
                          <Image
                            alt={'Email'}
                            src='/icons/envelope.svg'
                            width={16}
                            height={16}
                            style={{
                              position: 'absolute',
                              left: -20,
                            }}
                          />
                        )}

                        {!isEditMode && label?.toLowerCase()}
                      </th>
                    </tr>
                  ))}

              {(!!data?.emails?.length || !!data?.phoneNumbers?.length) && (
                <tr className={styles.divider} />
              )}
              {!hidePhoneNumberInReadOnlyIfNoData &&
                data?.phoneNumbers.map(({ label, ...rest }, index) => (
                  <tr
                    key={`detail-item-phone-number-label-${rest.id}`}
                    className={classNames(styles.communicationItem)}
                  >
                    <th
                      className={classNames(styles.tableHeader, {
                        [styles.primary]: rest.primary,
                      })}
                      colSpan={1}
                    >
                      {isEditMode && (
                        <>
                          {index === 0 &&
                          data?.phoneNumbers?.length === 1 &&
                          !rest.rawPhoneNumber &&
                          !rest.e164 ? null : (
                            <DeleteIconButton
                              onDelete={() => {
                                if (
                                  index === 0 &&
                                  data?.phoneNumbers?.length === 1
                                ) {
                                  onRemovePhoneNumber(rest.id).then(
                                    ({ result }) => {
                                      if (result) {
                                        handleAddEmptyPhoneNumber();
                                      }
                                    },
                                  );
                                  return;
                                }
                                return onRemovePhoneNumber(rest.id);
                              }}
                              style={{
                                position: 'absolute',
                                left: -20,
                              }}
                            />
                          )}

                          <Button
                            mode='link'
                            style={{
                              display: 'inline-flex',
                              padding: 0,
                            }}
                            onClick={(e: OverlayPanelEventType) =>
                              addPhoneNumberContainerRef?.current?.[
                                index
                                //@ts-expect-error revisit later
                              ]?.toggle(e)
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
                            ref={(element) =>
                              //@ts-expect-error revisit later
                              (addPhoneNumberContainerRef.current[index] =
                                element)
                            }
                            model={getLabelOptions(
                              PhoneNumberLabel,
                              (newLabel: PhoneNumberLabel) => {
                                onUpdatePhoneNumber({
                                  label: newLabel,
                                  id: rest.id,
                                  primary: rest?.primary || true,
                                  phoneNumber: rest.rawPhoneNumber || rest.e164,
                                });
                              },
                              'phone',
                              index,
                            )}
                          />
                        </>
                      )}
                      {!isEditMode && index === 0 && (
                        <Image
                          alt={'Phone number'}
                          src='/icons/phone.svg'
                          width={16}
                          height={16}
                          style={{
                            position: 'absolute',
                            left: -20,
                          }}
                        />
                      )}
                      {!isEditMode && label?.toLowerCase()}
                    </th>
                  </tr>
                ))}
            </thead>
            <tbody>
              {!hideEmailInReadOnlyIfNoData &&
                data?.emails
                  .filter((email) => (isEditMode ? true : email.email?.length))
                  .map(({ label, email, primary, id: emailId }, index) => {
                    return (
                      <tr
                        key={`detail-item-email-content-${index}-${emailId}`}
                        className={classNames(styles.communicationItem, {
                          [styles.primary]: primary && !isEditMode,
                        })}
                      >
                        <td
                          className={classNames(styles.communicationItem, {})}
                        >
                          <EditableContentInput
                            id={`communication-details-email-${index}-${emailId}`}
                            label='Email'
                            onChange={(value: string) =>
                              onUpdateEmail({
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
                          <td className={styles.checkboxContainer}>
                            <Checkbox
                              checked={primary}
                              type='radio'
                              label='Primary'
                              onChange={() =>
                                onUpdateEmail({
                                  id: emailId,
                                  label,
                                  email,
                                  primary: !primary,
                                })
                              }
                            />
                          </td>
                        )}

                        {index === data?.emails.length - 1 && isEditMode && (
                          <td>
                            <AddIconButton
                              onAdd={() =>
                                onAddEmail({
                                  label: EmailLabel.Work,
                                  primary: false,
                                  email: '',
                                })
                              }
                            />
                          </td>
                        )}
                      </tr>
                    );
                  })}
              {(hideEmailInReadOnlyIfNoData ||
                !!data?.phoneNumbers?.length) && (
                <tr className={styles.divider} />
              )}
              {!hidePhoneNumberInReadOnlyIfNoData &&
                data?.phoneNumbers
                  .filter((phoneNumber) =>
                    isEditMode
                      ? true
                      : phoneNumber.rawPhoneNumber?.length ||
                        phoneNumber.e164?.length,
                  )
                  .map(
                    (
                      {
                        label,
                        rawPhoneNumber,
                        e164,
                        primary,
                        id: phoneNumberId,
                      },
                      index,
                    ) => {
                      return (
                        <tr
                          key={`detail-item-phone-number-content-${index}-${phoneNumberId}`}
                          className={classNames(styles.communicationItem)}
                        >
                          <td
                            className={classNames(styles.communicationItem, {
                              [styles.primary]: primary && !isEditMode,
                            })}
                          >
                            <EditableContentInput
                              id={`communication-details-phone-number-${index}-${phoneNumberId}`}
                              label='Phone number'
                              isEditMode={isEditMode}
                              onChange={(value: string) =>
                                onUpdatePhoneNumber({
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
                            <td className={styles.checkboxContainer}>
                              <Checkbox
                                checked={primary}
                                type='radio'
                                label='Primary'
                                onChange={() =>
                                  onUpdatePhoneNumber({
                                    id: phoneNumberId,
                                    label,
                                    phoneNumber: rawPhoneNumber || e164 || '',
                                    primary: !primary,
                                  })
                                }
                              />
                            </td>
                          )}
                          {index === data?.phoneNumbers.length - 1 &&
                            isEditMode && (
                              <td>
                                <AddIconButton
                                  onAdd={() =>
                                    onAddPhoneNumber({
                                      phoneNumber: '',
                                      label: PhoneNumberLabel.Work,
                                      primary: false,
                                    })
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
