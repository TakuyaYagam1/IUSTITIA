import { Button } from '@shared/ui';
import { useEffect, useMemo, useRef } from 'react';
import type { Document } from '../api/docgen.api';
import styles from './Preview.module.css';

interface Props {
  doc: Document | null;
  onApprove: (doc: Document) => void;
}

const wrapDocumentHtml = (content: string): string => `<!doctype html>
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
    if (ev.data && ev.data.type === 'approve-request') {
      window.parent.postMessage({ type: 'approve-ack', ok: true }, '*');
    }
  });
  window.parent.postMessage({ type: 'ready' }, '*');
</script></body></html>`;

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
    const handler = (ev: MessageEvent) => {
      const payload = ev.data as { type?: string; ok?: boolean } | undefined;
      if (payload?.type === 'approve-ack' && payload.ok && docRef.current) {
        onApproveRef.current(docRef.current);
      }
    };
    window.addEventListener('message', handler);
    return () => window.removeEventListener('message', handler);
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
