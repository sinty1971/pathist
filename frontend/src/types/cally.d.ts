import 'react';

declare global {
  namespace JSX {
    interface IntrinsicElements {
      'calendar-date': React.DetailedHTMLProps<React.HTMLAttributes<HTMLElement>, HTMLElement>;
      'calendar-month': React.DetailedHTMLProps<React.HTMLAttributes<HTMLElement>, HTMLElement>;
    }
  }
}

declare module 'react' {
  interface HTMLAttributes<T> {
    slot?: string;
  }
  interface SVGProps<T> {
    slot?: string;
  }
}