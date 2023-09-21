import React from 'react';
import { ExternalSystem } from '@graphql/types';
import { Box } from '@ui/layout/Box';
import Link from 'next/link';
import { Hubspot } from '@ui/media/logos/Hubspot';
import { Salesforce } from '@ui/media/logos/Salesforce';
import { getExternalUrl } from '@spaces/utils/getExternalLink';

const getIcon = (type: string) => {
  switch (type) {
    case 'SALESFORCE':
      return <Salesforce boxSize='5' mr={2} />;
    case 'HUBSPOT':
      return <Hubspot boxSize='5' mr={2} />;

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
    <Box as={Link} href={link} mt={4}>
      {icon}
      View in {externalLink.type.toLowerCase()}
    </Box>
  );
};
