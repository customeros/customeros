import React, { useEffect, useState } from 'react';
import classNames from 'classnames';
import {
  Email,
  EmailInput,
  EmailLabel,
  EmailUpdateInput,
  PhoneNumber,
} from '@spaces/graphql';
import { EditableContentInput } from '@spaces/atoms/input/EditableContentInput';
import { AddIconButton } from '@spaces/atoms/icon-button/AddIconButton';
import { Checkbox } from '@spaces/atoms/checkbox/Checkbox';
import styles from './communication-details.module.scss';
import { EmailValidationMessage } from '@spaces/molecules/communication-details/EmailValidationMessage';

interface Props {
  onAddEmail: (input: EmailInput) => void;
  onRemoveEmail: (id: string) => Promise<any>;
  onUpdateEmail: (input: EmailUpdateInput) => Promise<any>;
  data: {
    emails: Array<Email>;
    phoneNumbers: Array<PhoneNumber>;
  };
  isEditMode: boolean;
  primary: boolean;
  emailId: string;
  email?: string | null;
  emailLabel: EmailLabel;
  index: number;
  emailValidationDetails: any;
}

export const EmailDetails = ({
  onAddEmail,
  onUpdateEmail,
  data,
  isEditMode,
  primary,
  emailId,
  email,
  index,
  emailLabel,
  emailValidationDetails,
}: Props) => {
  const [showValidationMessage, setShowValidationMessage] = useState(false);

  const label = emailLabel || EmailLabel.Other;

  return (
    <tr
      className={classNames(styles.communicationItem, {
        [styles.primary]: primary && !isEditMode,
      })}
    >
      <td
        className={classNames(styles.communicationItem, {})}
        onMouseEnter={() => setShowValidationMessage(true)}
        onMouseLeave={() => setShowValidationMessage(false)}
      >
        <EditableContentInput
          id={`communication-details-email-${index}-${emailId}`}
          label='Email'
          onBlur={(value: string) => {
            onUpdateEmail({
              id: emailId,
              label,
              primary: primary,
              email: value,
            });
          }}
          inputSize='xxxxs'
          value={email || ''}
          placeholder='email'
          isEditMode={isEditMode}
        />

        {!!email?.length && (
          <EmailValidationMessage
            email={email}
            showValidationMessage={showValidationMessage}
            isEditMode={isEditMode}
            validationDetails={emailValidationDetails}
          />
        )}
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
};
