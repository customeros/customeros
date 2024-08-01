export function extractPlainText(htmlString: string): string {
  const htmlWithNewlines = htmlString.replace(/<br\s*\/?>/gi, '\n');
  const parser = new DOMParser();
  const doc = parser.parseFromString(htmlWithNewlines, 'text/html');

  const plaintext = doc.body.textContent || '';
  return plaintext.trimEnd();
}
