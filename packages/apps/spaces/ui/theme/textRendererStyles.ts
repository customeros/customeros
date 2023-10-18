export const getTextRendererStyles = (showAsInlineText?: boolean) => ({
  '& ol, ul': {
    pl: showAsInlineText ? 0 : '5',
  },
  '& h1, & h2, & h3, & h4, & h5, & h6': {
    fontSize: showAsInlineText ? 'sm' : 'md',
    fontWeight: showAsInlineText ? 'medium' : 'semibold',
  },
  '& pre': {
    whiteSpace: 'normal',
    fontSize: '12px',
    color: 'gray.700',
    border: '1px solid',
    borderColor: 'gray.300',
    borderRadius: '4',
    p: '2',
    py: '1',
    my: '2',
  },
  '& hr': {
    paddingY: 2,
  },
  '& blockquote': {
    position: 'relative',
    pl: showAsInlineText ? 0 : '3',
    borderRadius: 0,
    verticalAlign: 'bottom',

    '&:before': {
      content: '""',
      position: 'absolute',
      left: 0,
      background: 'gray.300',
      width: showAsInlineText ? 0 : '3px',
      height: '100%',
      borderRadius: '8px',
      bottom: 0,
    },
    '& p': {
      color: 'gray.500',
    },
    '& .customeros-tag': {
      color: 'gray.700',
      fontWeight: 'medium',

      '&:before': {
        content: '"#"',
      },
    },
    '& .customeros-mention': {
      color: 'gray.700',
      fontWeight: 'medium',

      '&:before': {
        content: '"@"',
      },
    },
  },
  "[aria-hidden='true']": {
    display: 'none',
  },

  // code to nicely present google meeting email notifications
  '& h2.primary-text': {
    color: 'gray.700',
    fontWeight: 'medium',
  },
  '& a.primary-button-text': {
    paddingY: 1,
    paddingX: 2,
    mb: 2,
    border: '1px solid',
    borderColor: 'primary.200',
    color: 'primary.700',
    background: 'primary.50',
    borderRadius: 'lg',
    width: 'fit-content',
    '&:hover': {
      textDecoration: 'none',
      bg: `primary.100`,
      color: `primary.700`,
      borderColor: `primary.200`,
    },
    '&:focus-visible': {
      textDecoration: 'none',
      bg: `primary.100`,
      color: `primary.700`,
      borderColor: `primary.200`,
      boxShadow: `0 0 0 4px var(--chakra-colors-primary-100)`,
    },
    '&:active': {
      textDecoration: 'none',
      bg: `primary.100`,
      color: `primary.700`,
      borderColor: `primary.200`,
    },
  },
  '& .body-container': {
    mt: 3,
    padding: 4,
    display: 'block',
    border: '1px solid',
    borderColor: 'gray.300',
    borderRadius: 'md',

    '& tr': {
      mr: 2,
      display: 'flex',
    },
    '& tbody': {
      mb: 2,
    },

    '& table': {
      marginInlineStart: 0,
      marginInlineEnd: 0,
    },
  },
  '& .main-column-table-ltr': {
    my: 3,
  },
  '& .grey-button-text:not(a)': {
    color: 'gray.700',
    width: 'fit-content',
    fontWeight: 'medium',
  },
  ...(showAsInlineText
    ? {
        '&': {
          display: 'inline',
        },
        '& *': {
          display: 'inline',
        },
      }
    : {}),
});
