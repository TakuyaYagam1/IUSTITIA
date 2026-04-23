import { forwardRef, type InputHTMLAttributes, type TextareaHTMLAttributes } from 'react';
import styles from './Input.module.css';

interface InputProps extends InputHTMLAttributes<HTMLInputElement> {
  label?: string;
  error?: string;
}

export const Input = forwardRef<HTMLInputElement, InputProps>(
  ({ label, error, className, id, ...rest }, ref) => {
    const inputClass = [styles.input, error && styles.invalid, className].filter(Boolean).join(' ');
    return (
      <div className={styles.field}>
        {label && (
          <label htmlFor={id} className={styles.label}>
            {label}
          </label>
        )}
        <input ref={ref} id={id} className={inputClass} {...rest} />
        {error && <span className={styles.error}>{error}</span>}
      </div>
    );
  },
);

Input.displayName = 'Input';

interface TextareaProps extends TextareaHTMLAttributes<HTMLTextAreaElement> {
  label?: string;
  error?: string;
}

export const Textarea = forwardRef<HTMLTextAreaElement, TextareaProps>(
  ({ label, error, className, id, ...rest }, ref) => {
    const areaClass = [styles.textarea, error && styles.invalid, className]
      .filter(Boolean)
      .join(' ');
    return (
      <div className={styles.field}>
        {label && (
          <label htmlFor={id} className={styles.label}>
            {label}
          </label>
        )}
        <textarea ref={ref} id={id} className={areaClass} {...rest} />
        {error && <span className={styles.error}>{error}</span>}
      </div>
    );
  },
);

Textarea.displayName = 'Textarea';
