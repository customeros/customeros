import { useSearchParams } from 'react-router-dom';
import React, { RefObject, useEffect, startTransition } from 'react';

import { FilterItem } from '@store/types';
import { observer } from 'mobx-react-lite';
import { CheckedState } from '@radix-ui/react-checkbox';

import { useStore } from '@shared/hooks/useStore';
import { ColumnViewType, ComparisonOperator } from '@graphql/types';
import { EmailFilterValidationOptionGroup } from '@organizations/components/Columns/contacts/Filters/Email/EmailFilterValidationOptionGroup';

import {
  FilterHeader,
  DebouncedSearchInput,
} from '../../../shared/Filters/abstract';
import {
  getCategoryString,
  DeliverabilityStatus,
  getOptionsForCategory,
  EmailVerificationStatus,
} from './utils';

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

    useEffect(() => {
      if (verificationFilter.active && !filter.active) {
        tableViewDef?.setFilter({ ...filter, active: true });
      }
    }, [verificationFilter.active, filter.active]);

    const handleFilterCategory = (
      category: DeliverabilityStatus,
      value: EmailVerificationStatus,
      checked?: CheckedState,
    ) => {
      const categoryString = getCategoryString(category);
      const currentValues =
        verificationFilter.value as EmailVerificationFilterValue[];

      const categoryMap = new Map(
        currentValues.map((item) => [item.category, item.values]),
      );

      let categoryValues = categoryMap.get(categoryString) || [];

      if (checked) {
        categoryValues = [...new Set([...categoryValues, value])];
      } else {
        categoryValues = categoryValues.filter((v) => v !== value);
      }

      categoryMap.set(categoryString, categoryValues);

      const allOptions = new Set(
        getOptionsForCategory(category).map((option) => option.value),
      );
      const allOptionsSelected =
        allOptions.size === categoryValues.length &&
        categoryValues.every((v) => allOptions.has(v));

      if (allOptionsSelected) {
        categoryValues.push(categoryString);
      } else {
        const index = categoryValues.indexOf(categoryString);

        if (index !== -1) {
          categoryValues.splice(index, 1);
        }
      }

      const newValues = Array.from(categoryMap)
        .map(([category, values]) => ({ category, values }))
        .filter((item) => item.values.length > 0);

      tableViewDef?.setFilter({
        ...verificationFilter,
        value: newValues,
        active: newValues.length > 0,
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

    const isCategoryChecked = (
      category: DeliverabilityStatus,
    ): CheckedState => {
      const currentValues =
        verificationFilter.value as EmailVerificationFilterValue[];
      const categoryItem = currentValues.find(
        (item) => item.category === getCategoryString(category),
      );

      if (!categoryItem) return false;

      const totalOptions = getOptionsForCategory(category).filter(
        (option) => !option.disabled,
      ).length;

      const validOptions = categoryItem?.values.filter(
        (e) => e !== categoryItem.category && e !== 'group_mailbox',
      ).length;

      if (validOptions === totalOptions) return true;

      return 'indeterminate';
    };

    const handleToggleCategory = (category: DeliverabilityStatus) => {
      const currentValues =
        verificationFilter.value as EmailVerificationFilterValue[];
      const categoryString = getCategoryString(category);

      // Use Map for O(1) lookup and modification
      const categoryMap = new Map(
        currentValues.map((item) => [item.category, item.values]),
      );

      const categoryValues = categoryMap.get(categoryString);

      if (!categoryValues || !categoryValues.includes(categoryString)) {
        const allOptions = new Set(
          getOptionsForCategory(category)
            .filter((option) => !option.disabled)
            .map((option) => option.value),
        );

        allOptions.add(categoryString);
        categoryMap.set(categoryString, Array.from(allOptions));
      } else {
        categoryMap.delete(categoryString);
      }

      const newValues = Array.from(categoryMap)
        .map(([category, values]) => ({ category, values }))
        .filter((item) => item.values.length > 0);

      tableViewDef?.setFilter({
        ...verificationFilter,
        value: newValues,
        active: newValues.length > 0,
      });
    };

    const handleOpenInfoModal = (
      e: React.MouseEvent,
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
          ref={initialFocusRef}
          onChange={handleChange}
          value={filter.value as string}
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
