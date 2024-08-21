import { useSearchParams } from 'react-router-dom';
import React, { RefObject, startTransition } from 'react';

import { FilterItem } from '@store/types';
import { observer } from 'mobx-react-lite';
import { CheckedState } from '@radix-ui/react-checkbox';

import { useStore } from '@shared/hooks/useStore';
import { Checkbox, CheckMinus } from '@ui/form/Checkbox/Checkbox';
import { ColumnViewType, ComparisonOperator } from '@graphql/types';
import {
  CollapsibleRoot,
  CollapsibleContent,
  CollapsibleTrigger,
} from '@ui/transitions/Collapse/Collapse.tsx';

import {
  FilterHeader,
  DebouncedSearchInput,
} from '../../../shared/Filters/abstract';
import {
  CategoryHeaderLabel,
  DeliverabilityStatus,
  EmailVerificationStatus,
} from './utils.ts';

interface EmailFilterProps {
  property?: ColumnViewType;
  initialFocusRef: RefObject<HTMLInputElement>;
}

interface CheckboxOption {
  label: string;
  disabled?: boolean;
  value: EmailVerificationStatus;
}

interface CheckboxGroupProps {
  options: CheckboxOption[];
  isCategoryChecked: boolean;
  category: DeliverabilityStatus;
  onToggleCategory: (category: DeliverabilityStatus) => void;
  isOptionChecked: (value: EmailVerificationStatus) => boolean;
  onToggleOption: (
    value: EmailVerificationStatus,
    checked?: CheckedState,
  ) => void;
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

const CheckboxGroup: React.FC<CheckboxGroupProps> = ({
  category,
  options,
  onToggleCategory,
  onToggleOption,
  isOptionChecked,
  isCategoryChecked,
}) => (
  <CollapsibleRoot open={isCategoryChecked} className='flex flex-col w-full'>
    <div className='flex justify-between w-full items-center'>
      <CollapsibleTrigger asChild={false}>
        <Checkbox
          icon={<CheckMinus />}
          isChecked={isCategoryChecked}
          onChange={() => onToggleCategory(category)}
        >
          <p className='text-sm'>{CategoryHeaderLabel[category]}</p>
        </Checkbox>
      </CollapsibleTrigger>
    </div>
    <CollapsibleContent>
      <div className='flex flex-col w-full gap-2 ml-6 mt-2'>
        {options.map((option) => (
          <Checkbox
            key={option.value}
            disabled={option?.disabled}
            isChecked={isOptionChecked(option.value)}
            onChange={(checked) => onToggleOption(option.value, checked)}
          >
            <p className='text-sm'>{option.label}</p>
          </Checkbox>
        ))}
      </div>
    </CollapsibleContent>
  </CollapsibleRoot>
);

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

    const renderCheckboxGroup = (category: DeliverabilityStatus) => (
      <CheckboxGroup
        category={category}
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
