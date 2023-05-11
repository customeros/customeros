import React from 'react';
import { CustomField } from '../../../../graphQL/__generated__/generated';

import styles from './custom-fields.module.scss';
import linkifyHtml from 'linkify-html';
import sanitizeHtml from 'sanitize-html';

interface Props {
  data: {
    customFields: Array<CustomField>;
  };
  loading: boolean;
}

export const CustomFields = ({ data, loading }: Props) => {
  return (
    <div className={styles.contactDetails}>
      <div className={styles.detailsList}>
        <>
          <table className={styles.table}>
            <thead>
              {data?.customFields?.map((customField, index) => (
                <tr key={`custom-field-item-label-${index}`}>
                  <th className={styles.th}>
                    <div className={styles.label}>{customField.name}</div>
                    <div
                      dangerouslySetInnerHTML={{
                        __html: sanitizeHtml(
                          linkifyHtml(customField.value, {
                            defaultProtocol: 'https',
                            rel: 'noopener noreferrer',
                          }),
                        ),
                      }}
                    ></div>
                  </th>
                </tr>
              ))}
            </thead>
            <tbody></tbody>
          </table>
        </>
      </div>
    </div>
  );
};
