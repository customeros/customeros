import { useNavigate, useLocation, useSearchParams } from 'react-router-dom';

import { useLocalStorage } from 'usehooks-ts';

export const useNavigationManager = () => {
  const navigate = useNavigate();
  const { pathname } = useLocation();
  const [searchParams] = useSearchParams();
  const [lastActivePosition, setLastActivePosition] = useLocalStorage(
    'customeros-player-last-position',
    { root: 'organization' },
  );

  const [lastSearchForPreset, setLastSearchForPreset] = useLocalStorage<{
    [key: string]: string;
  }>(`customeros-last-search-for-preset`, { root: 'root' });

  const handleItemClick = (path: string) => {
    const preset = searchParams.get('preset');

    setLastActivePosition({ ...lastActivePosition, root: path });

    if (preset) {
      const search = searchParams.get('search');

      setLastSearchForPreset({
        ...lastSearchForPreset,
        [preset]: search ?? '',
      });
    }
    navigate(`/${path}`);
  };

  const checkIsActive = (
    path: string,
    options?: { preset: string | Array<string> },
  ) => {
    const _pathName = path.split('?')[0];
    const presetParam = searchParams.get('preset');

    if (options?.preset) {
      if (Array.isArray(options.preset)) {
        return (
          pathname.startsWith(`/${_pathName}`) &&
          options.preset.includes(presetParam ?? '')
        );
      } else {
        return (
          pathname.startsWith(`/${_pathName}`) && presetParam === options.preset
        );
      }
    } else {
      return pathname.startsWith(`/${_pathName}`) && !presetParam;
    }
  };

  return { handleItemClick, checkIsActive };
};
