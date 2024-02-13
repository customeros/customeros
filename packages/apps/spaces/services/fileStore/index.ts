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
