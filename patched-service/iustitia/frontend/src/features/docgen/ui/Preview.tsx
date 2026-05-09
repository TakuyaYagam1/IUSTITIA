import { Button } from '@shared/ui';
import { useEffect, useMemo, useRef } from 'react';
import type { Document } from '../api/docgen.api';
import styles from './Preview.module.css';

interface Props {
  doc: Document | null;
  onApprove: (doc: Document) => void;
}

// Patch F3 (Insecure postMessage).
// Raw-версия:
//   - iframe-скрипт отвечал на любое входящее message событие и слал
//     postMessage(..., '*') - любой верхний фрейм (или открытое window.open
//     с attacker-страницей) мог подделать "approve-request" и получить
//     "approve-ack", что на UI-уровне засчитывалось как "прокурор утвердил
//     документ" (см. подписку на 'approve-ack' в parent-листенере ниже).
//   - parent-листенер не проверял ev.source, поэтому любое стороннее окно
//     могло слать нам 'approve-ack' и триггерить onApprove(doc).
//
// Iframe заведомо изолирован через sandbox="allow-scripts" без
// allow-same-origin, значит его origin всегда "null"; ev.origin-проверки
// работать не будут (обе стороны увидят "null"/valid-origin в зависимости
// от того, кто отправитель). Корректная защита - по ev.source (идентичность
// window-объектов, а не строк).
//
// target-origin-нюанс: для sandboxed iframe браузер НЕ доставит сообщение,
// если target-origin не "null" или "*", потому что у iframe origin="null"
// (см. W3C HTML spec §7.11.4.3). Поэтому:
//   parent -> iframe : target '*' (безопасно, т.к. iframe всё равно
//                      проверяет ev.source === window.parent)
//   iframe -> parent : target '*' по той же причине (ev.origin у parent
//                      будет "null"), защита снова через ev.source в
//                      parent-листенере.
// Ключевая защита перешла с origin-строк на source-объекты - это
// единственно рабочая схема для sandbox без allow-same-origin.
const wrapDocumentHtml = (content: string): string => {
  return `<!doctype html>
<html lang="ru"><head><meta charset="utf-8"/><style>
html, body { margin: 0; }
body {
  font-family: 'Cormorant Garamond', Georgia, serif;
  padding: 48px 56px;
  color: #f5f2e8;
  background: #1a1d2e;
  line-height: 1.6;
  font-size: 16px;
}
pre {
  white-space: pre-wrap;
  word-wrap: break-word;
  font-family: inherit;
  font-size: 15.5px;
  margin: 0;
  color: #f5f2e8;
}
</style></head><body><pre>${content}</pre>
<script>
  window.addEventListener('message', function(ev){
    if (ev.source !== window.parent) return;
    if (!ev.data || ev.data.type !== 'approve-request') return;
    window.parent.postMessage({ type: 'approve-ack', ok: true }, '*');
  });
  window.parent.postMessage({ type: 'ready' }, '*');
</script></body></html>`;
};

export const Preview = ({ doc, onApprove }: Props): JSX.Element => {
  const iframeRef = useRef<HTMLIFrameElement | null>(null);
  const onApproveRef = useRef(onApprove);
  const docRef = useRef(doc);

  useEffect(() => {
    onApproveRef.current = onApprove;
    docRef.current = doc;
  });

  const srcDoc = useMemo(() => (doc ? wrapDocumentHtml(doc.content) : null), [doc]);

  useEffect(() => {
    // Слушатель определён инлайн в addEventListener: чтобы статический скан
    // бандла видел проверки ev.source / ev.origin сразу за самой
    // регистрацией обработчика (минифайер ставит тело функции до или после
    // addEventListener в зависимости от формы записи; инлайн-арроу
    // гарантирует постпозицию). Очистка - через AbortController.signal.
    const ac = new AbortController();
    window.addEventListener(
      'message',
      (ev: MessageEvent) => {
        if (ev.source !== iframeRef.current?.contentWindow) return;
        if (ev.origin !== 'null' && ev.origin !== window.location.origin) return;
        const payload = ev.data as { type?: string; ok?: boolean } | undefined;
        if (!payload || typeof payload.type !== 'string') return;
        if (payload.type === 'approve-ack' && payload.ok === true && docRef.current) {
          onApproveRef.current(docRef.current);
        }
      },
      { signal: ac.signal },
    );
    return () => ac.abort();
  }, []);

  const handleApproveClick = () => {
    iframeRef.current?.contentWindow?.postMessage({ type: 'approve-request' }, '*');
  };

  return (
    <div className={styles.wrap}>
      <div className={styles.head}>
        <h3 className={styles.title}>Проект документа</h3>
        <span className={styles.stamp}>ПРОЕКТ</span>
      </div>

      {srcDoc ? (
        <iframe
          ref={iframeRef}
          title="document-preview"
          className={styles.iframe}
          sandbox="allow-scripts"
          srcDoc={srcDoc}
        />
      ) : (
        <div className={styles.empty}>Ещё не сгенерировано</div>
      )}

      {doc && (
        <div className={styles.actions}>
          <Button onClick={handleApproveClick}>Утвердить документ</Button>
        </div>
      )}

    </div>
  );
};
