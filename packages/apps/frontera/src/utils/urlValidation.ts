import { getExternalUrl } from '@utils/getExternalLink';

export function isValidUrl(string: string) {
  const urlWithProtocol = getExternalUrl(string);

  try {
    new URL(urlWithProtocol);

    return true;
  } catch (err) {
    return false;
  }
}
