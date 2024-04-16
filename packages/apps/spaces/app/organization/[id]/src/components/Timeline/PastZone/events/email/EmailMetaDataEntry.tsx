'use client';
import React, { FC } from 'react';

import { Tooltip } from '@ui/overlay/Tooltip/Tooltip';
import { InteractionEventParticipant } from '@graphql/types';
import { getEmailParticipantsNameAndEmail } from '@spaces/utils/getParticipantsName';

interface EmailMetaDataEntry {
  entryType: string;
  content: InteractionEventParticipant[] | string;
}
interface EmailMetaData {
  label: string;
  [x: string]: string;
}

export const EmailMetaDataEntry: FC<EmailMetaDataEntry> = ({
  entryType,
  content,
}) => {
  const data: boolean | Array<EmailMetaData> =
    typeof content !== 'string' &&
    getEmailParticipantsNameAndEmail(content, 'email').filter(
      (e) => e?.label || e?.email,
    );
  if (typeof data !== 'boolean' && !data?.length) return null;

  return (
    <div className='flex overflow-hidden max-w-[100%]'>
      <span className='text-gray-700 font-semibold mr-1'>{entryType}:</span>

      <p className='text-gray-500 whitespace-nowrap text-ellipsis overflow-hidden'>
        <>
          {typeof content === 'string' && content}
          {typeof content !== 'string' &&
            !!data &&
            data.map((e, i) => {
              if (!e.label) {
                return (
                  <>
                    {e.email}
                    {i !== data.length - 1 ? ', ' : ''}
                  </>
                );
              }

              return (
                <React.Fragment
                  key={`email-participant-tag-${e.label}-${e.email}`}
                >
                  <Tooltip label={e.email} aria-label={`${e.email}`} side='top'>
                    {e.label}
                  </Tooltip>
                  {i !== data.length - 1 ? ',  ' : ''}
                </React.Fragment>
              );
            })}
        </>
      </p>
    </div>
  );
};
