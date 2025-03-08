# Mailstack

- The scope of this service is to manage sending and receiving emails for any mailbox of a CustomerOS user.
- All email sending is managed by a single API
- All emails received are sent to a registered webhook

## Things that are in scope

- Full lifecycle of managing Mailstack mailboxes
- Management of all 3rd party mailboxes, like Google, Microsoft, and any IMAP/SMTP connection
- Spam filtering and the elimination of non-human emails
- Proxying all links in emails and tracking clicks
- Mailbox reputation monitoring
- Email validation before sending

## Things that are not in scope

- Email analysis or processing beyond initial spam filtering
- Notifications beyond sending the message via webhook
