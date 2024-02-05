import axios from 'axios';

import { useCustomerLogo } from '@shared/state/CustomerLogo.atom';
import { toastError, toastSuccess } from '@ui/presentation/Toast';

export async function DownloadInvoice(
  invoiceId: string,
  name: string,
): Promise<unknown> {
  return fetch(`/fs/file/${invoiceId}/download`)
    .then(async (response) => {
      if (!response.ok) {
        throw new Error(`Download failed with status ${response.status}`);
      }
      // @ts-expect-error this error should not happen and type needs to be declared for desired outcome
      const blob = await response.blob({ type: 'application/pdf' });

      // Create a URL for the blob
      const blobUrl = window.URL.createObjectURL(blob);

      // Create a temporary anchor element and trigger the download
      const a = document.createElement('a');
      a.href = blobUrl;
      a.download = `invoice-${name}.pdf`; // Set the filename here
      document.body.appendChild(a); // Append the anchor to the body
      a.click(); // Simulate a click on the anchor
      document.body.removeChild(a); // Remove the anchor after clicking
      // Open the blob URL in a new tab
      window.URL.revokeObjectURL(blobUrl); // Free up memory by releasing the blob URL

      toastSuccess('Invoice downloaded successfully!', 'download-invoice');
    })
    .catch((reason) => {
      toastError(
        'Something went wrong while downloading the invoice',
        'download-invoice-error',
      );
    });
}

export function FetchLogo({
  id,
  setLogoUrl,
}: {
  id: string;
  setLogoUrl: (args: any) => void;
}) {
  return fetch(`/fs/file/${id}/download`)
    .then(async (response: any) => {
      const blob = await response.blob();
      console.log('🏷️ ----- response: ', response);
      const reader = new FileReader();
      reader.onload = function () {
        const img = new Image();
        img.src = reader.result as string;
        const dataUrl = reader.result as any;
        if (dataUrl) {
          setLogoUrl({
            logoUrl: dataUrl,
            dimensions: {
              width: img.width || '136px',
              height: img.height || '36px',
            },
          });
        } else {
          console.log('🏷️ ----- : ERROR');
        }
      };
      reader.readAsDataURL(blob);
    })
    .catch((reason: any) => {
      console.log('🏷️ ----- : ERROR CATCH', reason);
    });
}
