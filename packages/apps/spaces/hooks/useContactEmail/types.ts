export {
  useGetContactCommunicationChannelsQuery,
  useGetContactPersonalDetailsQuery,
  useCreateContactMutation,
  useUpdateContactEmailMutation,
  useAddEmailToContactMutation,
  useRemoveEmailFromContactMutation,
  useUpdateContactPhoneNumberMutation,
  useRemovePhoneNumberFromContactMutation,
  useAddPhoneToContactMutation,
} from '../../graphQL/generated';
export type {
  Contact,
  GetContactCommunicationChannelsQuery,
  GetContactPersonalDetailsQuery,
  ContactInput,
  CreateContactMutation,
  Email,
  UpdateContactEmailMutation,
} from '../../graphQL/generated';
