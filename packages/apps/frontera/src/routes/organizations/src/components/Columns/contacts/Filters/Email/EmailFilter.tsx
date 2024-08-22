import { useSearchParams } from 'react-router-dom';
import React, { RefObject, MouseEvent, startTransition } from 'react';

import { FilterItem } from '@store/types';
import { observer } from 'mobx-react-lite';
import { CheckedState } from '@radix-ui/react-checkbox';

import { useStore } from '@shared/hooks/useStore';
import { ColumnViewType, ComparisonOperator } from '@graphql/types';

import { DeliverabilityStatus, EmailVerificationStatus } from './utils.ts';
import { EmailFilterValidationOptionGroup } from './EmailFilterValidationOptionGroup';
import {
  FilterHeader,
  DebouncedSearchInput,
} from '../../../shared/Filters/abstract';

interface EmailFilterProps {
  property?: ColumnViewType;
  initialFocusRef: RefObject<HTMLInputElement>;
}

interface EmailVerificationFilterValue {
  category: string;
  values: EmailVerificationStatus[];
}

const DEFAULT_FILTER: FilterItem = {
  property: ColumnViewType.ContactsEmails,
  value: '',
  active: false,
  caseSensitive: false,
  includeEmpty: true,
  operation: ComparisonOperator.Contains,
};

const DEFAULT_VERIFICATION_FILTER: FilterItem = {
  property: 'email_verification',
  value: [],
  active: false,
  caseSensitive: false,
  includeEmpty: false,
  operation: ComparisonOperator.In,
};

export const EmailFilter: React.FC<EmailFilterProps> = observer(
  ({ initialFocusRef, property }) => {
    const [searchParams] = useSearchParams();
    const preset = searchParams.get('preset');

    const store = useStore();
    const tableViewDef = store.tableViewDefs.getById(preset ?? '');

    const filter = tableViewDef?.getFilter(
      property || DEFAULT_FILTER.property,
    ) ?? { ...DEFAULT_FILTER, property: property || DEFAULT_FILTER.property };

    const verificationFilter =
      tableViewDef?.getFilter('email_verification') ??
      DEFAULT_VERIFICATION_FILTER;

    const toggle = () => {
      const isActive = filter.active || verificationFilter.active;

      tableViewDef?.setFilter({ ...filter, active: !isActive });
      tableViewDef?.setFilter({
        ...verificationFilter,
        active: false,
        value: [],
      });
    };

    const handleChange = (value: string) => {
      startTransition(() => {
        tableViewDef?.setFilter({ ...filter, value, active: true });
      });
    };

    const handleFilterCategory = (
      category: DeliverabilityStatus,
      value: EmailVerificationStatus,
      checked?: CheckedState,
    ) => {
      const currentValues =
        verificationFilter.value as EmailVerificationFilterValue[];
      const categoryIndex = currentValues.findIndex(
        (item) => item.category === getCategoryString(category),
      );

      let newValues: EmailVerificationFilterValue[];

      if (categoryIndex === -1) {
        newValues = [
          ...currentValues,
          {
            category: getCategoryString(category),
            values: checked ? [value] : [],
          },
        ];
      } else {
        newValues = currentValues.map((item, index) => {
          if (index === categoryIndex) {
            const newCategoryValues = checked
              ? [...item.values, value]
              : item.values.filter((v) => v !== value);

            return { ...item, values: newCategoryValues };
          }

          return item;
        });
      }

      tableViewDef?.setFilter({
        ...verificationFilter,
        value: newValues,
        active: newValues.some(
          (item) =>
            item.category === getCategoryString(category) ||
            item.values.length > 0,
        ),
      });
    };

    const isOptionChecked = (
      category: DeliverabilityStatus,
      value: EmailVerificationStatus,
    ): boolean => {
      const currentValues =
        verificationFilter.value as EmailVerificationFilterValue[];
      const categoryItem = currentValues.find(
        (item) => item.category === getCategoryString(category),
      );

      return categoryItem ? categoryItem.values.includes(value) : false;
    };

    const isCategoryChecked = (category: DeliverabilityStatus): boolean => {
      const currentValues =
        verificationFilter.value as EmailVerificationFilterValue[];

      return currentValues.some(
        (item) => item.category === getCategoryString(category),
      );
    };

    const getCategoryString = (category: DeliverabilityStatus): string => {
      switch (category) {
        case DeliverabilityStatus.Deliverable:
          return 'is_deliverable';
        case DeliverabilityStatus.NotDeliverable:
          return 'is_not_deliverable';
        case DeliverabilityStatus.Unknown:
          return 'is_deliverable_unknown';
      }
    };

    const getOptionsForCategory = (category: DeliverabilityStatus) => {
      switch (category) {
        case DeliverabilityStatus.Deliverable:
          return [
            { label: 'No risk', value: EmailVerificationStatus.NoRisk },
            {
              label: 'Firewall protected',
              value: EmailVerificationStatus.FirewallProtected,
            },
            {
              label: 'Free account',
              value: EmailVerificationStatus.FreeAccount,
            },
            {
              disabled: true,
              label: 'Group mailbox',
              value: EmailVerificationStatus.GroupMailbox,
            },
          ];
        case DeliverabilityStatus.NotDeliverable:
          return [
            {
              label: 'Invalid mailbox',
              value: EmailVerificationStatus.InvalidMailbox,
            },
            {
              label: 'Mailbox full',
              value: EmailVerificationStatus.MailboxFull,
            },
            {
              label: 'Incorrect email format',
              value: EmailVerificationStatus.IncorrectFormat,
            },
          ];
        case DeliverabilityStatus.Unknown:
          return [
            { label: 'Catch all', value: EmailVerificationStatus.CatchAll },
            {
              label: 'Not verified yet',
              value: EmailVerificationStatus.NotVerified,
            },
            {
              label: 'Verification in progress',
              value: EmailVerificationStatus.VerificationInProgress,
            },
          ];
        default:
          return [];
      }
    };

    const handleToggleCategory = (category: DeliverabilityStatus) => {
      const currentValues =
        verificationFilter.value as EmailVerificationFilterValue[];
      const categoryString = getCategoryString(category);
      const categoryIndex = currentValues.findIndex(
        (item) => item.category === categoryString,
      );

      let newValues: EmailVerificationFilterValue[];

      if (categoryIndex === -1) {
        newValues = [
          ...currentValues,
          { category: categoryString, values: [] },
        ];
      } else {
        newValues = currentValues.filter(
          (item) => item.category !== categoryString,
        );
      }

      tableViewDef?.setFilter({
        ...verificationFilter,
        value: newValues,
        active: newValues.length > 0,
      });
    };

    const handleOpenInfoModal = (
      e: MouseEvent<HTMLButtonElement>,
      status: EmailVerificationStatus,
    ) => {
      e.stopPropagation();
      e.preventDefault();
      store.ui.commandMenu.setType('ContactEmailVerificationInfoModal');
      store.ui.commandMenu.setContext({
        ids: [],
        entity: 'Contact',
        property: status,
      });
      store.ui.commandMenu.setOpen(true);
    };

    const renderCheckboxGroup = (category: DeliverabilityStatus) => (
      <EmailFilterValidationOptionGroup
        category={category}
        onOpenInfoModal={handleOpenInfoModal}
        onToggleCategory={handleToggleCategory}
        options={getOptionsForCategory(category)}
        isCategoryChecked={isCategoryChecked(category)}
        isOptionChecked={(value) => isOptionChecked(category, value)}
        onToggleOption={(value, checked) =>
          handleFilterCategory(category, value, checked)
        }
      />
    );

    return (
      <>
        <FilterHeader
          onToggle={toggle}
          onDisplayChange={() => {}}
          isChecked={filter.active || verificationFilter.active || false}
        />

        <DebouncedSearchInput
          value={filter.value}
          ref={initialFocusRef}
          onChange={handleChange}
          placeholder='e.g. john.doe@acme.com'
        />
        <div className='flex flex-col gap-2 mt-2 items-start'>
          {renderCheckboxGroup(DeliverabilityStatus.Deliverable)}
          {renderCheckboxGroup(DeliverabilityStatus.NotDeliverable)}
          {renderCheckboxGroup(DeliverabilityStatus.Unknown)}
        </div>
      </>
    );
  },
);

export default EmailFilter;
