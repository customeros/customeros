import React, { CSSProperties, FC } from 'react';
import classNames from 'classnames';
import styles from '@spaces/ui/form/select/components/select.module.scss';
import indicatorStyles from './health-selector.module.scss';
import { InlineLoader } from '@ui/presentation/inline-loader';
import { useSingleSelect } from '@spaces/ui/form/select/components/single-select/SingleSelect';

export const HealthIndicatorToggleButton: FC<{
  saving?: boolean;
  readOnly?: boolean;
  placeholder: string;
  customStyles?: CSSProperties | undefined;
  showIcon?: boolean;
}> = ({ saving, placeholder, customStyles, showIcon }) => {
  const { state, getToggleButtonProps } = useSingleSelect();
  const index = state.items.findIndex((e) => e?.value === state?.selection);
  const selectedLabel = state.items?.[index]?.label;

  return (
    <>
      <div
        className={classNames(
          styles.dropdownInput,
          indicatorStyles.selectButton,
          {
            [indicatorStyles.selectButtonEditable]:
              state.isEditing &&
              window.getSelection()?.isCollapsed &&
              selectedLabel,
          },
        )}
        role={'button'}
        tabIndex={0}
        style={{
          ...customStyles,
        }}
        aria-expanded={state.isOpen}
        aria-haspopup='menu'
        {...getToggleButtonProps()}
      >
        {
          <>
            {(selectedLabel || showIcon) && (
              <div
                className={indicatorStyles.colorIndicator}
                style={{
                  background: `var(--health-indicator-${
                    selectedLabel?.toLowerCase() || 'unset'
                  })`,
                }}
              />
            )}
            {selectedLabel || (
              <span
                className={classNames(indicatorStyles.placeholder, {
                  [indicatorStyles.placeholderEditable]: state.isEditing,
                })}
              >
                {placeholder}
              </span>
            )}
          </>
        }
      </div>

      {saving && <InlineLoader />}
    </>
  );
};
