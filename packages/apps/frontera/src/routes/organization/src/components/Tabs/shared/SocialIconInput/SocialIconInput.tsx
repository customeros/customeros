import { useRef, useMemo } from 'react';
import { useParams } from 'react-router-dom';

import { observer } from 'mobx-react-lite';

import { useStore } from '@shared/hooks/useStore';

import { SocialInput } from './SocialInput';

interface SocialIconInputProps {
  name?: string;
  dataTest?: string;
  isReadOnly?: boolean;
  placeholder?: string;
  leftElement?: React.ReactNode;
  value?: { label: string; value: string }[];
}

export const SocialIconInput = observer(
  ({
    value,
    name = 'socialMedia',
    leftElement,
    isReadOnly,
    dataTest,
    ...rest
  }: SocialIconInputProps) => {
    const store = useStore();
    const id = useParams()?.id as string;
    const organization = store.organizations.getById(id);

    const _leftElement = useMemo(() => leftElement, [leftElement]);
    const newInputRef = useRef<HTMLInputElement>(null);

    if (!organization || !organization?.value) return null;

    const handleSocialChange = (e: React.ChangeEvent<HTMLInputElement>) => {
      const id = (e.target as HTMLInputElement).id;
      const value = e.target.value;

      if (organization) {
        const idx = organization.value?.socialMedia.findIndex(
          (s) => s.id === id,
        );

        if (typeof idx !== 'number' || idx < 0) return;

        organization.draft();
        organization.value!.socialMedia[idx].url = value;
      }
    };

    const handleSocialBlur = (e: React.ChangeEvent<HTMLInputElement>) => {
      const id = (e.target as HTMLInputElement).id;

      const idx = organization.value?.socialMedia.findIndex((s) => s.id === id);

      if (typeof idx !== 'number' || idx < 0) return;

      if (e.target?.value === '') {
        organization.draft();
        organization.value.socialMedia.splice(idx, 1);
        newInputRef.current?.focus();
        organization.commit();
      } else {
        organization.draft();
        organization.value!.socialMedia[idx].url = e.target.value;
        organization.commit();
      }
    };

    const handleSocialKeyDown = (e: React.KeyboardEvent<HTMLInputElement>) => {
      const id = (e.target as HTMLInputElement).id;

      const idx = organization.value?.socialMedia.findIndex((s) => s.id === id);

      if (typeof idx !== 'number' || idx < 0) return;
      const social = organization.value?.socialMedia[idx];

      if (!social) return;

      if ((e.target as HTMLInputElement).value === '') {
        organization.draft();
        organization.value.socialMedia.splice(idx, 1);
        organization.commit();
      } else {
        organization.draft();
        organization.value!.socialMedia[idx].url = (
          e.target as HTMLInputElement
        ).value;

        organization.commit();
      }
      newInputRef.current?.focus();
    };

    return (
      <>
        {value?.map(({ value: v, label: l }) => (
          <SocialInput
            id={v}
            key={v}
            value={l}
            name={name}
            dataTest={dataTest}
            isReadOnly={isReadOnly}
            onBlur={handleSocialBlur}
            leftElement={_leftElement}
            onChange={handleSocialChange}
            onKeyDown={handleSocialKeyDown}
            {...rest}
          />
        ))}
      </>
    );
  },
);
