import React, { PropsWithChildren, useCallback } from 'react';
import styles from './timeline-item.module.scss';
import { DateTimeUtils } from '../../../../utils';
import Image from 'next/image';
import { DataSource, ExternalSystem } from '@spaces/graphql';
import { getZendeskIssueBaseUrl } from '@spaces/utils/getZendeskBaseUrl';

interface Props {
  children: React.ReactNode;
  createdAt?: string | number;
  first?: boolean;
  contentClassName?: any;
  hideTimeTick?: boolean;
  source: string;
  externalLinks?: ExternalSystem[];
}

export const TimelineItem: React.FC<Props> = ({
  children,
  createdAt,
  contentClassName,
  hideTimeTick,
  source = '',
  externalLinks,
  ...rest
}) => {
  const getSourceLogo = useCallback(() => {
    if (source === DataSource.ZendeskSupport) return 'zendesksupport';
    if (source === DataSource.Hubspot) return 'hubspot';
    return 'customer-os-small';
  }, [source]);

  return (
    <div className={`${styles.timelineItem}`}>
      {!hideTimeTick && (
        <>
          {createdAt ? (
            <div className={styles.when}>
              <div className={styles.timeAgo}>
                {DateTimeUtils.timeAgo(createdAt, {
                  addSuffix: true,
                })}
              </div>
              <div className={styles.metadata}>
                {DateTimeUtils.format(createdAt)}{' '}
                {!!source.length && (
                  <SourceIcon source={source} externalLinks={externalLinks}>
                    <Image
                      className={styles.logo}
                      src={`/logos/${getSourceLogo()}.svg`}
                      alt={source}
                      height={16}
                      width={16}
                    />
                  </SourceIcon>
                )}
              </div>
            </div>
          ) : (
            'Date not available'
          )}
        </>
      )}

      <div className={`${styles.content} ${contentClassName}`} {...rest}>
        {children}
      </div>
    </div>
  );
};

interface SourceIconProps {
  source: DataSource | string;
  externalLinks?: ExternalSystem[];
}

const SourceLabels: Record<DataSource, string> = {
  HUBSPOT: 'Hubspot',
  ZENDESK_SUPPORT: 'Zendesk Support',
  OPENLINE: 'Openline',
  NA: 'N/A',
};

const SourceIcon = ({
  source,
  children,
  externalLinks,
}: PropsWithChildren<SourceIconProps>) => {
  const issueExternalId = externalLinks?.[0]?.externalId ?? '';
  const issueExternalApiUrl = externalLinks?.[0]?.externalUrl ?? '';

  const commonProps = {
    className: styles.sourceLogo,
    'data-tooltip': `From ${SourceLabels[source as DataSource]}`,
  };

  if (source === DataSource.ZendeskSupport && externalLinks) {
    const zendeskUrl = `${getZendeskIssueBaseUrl(
      issueExternalApiUrl,
    )}/${issueExternalId}`;

    return (
      <a href={zendeskUrl} target='_blank' {...commonProps}>
        {children}
      </a>
    );
  }

  return <div {...commonProps}>{children}</div>;
};
