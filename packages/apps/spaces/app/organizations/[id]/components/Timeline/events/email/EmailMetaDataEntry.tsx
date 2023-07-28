'use client';
import React, { FC } from 'react';

import { Text } from '@ui/typography/Text';
import { Flex } from '@ui/layout/Flex';
import { Tooltip } from '@ui/presentation/Tooltip';
import { InteractionEventParticipant } from '@graphql/types';
import { getEmailParticipantsNameAndEmail } from '@spaces/utils/getParticipantsName';

interface EmailMetaDataEntry {
  entryType: string;
  content: InteractionEventParticipant[] | string;
}
interface EmailMetaData {
  [x: string]: string;
  label: string;
}

export const EmailMetaDataEntry: FC<EmailMetaDataEntry> = ({
  entryType,
  content,
}) => {
  const data: boolean | Array<EmailMetaData> =
    typeof content !== 'string' &&
    getEmailParticipantsNameAndEmail(content, 'email');

  return (
    <Flex overflow='hidden' maxWidth={'100%'}>
      <Text as={'span'} color='gray.700' fontWeight={600} mr={1}>
        {entryType}:
      </Text>

      <Text
        color='gray.500'
        whiteSpace='nowrap'
        textOverflow='ellipsis'
        overflow='hidden'
      >
        <>
          {typeof content === 'string' && content}
          {typeof content !== 'string' &&
            !!data &&
            data.map((e, i) => {
              if (!e.label) {
                return (
                  <>
                    {e.email}
                    {i !== data.length - 1 && ', '}
                  </>
                );
              }
              return (
                <Tooltip
                  key={`email-participant-tag-${e.label}-${e.email}`}
                  label={e.email}
                  aria-label={`${e.email}`}
                  placement='top'
                  zIndex={100}
                >
                  {`${e.label} ${i !== data.length - 1 && ', '}`}
                </Tooltip>
              );
            })}
        </>
      </Text>
    </Flex>
  );
};
