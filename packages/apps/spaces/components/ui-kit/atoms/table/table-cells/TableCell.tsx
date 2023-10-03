import styles from './table-cells.module.scss';
import React, { CSSProperties, ReactNode } from 'react';
import classNames from 'classnames';

export const TableCell = ({
  label,
  customStyleLabel,
  subLabel,
  customStyleSubLabel,
  children,
  className,
}: {
  label: string | ReactNode;
  customStyleLabel?: CSSProperties | undefined;
  subLabel?: string | ReactNode;
  customStyleSubLabel?: CSSProperties | undefined;
  className?: string;
  children?: ReactNode;
}) => {
  return (
    <div className={classNames(styles.cell)}>
      {children}

      <div
        className={classNames({ [styles.textContent]: children })}
        style={{ width: '100%' }}
      >
        <span
          className={classNames(className, styles.cellData)}
          style={{ ...customStyleLabel }}
        >
          {label}
        </span>
        {subLabel && (
          <span
            className={classNames(styles.subLabel, styles.cellData)}
            style={{ ...customStyleSubLabel }}
          >
            {subLabel}
          </span>
        )}
      </div>
    </div>
  );
};
