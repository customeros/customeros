import React from 'react';
import Link from 'next/link';

import { ExternalSystem } from '@graphql/types';
import { Hubspot } from '@ui/media/logos/Hubspot';
import { Salesforce } from '@ui/media/logos/Salesforce';
import { getExternalUrl } from '@spaces/utils/getExternalLink';

const getIcon = (type: string) => {
  switch (type) {
    case 'SALESFORCE':
      return <Salesforce className='size-5 mr-2' />;
    case 'HUBSPOT':
      return <Hubspot className='size-5 mr-2' />;

    default:
      return '';
  }
};
export const LogEntryExternalLink: React.FC<{
  externalLink: ExternalSystem;
}> = ({ externalLink }) => {
  const icon = (() => getIcon(externalLink.type))();
  const link = getExternalUrl(`${externalLink?.externalUrl}`);

  return (
    <div className='mt-4'>
      <Link href={link}>
        {icon}
        View in {externalLink.type.toLowerCase()}
      </Link>
    </div>
  );
};
