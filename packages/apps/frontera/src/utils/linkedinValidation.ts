export function validLinkedInProfileUrl(url: string): boolean {
  const re =
    /(https?:\/\/(www\.)|(www\.))?linkedin\.com\/(mwlite\/|m\/)?in\/[a-zA-Z0-9-]+\/?$/;

  return re.test(url);
}
