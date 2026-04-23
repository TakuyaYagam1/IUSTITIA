// Patch F1 (Stored XSS render, Toast-путь).
// Raw-версия пробрасывала t.message напрямую в dangerouslySetInnerHTML,
// и любой серверный toast.push() с атакуемым сообщением давал XSS - даже
// если бек-санитайзер закрыл один конкретный путь хранения. Toast
// принимает "rich text" для форматирования, так что просто текстовая
// замена (innerText) сломала бы стилистику; вместо этого используем
// DOMPurify с ultra-minimal whitelist-ом (только inline-форматирование),
// без ALLOWED_ATTR - ни class, ни href, ни on*.
import DOMPurify from 'dompurify';
import { useCallback, useMemo, useState, type ReactNode } from 'react';
import styles from './Toast.module.css';
import { ToastContext, type ToastItem, type ToastKind } from './context';

const TOAST_SANITIZE_CONFIG = {
  ALLOWED_TAGS: ['b', 'i', 'em', 'strong', 'code', 'br'],
  ALLOWED_ATTR: [],
};

const kindClass: Record<ToastKind, string> = {
  info: styles['info']!,
  success: styles['success']!,
  error: styles['error']!,
};

const defaultTitles: Record<ToastKind, string> = {
  info: 'Извещение',
  success: 'Принято',
  error: 'Отклонено',
};

export const ToastProvider = ({ children }: { children: ReactNode }) => {
  const [toasts, setToasts] = useState<ToastItem[]>([]);

  const push = useCallback((toast: Omit<ToastItem, 'id'>) => {
    const id = `${Date.now()}-${Math.random().toString(36).slice(2, 8)}`;
    const item: ToastItem = { ...toast, id };
    setToasts((prev) => [...prev, item]);
    window.setTimeout(() => {
      setToasts((prev) => prev.filter((t) => t.id !== id));
    }, 5000);
  }, []);

  const value = useMemo(() => ({ push }), [push]);

  return (
    <ToastContext.Provider value={value}>
      {children}
      <div className={styles.container} role="status" aria-live="polite">
        {toasts.map((t) => (
          <div key={t.id} className={`${styles.toast} ${kindClass[t.kind]}`}>
            <span className={styles.title}>{t.title || defaultTitles[t.kind]}</span>
            <span
              className={styles.message}
              dangerouslySetInnerHTML={{
                __html: DOMPurify.sanitize(t.message, TOAST_SANITIZE_CONFIG),
              }}
            />
          </div>
        ))}
      </div>
    </ToastContext.Provider>
  );
};
