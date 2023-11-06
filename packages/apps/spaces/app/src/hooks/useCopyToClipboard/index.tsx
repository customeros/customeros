import { useState } from 'react';

import { toastSuccess } from '@ui/presentation/Toast';

type CopiedValue = string | null;
type CopyFn = (text: string, message?: string) => Promise<boolean>; // Return success

export function useCopyToClipboard(): [CopiedValue, CopyFn] {
  const [copiedText, setCopiedText] = useState<CopiedValue>(null);

  const copy: CopyFn = async (text, message = 'Link copied') => {
    if (!navigator?.clipboard) {
      return false;
    }

    // Try to save to clipboard then save it in the state if worked
    try {
      await navigator.clipboard.writeText(text);
      setCopiedText(text);
      toastSuccess(message, `copied-to-clipboard${text}`);

      return true;
    } catch (error) {
      setCopiedText(null);

      return false;
    }
  };

  return [copiedText, copy];
}
