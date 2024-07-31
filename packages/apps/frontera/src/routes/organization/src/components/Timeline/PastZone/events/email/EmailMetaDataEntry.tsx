import React, { FC } from 'react';

import { Tooltip } from '@ui/overlay/Tooltip/Tooltip';
import { InteractionEventParticipant } from '@graphql/types';
import { getEmailParticipantsNameAndEmail } from '@utils/getParticipantsName';

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
      <Tooltip
        label={
          typeof content !== 'string' && !!data
            ? data?.map((e) => e.email).join(', ')
            : ''
        }
      >
        <p className='text-gray-500 whitespace-nowrap text-ellipsis overflow-hidden inline'>
          <>
            {typeof content === 'string' && content}
            {typeof content !== 'string' &&
              !!data &&
              data.map((e, i) => {
                if (!e.label) {
                  return (
                    <React.Fragment key={`email-participant-tag--${e?.email}`}>
                      {e.email}
                      {i !== data.length - 1 ? ', ' : ''}
                    </React.Fragment>
                  );
                }

                return (
                  <React.Fragment
                    key={`email-participant-tag-${e?.label}-${e?.email}`}
                  >
                    <Tooltip
                      side='top'
                      label={e.email}
                      className={'inline'}
                      aria-label={`${e.email}`}
                    >
                      <span>{e.label}</span>
                    </Tooltip>
                    {i !== data.length - 1 ? ',  ' : ''}
                  </React.Fragment>
                );
              })}
          </>
        </p>
      </Tooltip>
    </div>
  );
};
