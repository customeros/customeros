import { atom, useRecoilState } from 'recoil';

interface ContactCardMeta {
  expandedId?: string | null;
  initialFocusedField?: 'email' | 'name' | null;
}

export const ContactCardMetaAtom = atom<ContactCardMeta>({
  key: 'ContactCardMeta',
  default: {
    expandedId: null,
    initialFocusedField: null,
  },
});

export const useContactCardMeta = () => {
  return useRecoilState(ContactCardMetaAtom);
};
