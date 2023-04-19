import { graphql } from 'msw';

export const meeting = {
  createdBy: {
    name: 'John Doe',
    email: 'johndoe@example.com',
    phone: '+1 123-456-7890',
  },
  createdAt: new Date('2023-04-15T10:30:00Z'),
  startDate: new Date('2023-04-18T14:00:00Z'),
  endDate: new Date('2023-04-18T15:00:00Z'),
  location: {
    name: 'Zoom Meeting',
    address: 'https://zoom.us/j/1234567890',
  },
  attendees: [
    { id: '262bc5dd-fdf9-4c83-acd8-a99f8e2578e6' },
    { id: '2c8715af-ceff-4bbf-b855-40caec4c81cd' },
    { id: '20970178-4c53-425b-8dcd-5d12a0184811' },
    { id: '262bc5dd-fdf9-4c83-acd8-a99f8e2578e6' },
    { id: '2c8715af-ceff-4bbf-b855-40caec4c81cd' },
    { id: '20970178-4c53-425b-8dcd-5d12a0184811' },
  ],
  agenda: {
    html: '<ul><li>Introductions</li><li>Project update</li><li>Next steps</li></ul>',
    json: '{"items": ["Introductions", "Project update", "Next steps"]}',
  },
  attachments: [],
  note: '',
  transcription: {
    summary: '',
    transcription: '',
  },
  recording: null,
};

export const handlers = [
  // Handles a "GetUserInfo" query
  graphql.mutation('CreateMeeting', (req, res, ctx) => {
    return res(
      ctx.data({
        ...meeting,
        __typename: 'Meeting',
      }),
    );
  }),
];
