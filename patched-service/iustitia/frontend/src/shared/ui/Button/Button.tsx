import type { ButtonHTMLAttributes, ReactNode } from 'react';
import styles from './Button.module.css';

type ButtonVariant = 'primary' | 'ghost' | 'icon';

interface ButtonProps extends ButtonHTMLAttributes<HTMLButtonElement> {
  variant?: ButtonVariant;
  children: ReactNode;
}

const variantClass: Record<ButtonVariant, string> = {
  primary: styles['primary']!,
  ghost: styles['ghost']!,
  icon: styles['icon']!,
};

export const Button = ({
  variant = 'primary',
  className,
  children,
  type = 'button',
  ...rest
}: ButtonProps) => {
  const classes = [styles.button, variantClass[variant], className].filter(Boolean).join(' ');
  return (
    <button type={type} className={classes} {...rest}>
      {children}
    </button>
  );
};
