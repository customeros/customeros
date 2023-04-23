import { graphql } from 'msw';
import { Meeting } from '../graphQL/__generated__/generated';
import { uuid4 } from '@sentry/utils';

export const meeting: Meeting = {
  id: uuid4(),
  createdBy: {
    type: 'contact',
    contactParticipant: {
      id: '262bc5dd-fdf9-4c83-acd8-a99f8e2578e6',
    },
  },
  createdAt: new Date('2023-04-15T10:30:00Z'),
  start: new Date('2023-04-18T14:00:00Z'),
  end: new Date('2023-04-18T15:00:00Z'),
  location: 'https://zoom.us/j/1234567890',
  attendedBy: [
    {
      __typename: 'ContactParticipant',
      contactParticipant: { id: '262bc5dd-fdf9-4c83-acd8-a99f8e2578e6' },
    },
    {
      __typename: 'ContactParticipant',
      contactParticipant: { id: '2c8715af-ceff-4bbf-b855-40caec4c81cd' },
    },
    {
      __typename: 'ContactParticipant',
      contactParticipant: { id: '20970178-4c53-425b-8dcd-5d12a0184811' },
    },
    {
      __typename: 'ContactParticipant',
      contactParticipant: { id: '12215e8a-4848-45fc-9df9-b4e2e1e73cd3' },
    },
    {
      __typename: 'ContactParticipant',
      contactParticipant: { id: 'ebe65274-61b1-4833-9865-1a30764f4b7d' },
    },
  ],
  agenda: {
    html: '<ul><li>Introductions</li><li>Project update</li><li>Next steps</li></ul>',
    json: '{"items": ["Introductions", "Project update", "Next steps"]}',
  },
  attachments: [],
  note: {
    html: '',
  },

  recoding: null,
};

export const handlers = [
  // Handles a "GetUserInfo" query
  graphql.mutation('createMeeting', (req, res, ctx) => {
    console.log('🏷️ ----- req: create', req);
    console.log('🏷️ ----- res: create', res);
    return res(
      ctx.data({
        meeting_Create: {
          __typename: 'Meeting',
          ...meeting,
        },
      }),
    );
  }),
  graphql.mutation('updateMeeting', (req, res, ctx) => {
    console.log('🏷️ ----- req: update', req);
    console.log('🏷️ ----- res: update', res);
    return res(
      ctx.data({
        meeting_Update: {
          ...meeting,
          __typename: 'Meeting',
        },
      }),
    );
  }),
  graphql.mutation('meetingLinkAttachment', (req, res, ctx) => {
    console.log('🏷️ ----- req: link', req);
    console.log('🏷️ ----- res: link', res);
    return res(
      ctx.data({
        ...meeting,
        __typename: 'Meeting',
      }),
    );
  }),
  graphql.mutation('meetingUnlinkAttachment', (req, res, ctx) => {
    console.log('🏷️ ----- req: unlink', req);
    console.log('🏷️ ----- res: unlink', res);
    return res(
      ctx.data({
        ...meeting,
        __typename: 'Meeting',
      }),
    );
  }),
];
