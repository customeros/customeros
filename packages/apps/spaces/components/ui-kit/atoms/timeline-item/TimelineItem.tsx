import React, { PropsWithChildren, useCallback } from 'react';
import styles from './timeline-item.module.scss';
import { DateTimeUtils } from '../../../../utils';
import Image from 'next/image';
import { DataSource, ExternalSystem } from '@spaces/graphql';

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
    return 'openline_small';
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

const getZendeskBaseUrl = (externalApiUrl: string) => {
  const url = `${externalApiUrl.split('.')[0]}.zendesk.com/agent/tickets`;
  if (url.startsWith('https')) return url;
  return `https://${url}`;
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
    'data-tooltip': `From ${source.toLowerCase()}`,
  };

  if (source === DataSource.ZendeskSupport && externalLinks) {
    const zendeskUrl = `${getZendeskBaseUrl(
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
