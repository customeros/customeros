import React, { useEffect,  useState } from 'react';
import classNames from 'classnames';
import {
  Email,
  PhoneNumber,
  PhoneNumberInput,
  PhoneNumberLabel,
  PhoneNumberUpdateInput,
} from '@spaces/graphql';
import { EditableContentInput } from '@spaces/atoms/input/EditableContentInput';
import { AddIconButton } from '@spaces/atoms/icon-button/AddIconButton';
import { Checkbox } from '@spaces/atoms/checkbox/Checkbox';

import styles from './communication-details.module.scss';

interface Props {
  onAddPhoneNumber: (input: PhoneNumberInput) => void;
  onRemovePhoneNumber: (id: string) => Promise<any>;
  onUpdatePhoneNumber: (input: PhoneNumberUpdateInput) => Promise<any>;
  data: {
    emails: Array<Email>;
    phoneNumbers: Array<PhoneNumber>;
  };
  loading: boolean;
  isEditMode: boolean;
  phoneNumberId: string;
  rawPhoneNumber?: string | null;
  e164?: string | null;
  primary: boolean;
  index: number;
  phoneLabel: PhoneNumberLabel;
}

export const PhoneNumberDetails = ({
  onAddPhoneNumber,
  onUpdatePhoneNumber,
  data,
  loading,
  isEditMode,
  index,
  primary,
  phoneNumberId,
  rawPhoneNumber,
  e164,
  phoneLabel,
}: Props) => {
  const [canAddPhoneNumber, setAddPhoneNumber] = useState(true);

  const handleAddEmptyPhoneNumber = () =>
    onAddPhoneNumber({
      phoneNumber: '',
      label: PhoneNumberLabel.Main,
      primary: true,
    });

  useEffect(() => {
    if (!loading && isEditMode) {
      setTimeout(() => {
        if (data?.phoneNumbers?.length === 0 && canAddPhoneNumber) {
          handleAddEmptyPhoneNumber();
          setAddPhoneNumber(false);
        }
      }, 300);
    }
  }, [data]);

  const label = phoneLabel || PhoneNumberLabel.Other;
  return (
    <tr
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
          onBlur={(value: string) =>
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
      {index === data?.phoneNumbers.length - 1 && isEditMode && (
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
};
