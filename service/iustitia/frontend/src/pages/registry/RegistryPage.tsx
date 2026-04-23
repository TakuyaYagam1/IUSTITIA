import { ROUTES } from '@/app/router/routes';
import { useCaseAccept, useCaseDismiss } from '@features/case-accept';
import { useCaseSearch, useUsersByRole } from '@features/case-view';
import { ApiError, type Case } from '@shared/api';
import { Button, useToast } from '@shared/ui';
import { Send, XCircle } from 'lucide-react';
import { useState } from 'react';
import { Link } from 'react-router-dom';
import styles from './RegistryPage.module.css';

interface RowProps {
  c: Case;
  prosecutors: { id: string; username: string; dome: string }[];
  onDone: () => void;
}

const Row = ({ c, prosecutors, onDone }: RowProps): JSX.Element => {
  const toast = useToast();
  const [prosecutorId, setProsecutorId] = useState<string>('');
  const [dismissOpen, setDismissOpen] = useState(false);
  const [reason, setReason] = useState('');

  const acceptMut = useCaseAccept(c.id);
  const dismissMut = useCaseDismiss(c.id);

  const handleAccept = async () => {
    if (!prosecutorId) {
      toast.push({
        kind: 'error',
        title: 'Отклонено',
        message: 'Выберите прокурора перед назначением.',
      });
      return;
    }
    try {
      await acceptMut.mutateAsync({ prosecutor_id: prosecutorId });
      toast.push({
        kind: 'success',
        title: 'Назначено',
        message: `Дело №${c.seq_num}/Δ передано прокурору.`,
      });
      onDone();
    } catch (err) {
      toast.push({
        kind: 'error',
        title: 'Ошибка',
        message: err instanceof ApiError ? err.message : 'Не удалось назначить.',
      });
    }
  };

  const handleDismiss = async () => {
    if (reason.trim().length < 3) {
      toast.push({ kind: 'error', title: 'Отклонено', message: 'Укажите причину.' });
      return;
    }
    try {
      await dismissMut.mutateAsync({ reason });
      toast.push({
        kind: 'success',
        title: 'Отклонено',
        message: `Дело №${c.seq_num}/Δ прекращено без суда.`,
      });
      setDismissOpen(false);
      setReason('');
      onDone();
    } catch (err) {
      toast.push({
        kind: 'error',
        title: 'Ошибка',
        message: err instanceof ApiError ? err.message : 'Не удалось отклонить.',
      });
    }
  };

  return (
    <div className={styles.row}>
      <span className={styles.rowNumber}>№{c.seq_num}/Δ</span>
      <div className={styles.rowBody}>
        <Link to={ROUTES.caseView(c.id)} style={{ textDecoration: 'none' }}>
          <h3 className={styles.rowDefendant}>{c.defendant}</h3>
        </Link>
        <p className={styles.rowCrime}>{c.crime}</p>
        <div className={styles.rowMeta}>
          <span>{new Date(c.created_at).toLocaleString('ru-RU')}</span>
          <span>статус: {c.status}</span>
        </div>
        {dismissOpen && (
          <div className={styles.dismissBox}>
            <textarea
              className={styles.reasonInput}
              value={reason}
              onChange={(e) => setReason(e.target.value)}
              placeholder="Причина отказа в рассмотрении"
              rows={3}
            />
            <div className={styles.dismissActions}>
              <Button variant="ghost" onClick={() => setDismissOpen(false)}>
                Отмена
              </Button>
              <Button variant="primary" onClick={handleDismiss} disabled={dismissMut.isPending}>
                {dismissMut.isPending ? 'Отклонение…' : 'Подтвердить отказ'}
              </Button>
            </div>
          </div>
        )}
      </div>
      <div className={styles.rowActions}>
        <select
          className={styles.select}
          value={prosecutorId}
          onChange={(e) => setProsecutorId(e.target.value)}
          aria-label="Выбрать прокурора"
        >
          <option value="">- прокурор -</option>
          {prosecutors.map((p) => (
            <option key={p.id} value={p.id}>
              {p.username} ({p.dome})
            </option>
          ))}
        </select>
        <Button variant="primary" onClick={handleAccept} disabled={acceptMut.isPending}>
          <Send size={14} /> Принять
        </Button>
        {!dismissOpen && (
          <Button variant="ghost" onClick={() => setDismissOpen(true)}>
            <XCircle size={14} /> Отклонить
          </Button>
        )}
      </div>
    </div>
  );
};

export const RegistryPage = (): JSX.Element => {
  const { data, isLoading, isError, error, refetch } = useCaseSearch();
  const prosecutorsQuery = useUsersByRole('prosecutor');
  const prosecutors = prosecutorsQuery.data ?? [];

  const pending = (data ?? []).filter((c) => c.status === 'draft' || c.status === 'open');

  return (
    <div className={styles.wrap}>
      <div className={styles.header}>
        <h1 className={styles.title}>Канцелярия</h1>
        <span className={styles.subtitle}>Входящие дела · Регистратор</span>
      </div>

      {isLoading && <div className={styles.message}>Получение картотеки…</div>}
      {isError && (
        <div className={`${styles.message} ${styles.error}`}>
          Ошибка: {error instanceof Error ? error.message : 'неизвестно'}
        </div>
      )}
      {!isLoading && !isError && pending.length === 0 && (
        <div className={styles.message}>Нет дел, ожидающих назначения.</div>
      )}

      <div className={styles.list}>
        {pending.map((c) => (
          <Row key={c.id} c={c} prosecutors={prosecutors} onDone={() => refetch()} />
        ))}
      </div>
    </div>
  );
};
