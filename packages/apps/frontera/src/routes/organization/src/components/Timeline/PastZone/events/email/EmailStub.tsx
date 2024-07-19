import { FC, useMemo } from 'react';

import { convert } from 'html-to-text';

import { cn } from '@ui/utils/cn';
import { EmailParticipant } from '@graphql/types';
import { Card, CardFooter, CardContent } from '@ui/presentation/Card/Card';
import { getEmailParticipantsByType } from '@organization/components/Timeline/PastZone/events/email/utils';
import {
  getEmailParticipantsName,
  getEmailParticipantsNameAndEmail,
} from '@utils/getParticipantsName';
import { useTimelineEventPreviewMethodsContext } from '@organization/components/Timeline/shared/TimelineEventPreview/context/TimelineEventPreviewContext';

import { InteractionEventWithDate } from '../../../types';

import postStamp from '/backgrounds/organization/post-stamp.webp';

export const EmailStub: FC<{ email: InteractionEventWithDate }> = ({
  email,
}) => {
  const { openModal } = useTimelineEventPreviewMethodsContext();
  const text = convert(email?.content || '', {
    preserveNewlines: true,
    selectors: [
      {
        selector: 'a',
        options: { hideLinkHrefIfSameAsText: true, ignoreHref: true },
      },
      {
        selector: 'img',
        format: 'skip',
      },
    ],
  });

  const { to, cc, bcc } = useMemo(
    () => getEmailParticipantsByType(email?.sentTo || []),
    [email?.sentTo],
  );

  const cleanCC = useMemo(
    () =>
      getEmailParticipantsNameAndEmail(cc || [])
        .map((e) => e.label || e.email)
        .filter((data) => Boolean(data)),
    [cc],
  );
  const cleanBCC = useMemo(
    () =>
      getEmailParticipantsNameAndEmail(bcc || [])
        .map((e) => e.label || e.email)
        .filter((data) => Boolean(data)),
    [bcc],
  );
  const isSendByTenant = (email?.sentBy?.[0] as EmailParticipant)
    ?.emailParticipant?.users?.length;

  return (
    <>
      <Card
        className={cn(
          isSendByTenant ? 'ml-6' : 'ml-0',
          'shadow-xs cursor-pointer text-sm border border-gray-200 bg-white flex max-w-[549px]',
          'rounded-lg hover:shadow-md transition-all duration-200 ease-out',
        )}
        onClick={() => openModal(email.id)}
      >
        <CardContent className='px-3 py-2 pr-0 overflow-hidden flex flex-row flex-1 '>
          <div className='flex flex-col items-start gap-0'>
            <p className='line-clamp-1 leading-[21px]'>
              <span className='font-medium leading-[21px]'>
                {getEmailParticipantsName(
                  ([email?.sentBy?.[0]] as unknown as EmailParticipant[]) || [],
                )}
              </span>{' '}
              <span className='text-[#6C757D]'>emailed</span>{' '}
              <span className='font-medium mr-2'>
                {getEmailParticipantsName(to)}
              </span>{' '}
              {!!cleanBCC.length && (
                <>
                  <span className='text-[#6C757D]'>BCC:</span>{' '}
                  <span>{cleanBCC}</span>
                </>
              )}
              {!!cleanCC.length && (
                <>
                  <span className='text-[#6C757D]'>CC:</span>{' '}
                  <span>{cleanCC}</span>
                </>
              )}
            </p>

            <p className='font-semibold line-clamp-1 leading-[21px]'>
              {email.interactionSession?.name}
            </p>

            <p className='line-clamp-2 break-words'>{text}</p>
          </div>
        </CardContent>
        <CardFooter className='py-2 px-3 ml-1'>
          <img src={postStamp} alt='Email' width={48} height={70} />
        </CardFooter>
      </Card>
    </>
  );
};
