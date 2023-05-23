import { ArchiveContactMutation, useArchiveContactMutation } from './types';
import { useRouter } from 'next/router';
import { toast } from 'react-toastify';

interface Props {
  id: string;
}

interface Result {
  onArchiveContact: () => Promise<
    ArchiveContactMutation['contact_Archive'] | null
  >;
}
export const useArchiveContact = ({ id }: Props): Result => {
  const { push } = useRouter();
  const [archiveContactMutation, { loading, error, data }] =
    useArchiveContactMutation({
      variables: {
        id,
      },
    });
  const handleArchiveContact: Result['onArchiveContact'] = async () => {
    try {
      const response = await archiveContactMutation({
        update(cache) {
          const normalizedId = cache.identify({
            id: id,
            __typename: 'Contact',
          });
          cache.evict({ id: normalizedId });
          cache.gc();
        },
      });
      if (response) {
        push('/contact').then(() =>
          toast.success('Contact successfully archived!', {
            toastId: `contact-archive-success-${id}`,
          }),
        );
      }

      return response.data?.contact_Archive ?? null;
    } catch (err) {
      toast.error('Something went wrong while archiving contact!', {
        toastId: `archive-contact-error-${id}`,
      });
      return null;
    }
  };

  return {
    onArchiveContact: handleArchiveContact,
  };
};
